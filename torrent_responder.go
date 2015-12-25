package main

import (
	"bytes"
	"fmt"
	"github.com/Syfaro/telegram-bot-api"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/dustin/go-humanize"
	"github.com/zhulik/margelet"
	"runtime"
	"time"
)

func infoAsString(info *metainfo.MetaInfo) string {
	return fmt.Sprintf("Name: %s, Size: %s", info.Info.Name, humanize.Bytes(uint64(info.Info.TotalLength())))
}

type TorrentResponder struct {
	client             *torrent.Client
	torrentsRepository *TorrentsRepository
	authorizedUsername string
}

func NewTorrentResponder(authorizedUsername string, client *torrent.Client, repo *TorrentsRepository) (responder *TorrentResponder, err error) {
	responder = &TorrentResponder{client, repo, authorizedUsername}
	runtime.SetFinalizer(responder, func(t *TorrentResponder) {
		t.client.Close()
	})
	return
}

func (session TorrentResponder) Response(bot margelet.MargeletAPI, message tgbotapi.Message) error {
	if message.From.UserName != session.authorizedUsername {
		bot.QuickSend(message.Chat.ID, "Sorry, you are not allowed to control me!")
		return nil
	}

	document := message.Document

	if len(document.FileID) > 0 && document.MimeType == "application/x-bittorrent" {
		bot.Send(tgbotapi.NewChatAction(message.Chat.ID, tgbotapi.ChatTyping))

		url, err := bot.GetFileDirectURL(document.FileID)
		if err != nil {
			return err
		}

		data, err := download(url)
		if err != nil {
			return err
		}

		info, err := metainfo.Load(bytes.NewReader(data))
		if err != nil {
			return err
		}

		session.torrentsRepository.Add(message.Chat.ID, message.From.ID, data)

		bot.QuickSend(message.Chat.ID, infoAsString(info))
		bot.GetSessionRepository().Create(message.Chat.ID, message.From.ID, "/download")
		bot.HandleSession(message, session)

		return nil
	}
	return nil
}

func (session TorrentResponder) HandleResponse(bot margelet.MargeletAPI, message tgbotapi.Message, responses []string) (bool, error) {
	if message.From.UserName != session.authorizedUsername {
		bot.QuickSend(message.Chat.ID, "Sorry, you are not allowed to control me!")
		return true, nil
	}
	switch len(responses) {
	case 0:
		if session.torrentsRepository.Exists(message.Chat.ID, message.From.ID) {
			msg := tgbotapi.NewMessage(message.Chat.ID, "Would you like to download it?")
			msg.ReplyMarkup = tgbotapi.ReplyKeyboardMarkup{
				[][]string{[]string{"yes", "no"}},
				true,
				true,
				true,
			}
			bot.Send(msg)
			return false, nil
		}
		bot.QuickSend(message.Chat.ID, "Please send me torrent file, do not use /download directly")

		return true, nil
	case 1:
		switch message.Text {
		case "yes":
			data, err := session.torrentsRepository.Get(message.Chat.ID, message.From.ID)
			if err != nil {
				bot.QuickSend(message.Chat.ID, "Sorry, something is wrong. Please, try again")
				return true, nil
			}
			session.torrentsRepository.Delete(message.Chat.ID, message.From.ID)
			downloadTorrent(bot, message.Chat.ID, data, session.client)
			return true, nil
		case "no":
			bot.QuickSend(message.Chat.ID, "Ok, i will forgive it")
			session.torrentsRepository.Delete(message.Chat.ID, message.From.ID)
			return true, nil
		default:
			bot.QuickSend(message.Chat.ID, "Sorry, i don't understand, yes or no?")
			return false, fmt.Errorf("unknown answer")
		}
	default:
		return false, fmt.Errorf("never achieved")
	}
}

func (session TorrentResponder) HelpMessage() string {
	return "Download torrent, please do not use it directly"
}

func downloadTorrent(bot margelet.MargeletAPI, chatID int, data []byte, client *torrent.Client) error {
	info, err := metainfo.Load(bytes.NewReader(data))
	if err != nil {
		return err
	}
	t, err := client.AddTorrent(info)
	if err != nil {
		return err
	}

	t.DownloadAll()
	bot.QuickSend(chatID, fmt.Sprintf("%s is Downloading...", t.Name()))

	go func() {
		for t.BytesCompleted() != t.Info().TotalLength() {
			time.Sleep(time.Second)
		}
		bot.QuickSend(chatID, fmt.Sprintf("Downloading of %s is finished!", t.Name()))
	}()

	return nil
}

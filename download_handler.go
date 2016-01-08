package main

import (
	"bytes"
	"fmt"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/zhulik/margelet"
	"time"
)

type downloadHandler struct {
	client             torrentClient
	torrentsRepository *torrentsRepository
	authorizedUsername string
}

func newDownloadHandler(authorizedUsername string, client torrentClient, repo *torrentsRepository) (responder *downloadHandler, err error) {
	responder = &downloadHandler{client, repo, authorizedUsername}
	return
}

func (session downloadHandler) HelpMessage() string {
	return "Download torrent, please do not use it directly"
}

func (session downloadHandler) HandleResponse(bot margelet.MargeletAPI, message tgbotapi.Message, responses []tgbotapi.Message) (bool, error) {
	if message.From.UserName != session.authorizedUsername {
		bot.QuickSend(message.Chat.ID, "Sorry, you are not allowed to control me!")
		return true, nil
	}
	switch len(responses) {
	case 0:
		if session.torrentsRepository.Exists(message.Chat.ID, message.From.ID) {
			msg := tgbotapi.NewMessage(message.Chat.ID, "Would you like to download it?")
			msg.ReplyMarkup = yesNoReplyMarkup
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

			return true, downloadTorrent(bot, message.Chat.ID, data, session.client)
		case "no":
			msg := tgbotapi.NewMessage(message.Chat.ID, "Ok, i will forgive it")
			msg.ReplyMarkup = hideReplyMarkup
			bot.Send(msg)
			session.torrentsRepository.Delete(message.Chat.ID, message.From.ID)
			return true, nil
		default:
			msg := tgbotapi.NewMessage(message.Chat.ID, "Sorry, i don't understand.")
			msg.ReplyMarkup = yesNoReplyMarkup
			bot.Send(msg)
			return false, fmt.Errorf("unknown answer")
		}
	default:
		return false, fmt.Errorf("never achieved")
	}
}

func downloadTorrent(bot margelet.MargeletAPI, chatID int, data []byte, client torrentClient) error {
	str := string(data)

	if magnetRE.MatchString(str) {
		return downloadMagnet(bot, chatID, str, client)
	}

	return downloadTorrentFile(bot, chatID, data, client)
}

func run(t Torrent, chatID int, bot margelet.MargeletAPI) error {
	t.DownloadAll()
	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("%s is downloading...", t.Info().Name))
	msg.ReplyMarkup = hideReplyMarkup
	bot.Send(msg)

	go func() {
		for t.BytesCompleted() != t.Info().TotalLength() {
			time.Sleep(time.Second)
		}
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Downloading of %s is finished!", t.Info().Name))
		msg.ReplyMarkup = hideReplyMarkup
		bot.Send(msg)
	}()

	return nil
}

func downloadTorrentFile(bot margelet.MargeletAPI, chatID int, data []byte, client torrentClient) error {
	info, err := metainfo.Load(bytes.NewReader(data))
	if err != nil {
		return err
	}
	t, err := client.AddTorrent(info)
	if err != nil {
		return err
	}

	return run(t, chatID, bot)
}

func downloadMagnet(bot margelet.MargeletAPI, chatID int, url string, client torrentClient) error {
	t, err := client.AddMagnet(url)
	if err != nil {
		return err
	}

	return run(t, chatID, bot)
}

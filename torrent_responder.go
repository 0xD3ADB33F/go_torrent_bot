package main

import (
	"bytes"
	"fmt"
	"github.com/Syfaro/telegram-bot-api"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/zhulik/margelet"
	"regexp"
	"time"
)

var (
	magnetRE, _      = regexp.Compile(`^magnet:\?xt=urn:.+$`)
	yesNoReplyMarkup = tgbotapi.ReplyKeyboardMarkup{
		[][]string{{"yes", "no"}},
		true,
		true,
		true,
	}
	hideReplyMarkup = tgbotapi.ReplyKeyboardHide{
		true,
		true,
	}
)

type torrentResponder struct {
	client             *torrent.Client
	torrentsRepository *torrentsRepository
	authorizedUsername string
}

func newTorrentResponder(authorizedUsername string, client *torrent.Client, repo *torrentsRepository) (responder *torrentResponder, err error) {
	responder = &torrentResponder{client, repo, authorizedUsername}
	return
}

func (session torrentResponder) handleTorrentFile(bot margelet.MargeletAPI, message tgbotapi.Message) error {
	bot.Send(tgbotapi.NewChatAction(message.Chat.ID, tgbotapi.ChatTyping))

	url, err := bot.GetFileDirectURL(message.Document.FileID)
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

	bot.QuickSend(message.Chat.ID, infoAsString(info))
	session.torrentsRepository.Add(message.Chat.ID, message.From.ID, data)
	bot.GetSessionRepository().Create(message.Chat.ID, message.From.ID, "/download")
	bot.HandleSession(message, session)
	return nil
}

func (session torrentResponder) handleMagnetLink(bot margelet.MargeletAPI, message tgbotapi.Message) error {
	bot.Send(tgbotapi.NewChatAction(message.Chat.ID, tgbotapi.ChatTyping))

	t, err := session.client.AddMagnet(message.Text)

	if err != nil {
		return err
	}

	<-t.GotInfo()
	bot.QuickSend(message.Chat.ID, infoAsString(t.MetaInfo()))
	session.torrentsRepository.Add(message.Chat.ID, message.From.ID, []byte(message.Text))
	bot.GetSessionRepository().Create(message.Chat.ID, message.From.ID, "/download")
	bot.HandleSession(message, session)
	t.Drop()

	return nil
}

func (session torrentResponder) Response(bot margelet.MargeletAPI, message tgbotapi.Message) error {
	if message.From.UserName != session.authorizedUsername {
		bot.QuickSend(message.Chat.ID, "Sorry, you are not allowed to control me!")
		return nil
	}

	if len(message.Document.FileID) > 0 && message.Document.MimeType == "application/x-bittorrent" {
		return session.handleTorrentFile(bot, message)
	}

	if magnetRE.MatchString(message.Text) {
		return session.handleMagnetLink(bot, message)
	}

	return nil
}

func (session torrentResponder) HandleResponse(bot margelet.MargeletAPI, message tgbotapi.Message, responses []tgbotapi.Message) (bool, error) {
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

func (session torrentResponder) HelpMessage() string {
	return "Download torrent, please do not use it directly"
}

func downloadTorrent(bot margelet.MargeletAPI, chatID int, data []byte, client *torrent.Client) error {
	str := string(data)

	if magnetRE.MatchString(str) {
		return downloadMagnet(bot, chatID, str, client)
	}

	return downloadTorrentFile(bot, chatID, data, client)
}

func run(t torrent.Torrent, chatID int, bot margelet.MargeletAPI) error {
	t.DownloadAll()
	bot.QuickSend(chatID, fmt.Sprintf("%s is downloading...", t.Name()))

	go func() {
		for t.BytesCompleted() != t.Info().TotalLength() {
			time.Sleep(time.Second)
		}
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Downloading of %s is finished!", t.Name()))
		msg.ReplyMarkup = hideReplyMarkup
		bot.Send(msg)
	}()

	return nil
}

func downloadTorrentFile(bot margelet.MargeletAPI, chatID int, data []byte, client *torrent.Client) error {
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

func downloadMagnet(bot margelet.MargeletAPI, chatID int, url string, client *torrent.Client) error {
	t, err := client.AddMagnet(url)
	if err != nil {
		return err
	}

	return run(t, chatID, bot)
}

package main

import (
	"bytes"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/zhulik/margelet"
)

var (
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
	client             torrentClient
	torrentsRepository *torrentsRepository
	authorizedUsername string
	downloadSession    margelet.SessionHandler
	download           Downloader
}

func newTorrentResponder(authorizedUsername string, client torrentClient,
	repo *torrentsRepository, downloadSession margelet.SessionHandler, d Downloader) (responder *torrentResponder, err error) {
	responder = &torrentResponder{client, repo, authorizedUsername, downloadSession, d}
	return
}

func (session torrentResponder) Response(bot margelet.MargeletAPI, message tgbotapi.Message) error {
	if message.From.UserName != session.authorizedUsername {
		bot.QuickSend(message.Chat.ID, "Sorry, you are not allowed to control me!")
		return nil
	}

	if len(message.Document.FileID) == 0 || message.Document.MimeType != "application/x-bittorrent" {
		return nil
	}

	bot.Send(tgbotapi.NewChatAction(message.Chat.ID, tgbotapi.ChatTyping))

	url, err := bot.GetFileDirectURL(message.Document.FileID)
	if err != nil {
		return err
	}

	data, err := session.download(url)
	if err != nil {
		return err
	}

	info, err := metainfo.Load(bytes.NewReader(data))
	if err != nil {
		return err
	}

	bot.QuickSend(message.Chat.ID, infoAsString(&info.Info.Info))
	session.torrentsRepository.Add(message.Chat.ID, message.From.ID, data)
	bot.GetSessionRepository().Create(message.Chat.ID, message.From.ID, "/download")
	bot.HandleSession(message, session.downloadSession)
	return nil
}

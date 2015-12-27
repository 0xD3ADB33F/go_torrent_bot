package main

import (
	"github.com/Syfaro/telegram-bot-api"
	"github.com/zhulik/margelet"
	"regexp"
)

var (
	magnetRE, _ = regexp.Compile(`^magnet:\?xt=urn:.+$`)
)

type magnetResponder struct {
	client             torrentClient
	torrentsRepository *torrentsRepository
	authorizedUsername string
	downloadSession    margelet.SessionHandler
}

func newMagnetResponder(authorizedUsername string, client torrentClient, repo *torrentsRepository, downloadSession margelet.SessionHandler) (responder *magnetResponder, err error) {
	responder = &magnetResponder{client, repo, authorizedUsername, downloadSession}
	return
}

func (session magnetResponder) Response(bot margelet.MargeletAPI, message tgbotapi.Message) error {
	if message.From.UserName != session.authorizedUsername {
		bot.QuickSend(message.Chat.ID, "Sorry, you are not allowed to control me!")
		return nil
	}

	if !magnetRE.MatchString(message.Text) {
		return nil
	}

	bot.Send(tgbotapi.NewChatAction(message.Chat.ID, tgbotapi.ChatTyping))

	t, err := session.client.AddMagnet(message.Text)

	if err != nil {
		return err
	}
	defer t.Drop()

	<-t.GotInfo()
	bot.QuickSend(message.Chat.ID, infoAsString(t.Info()))
	session.torrentsRepository.Add(message.Chat.ID, message.From.ID, []byte(message.Text))
	bot.GetSessionRepository().Create(message.Chat.ID, message.From.ID, "/download")
	bot.HandleSession(message, session.downloadSession)

	return nil
}

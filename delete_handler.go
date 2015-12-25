package main

import (
	"fmt"
	"github.com/Syfaro/telegram-bot-api"
	"github.com/anacrolix/torrent"
	"github.com/zhulik/margelet"
	"strings"
)

var (
	yesNoCancelReplyMarkup = tgbotapi.ReplyKeyboardMarkup{
		[][]string{[]string{"yes", "no"}, []string{"cancel"}},
		true,
		true,
		true,
	}
)

type DeleteHandler struct {
	client             *torrent.Client
	authorizedUsername string
}

func NewDeleteHandler(authorizedUsername string, client *torrent.Client) *DeleteHandler {
	return &DeleteHandler{client, authorizedUsername}
}

func (handler DeleteHandler) HelpMessage() string {
	return "Deletes torrent without files ot with it"
}

func (handler DeleteHandler) handleDeleteCommand(bot margelet.MargeletAPI, message tgbotapi.Message) (bool, error) {
	hash := strings.TrimSpace(message.CommandArguments())
	if len(hash) > 0 {
		torrents := handler.client.Torrents()
		index := indexByHexHash(hash, torrents)
		if index != -1 {
			t := torrents[index]
			msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("You trying to download %s.\nWould you like to remove downloaded files?", infoAsString(t.MetaInfo())))
			msg.ReplyMarkup = yesNoCancelReplyMarkup
			bot.Send(msg)
			return false, nil
		}
		bot.QuickSend(message.Chat.ID, fmt.Sprintf("Cannot find download with hash %s", hash))
	}
	return true, nil
}

func (handler DeleteHandler) handleAnswer(bot margelet.MargeletAPI, message tgbotapi.Message) (bool, error) {
	switch message.Text {
	case "yes":
	case "no":
	case "cancel":
		return true, nil
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, "Sorry, i don't understand.")
	msg.ReplyMarkup = yesNoCancelReplyMarkup
	bot.Send(msg)
	return false, fmt.Errorf("unknown answer")
}

func (handler DeleteHandler) HandleResponse(bot margelet.MargeletAPI, message tgbotapi.Message, responses []string) (bool, error) {
	if message.From.UserName != handler.authorizedUsername {
		bot.QuickSend(message.Chat.ID, "Sorry, you are not allowed to control me!")
		return true, nil
	}

	switch len(responses) {
	case 0:
		return handler.handleDeleteCommand(bot, message)
	case 1:
		return handler.handleAnswer(bot, message)
	}

	return true, nil
}

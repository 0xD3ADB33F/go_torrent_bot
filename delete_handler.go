package main

import (
	"fmt"
	"github.com/Syfaro/telegram-bot-api"
	"github.com/anacrolix/torrent"
	"github.com/zhulik/margelet"
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
	torrent, err := findTorrentByMessage(handler.client, message)

	if err != nil {
		bot.QuickSend(message.Chat.ID, fmt.Sprintf("usage: /delete <download hash>"))
		return true, nil
	}

	if torrent != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("You trying to delete %s.\nWould you like to remove downloaded files?", infoAsString(torrent.MetaInfo())))
		msg.ReplyMarkup = yesNoCancelReplyMarkup
		bot.Send(msg)
		return false, nil
	}
	bot.QuickSend(message.Chat.ID, fmt.Sprintf("Cannot find download with hash %s", message.CommandArguments()))
	return true, nil

}

func (handler DeleteHandler) handleAnswer(bot margelet.MargeletAPI, prevMessage tgbotapi.Message, message tgbotapi.Message) (bool, error) {
	if message.Text == "cancel" {
		bot.QuickSend(message.Chat.ID, "Delete canceled!")
		return true, nil
	}

	switch message.Text {
	case "yes":
	case "no":
		fmt.Print("TEST")
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, "Sorry, i don't understand.")
	msg.ReplyMarkup = yesNoCancelReplyMarkup
	bot.Send(msg)
	return false, fmt.Errorf("unknown answer")
}

func (handler DeleteHandler) HandleResponse(bot margelet.MargeletAPI, message tgbotapi.Message, responses []tgbotapi.Message) (bool, error) {
	if message.From.UserName != handler.authorizedUsername {
		bot.QuickSend(message.Chat.ID, "Sorry, you are not allowed to control me!")
		return true, nil
	}

	switch len(responses) {
	case 0:
		return handler.handleDeleteCommand(bot, message)
	case 1:
		fmt.Println(responses)
		return handler.handleAnswer(bot, responses[0], message)
	}

	return true, nil
}

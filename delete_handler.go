package main

import (
	"fmt"
	"github.com/Syfaro/telegram-bot-api"
	"github.com/anacrolix/torrent"
	"github.com/zhulik/margelet"
	"os"
	"path"
)

var (
	yesNoCancelReplyMarkup = tgbotapi.ReplyKeyboardMarkup{
		[][]string{[]string{"yes", "no"}, []string{"cancel"}},
		true,
		true,
		true,
	}
)

type deleteHandler struct {
	client             *torrent.Client
	authorizedUsername string
	downloadPath       string
}

func newDeleteHandler(authorizedUsername, downloadPath string, client *torrent.Client) *deleteHandler {
	return &deleteHandler{client, authorizedUsername, downloadPath}
}

func (handler deleteHandler) HelpMessage() string {
	return "Deletes torrent without files ot with it"
}

func (handler deleteHandler) handleDeleteCommand(bot margelet.MargeletAPI, message tgbotapi.Message) (bool, error) {
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

func (handler deleteHandler) handleAnswer(bot margelet.MargeletAPI, prevMessage tgbotapi.Message, message tgbotapi.Message) (bool, error) {
	if message.Text == "cancel" {
		bot.QuickSend(message.Chat.ID, "Delete canceled!")
		return true, nil
	}

	torrent, _ := findTorrentByMessage(handler.client, prevMessage)

	switch message.Text {
	case "yes":
		torrent.Drop()
		go func() {
			err := os.RemoveAll(path.Join(handler.downloadPath, torrent.Info().Name))
			if err != nil {
				bot.QuickSend(message.Chat.ID, fmt.Sprintf("Sorry, something went wrong when i trying to delete %s files!", infoAsString(torrent.MetaInfo())))
			}
		}()
		bot.QuickSend(message.Chat.ID, fmt.Sprintf("Downloading of %s canceled, files removed!", infoAsString(torrent.MetaInfo())))
		return true, nil
	case "no":
		bot.QuickSend(message.Chat.ID, fmt.Sprintf("Downloading of %s canceled!", infoAsString(torrent.MetaInfo())))
		torrent.Drop()
		return true, nil
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, "Sorry, i don't understand.")
	msg.ReplyMarkup = yesNoCancelReplyMarkup
	bot.Send(msg)
	return false, fmt.Errorf("unknown answer")
}

func (handler deleteHandler) HandleResponse(bot margelet.MargeletAPI, message tgbotapi.Message, responses []tgbotapi.Message) (bool, error) {
	if message.From.UserName != handler.authorizedUsername {
		bot.QuickSend(message.Chat.ID, "Sorry, you are not allowed to control me!")
		return true, nil
	}

	switch len(responses) {
	case 0:
		return handler.handleDeleteCommand(bot, message)
	case 1:
		return handler.handleAnswer(bot, responses[0], message)
	}

	return true, nil
}

package main

import (
	"fmt"
	"github.com/Syfaro/telegram-bot-api"
	"github.com/anacrolix/torrent"
	"github.com/zhulik/margelet"
	"strings"
)

type StatusHandler struct {
	path               string
	client             *torrent.Client
	authorizedUsername string
}

func NewStatusHandler(authorizedUsername string, path string, client *torrent.Client) *StatusHandler {
	return &StatusHandler{path, client, authorizedUsername}
}

func (responder StatusHandler) Response(bot margelet.MargeletAPI, message tgbotapi.Message) error {
	if message.From.UserName != responder.authorizedUsername {
		bot.QuickSend(message.Chat.ID, "Sorry, you are not allowed to control me!")
		return nil
	}

	if len(responder.client.Torrents()) == 0 {
		bot.QuickSend(message.Chat.ID, "There is no downloads")
		return nil
	}

	hash := strings.TrimSpace(message.CommandArguments())
	if len(hash) > 0 {
		torrents := responder.client.Torrents()
		index := indexByHexHash(hash, torrents)
		if index != -1 {
			t := torrents[index]
			bot.QuickSend(message.Chat.ID, verboseTorrentStats(responder.path, t))
			return nil
		}
		bot.QuickSend(message.Chat.ID, fmt.Sprintf("Cannot find download with hash %s", hash))

		return nil
	}

	for _, t := range responder.client.Torrents() {
		bot.QuickSend(message.Chat.ID, torrentStats(t))
	}
	return nil
}

func (responder StatusHandler) HelpMessage() string {
	return "Shows status of your downloads"
}

package main

import (
	"fmt"
	"github.com/Syfaro/telegram-bot-api"
	"github.com/anacrolix/torrent"
	"github.com/dustin/go-humanize"
	"github.com/zhulik/margelet"
	"path"
	"strings"
)

func torrentStats(torrent torrent.Torrent) string {
	return fmt.Sprintf("Hash: %s, Name: %s, Size: %s, Progress: %.2f%%",
		torrent.InfoHash().HexString(),
		torrent.Info().Name,
		humanize.Bytes(uint64(torrent.Info().TotalLength())),
		float64(torrent.BytesCompleted())/float64(torrent.Info().TotalLength())*100,
	)
}

type StatusResponder struct {
	path               string
	client             *torrent.Client
	authorizedUsername string
}

func NewStatusResponder(authorizedUsername string, path string, client *torrent.Client) *StatusResponder {
	return &StatusResponder{path, client, authorizedUsername}
}

func indexByHexHash(hash string, torrents []torrent.Torrent) int {
	for index, t := range torrents {
		if t.InfoHash().HexString() == hash {
			return index
		}
	}
	return -1
}

func (responder StatusResponder) Response(bot margelet.MargeletAPI, message tgbotapi.Message) error {
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
			bot.QuickSend(message.Chat.ID, responder.verboseTorrentStats(t))
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

func (responder StatusResponder) HelpMessage() string {
	return "Shows status of your downloads"
}

func (responder StatusResponder) verboseTorrentStats(t torrent.Torrent) string {
	return fmt.Sprintf("Hash: %s\nName: %s\nSize: %s\nProgress: %.2f%%\nSeeding: %t\nPeers: %d\nLocation: %s",
		t.InfoHash().HexString(),
		t.Info().Name,
		humanize.Bytes(uint64(t.Info().TotalLength())),
		float64(t.BytesCompleted())/float64(t.Info().TotalLength())*100,
		t.Seeding(),
		len(t.Peers),
		path.Join(responder.path, t.Info().Name),
	)
}

package main

import (
	"fmt"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/dustin/go-humanize"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
)

type Downloader func(url string) ([]byte, error)

func download(url string) ([]byte, error) {
	response, err := http.Get(url)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("Respose code != 200")
	}

	return ioutil.ReadAll(response.Body)
}

func indexByHexHash(hash string, torrents []Torrent) int {
	for index, t := range torrents {
		if t.InfoHash().HexString() == hash {
			return index
		}
	}
	return -1
}

func verboseTorrentStats(downloadPath string, t Torrent) string {
	return fmt.Sprintf("%s\nName: %s\nSize: %s\nProgress: %.2f%%\nSeeding: %t\nLocation: %s",
		t.InfoHash().HexString(),
		t.Info().Name,
		humanize.Bytes(uint64(t.Info().TotalLength())),
		float64(t.BytesCompleted())/float64(t.Info().TotalLength())*100,
		t.Seeding(),
		path.Join(downloadPath, t.Info().Name),
	)
}

func torrentStats(torrent Torrent) string {
	return fmt.Sprintf("%s\nName: %s, Size: %s, Progress: %.2f%%",
		torrent.InfoHash().HexString(),
		torrent.Info().Name,
		humanize.Bytes(uint64(torrent.Info().TotalLength())),
		float64(torrent.BytesCompleted())/float64(torrent.Info().TotalLength())*100,
	)
}

func infoAsString(info *metainfo.Info) string {
	return fmt.Sprintf("Name: %s, Size: %s", info.Name, humanize.Bytes(uint64(info.TotalLength())))
}

func findTorrent(client torrentClient, hash string) Torrent {
	torrents := client.Torrents()
	index := indexByHexHash(hash, torrents)
	if index != -1 {
		return torrents[index]
	}
	return nil
}

func findHash(message tgbotapi.Message) string {
	if message.IsCommand() {
		return strings.TrimSpace(message.CommandArguments())
	}
	lines := strings.SplitN(message.Text, "\n", 2)
	if len(lines) == 2 {
		return strings.TrimSpace(lines[0])
	}
	return ""
}

func findTorrentByMessage(client torrentClient, message tgbotapi.Message) (Torrent, string, error) {
	hash := findHash(message)
	if len(hash) > 0 {
		if torrent := findTorrent(client, hash); torrent != nil {
			return torrent, hash, nil
		}
		return nil, hash, nil
	}
	return nil, hash, fmt.Errorf("No hash in message")
}

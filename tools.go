package main

import (
	"fmt"
	"github.com/Syfaro/telegram-bot-api"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/dustin/go-humanize"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
)

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

func indexByHexHash(hash string, torrents []torrent.Torrent) int {
	for index, t := range torrents {
		if t.InfoHash().HexString() == hash {
			return index
		}
	}
	return -1
}

func verboseTorrentStats(downloadPath string, t torrent.Torrent) string {
	return fmt.Sprintf("%s\nName: %s\nSize: %s\nProgress: %.2f%%\nSeeding: %t\nPeers: %d\nLocation: %s",
		t.InfoHash().HexString(),
		t.Info().Name,
		humanize.Bytes(uint64(t.Info().TotalLength())),
		float64(t.BytesCompleted())/float64(t.Info().TotalLength())*100,
		t.Seeding(),
		len(t.Peers),
		path.Join(downloadPath, t.Info().Name),
	)
}

func torrentStats(torrent torrent.Torrent) string {
	return fmt.Sprintf("%s\nName: %s, Size: %s, Progress: %.2f%%",
		torrent.InfoHash().HexString(),
		torrent.Info().Name,
		humanize.Bytes(uint64(torrent.Info().TotalLength())),
		float64(torrent.BytesCompleted())/float64(torrent.Info().TotalLength())*100,
	)
}

func infoAsString(info *metainfo.MetaInfo) string {
	return fmt.Sprintf("Name: %s, Size: %s", info.Info.Name, humanize.Bytes(uint64(info.Info.TotalLength())))
}

func findTorrent(client *torrent.Client, hash string) *torrent.Torrent {
	torrents := client.Torrents()
	index := indexByHexHash(hash, torrents)
	if index != -1 {
		return &torrents[index]
	}
	return nil
}

func findTorrentByMessage(client *torrent.Client, message tgbotapi.Message) (*torrent.Torrent, error) {
	hash := strings.TrimSpace(message.CommandArguments())
	if len(hash) > 0 {
		if torrent := findTorrent(client, hash); torrent != nil {
			return torrent, nil
		}
		return nil, nil
	}
	return nil, fmt.Errorf("No hash in message")
}

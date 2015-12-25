package main

import (
	"fmt"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/dustin/go-humanize"
	"io/ioutil"
	"net/http"
	"path"
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
	return fmt.Sprintf("Hash: %s\nName: %s\nSize: %s\nProgress: %.2f%%\nSeeding: %t\nPeers: %d\nLocation: %s",
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
	return fmt.Sprintf("Hash: %s, Name: %s, Size: %s, Progress: %.2f%%",
		torrent.InfoHash().HexString(),
		torrent.Info().Name,
		humanize.Bytes(uint64(torrent.Info().TotalLength())),
		float64(torrent.BytesCompleted())/float64(torrent.Info().TotalLength())*100,
	)
}

func infoAsString(info *metainfo.MetaInfo) string {
	return fmt.Sprintf("Name: %s, Size: %s", info.Info.Name, humanize.Bytes(uint64(info.Info.TotalLength())))
}

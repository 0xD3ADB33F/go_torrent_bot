package main

import (
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
)

type torrentClient interface {
	Torrents() []torrent.Download
	AddMagnet(url string) (torrent.Download, error)
	AddTorrent(info *metainfo.MetaInfo) (torrent.Download, error)
}

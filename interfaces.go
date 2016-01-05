package main

import (
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
)

type torrentClient interface {
	Torrents() []torrent.Torrent
	AddMagnet(url string) (torrent.Torrent, error)
	AddTorrent(info *metainfo.MetaInfo) (torrent.Torrent, error)
}

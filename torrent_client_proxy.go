package main

import (
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
)

type torrentClientProxy struct {
	client *torrent.Client
}

func (p *torrentClientProxy) Torrents() (result []Torrent) {
	torrents := p.client.Torrents()

	for _, t := range torrents {
		result = append(result, t)
	}

	return
}

func (p *torrentClientProxy) AddMagnet(url string) (Torrent, error) {
	return p.client.AddMagnet(url)
}

func (p *torrentClientProxy) AddTorrent(info *metainfo.MetaInfo) (Torrent, error) {
	return p.client.AddTorrent(info)
}

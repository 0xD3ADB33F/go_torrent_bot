package main

import (
	"github.com/anacrolix/missinggo/pubsub"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
)

type Torrent interface {
	InfoHash() torrent.InfoHash
	GotInfo() <-chan struct{}
	Info() *metainfo.Info
	NewReader() (ret *torrent.Reader)
	PieceStateRuns() []torrent.PieceStateRun
	NumPieces() int
	Drop()
	BytesCompleted() int64
	SubscribePieceStateChanges() *pubsub.Subscription
	Seeding() bool
	SetDisplayName(dn string)
	AddPeers(pp []torrent.Peer) error
	DownloadAll()
	Files() (ret []torrent.File)
}

type torrentClient interface {
	Torrents() []Torrent
	AddMagnet(url string) (Torrent, error)
	AddTorrent(info *metainfo.MetaInfo) (Torrent, error)
}

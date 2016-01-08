package main

import (
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/anacrolix/missinggo/pubsub"
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
//	Client() *torrent.Client
	AddPeers(pp []torrent.Peer) error
	DownloadAll()
//	Trackers() [][]tracker.Client
	Files() (ret []torrent.File)
//	Peers() map[PeersKey]Peer
}

type torrentClient interface {
	Torrents() []Torrent
	AddMagnet(url string) (Torrent, error)
	AddTorrent(info *metainfo.MetaInfo) (Torrent, error)
}

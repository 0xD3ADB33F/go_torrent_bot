package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/anacrolix/missinggo/pubsub"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/anacrolix/torrent/tracker"
	"github.com/zhulik/margelet"
	"gopkg.in/redis.v3"
)

type TorrentClientMock struct {
	torrents []torrent.Download
}

func newTorrentClientMock() *TorrentClientMock {
	return &TorrentClientMock{}
}

func (mock TorrentClientMock) AddMagnet(string) (torrent.Download, error) {
	return mock.torrents[0], nil
}

func (mock TorrentClientMock) AddTorrent(info *metainfo.MetaInfo) (torrent.Download, error) {
	return mock.torrents[0], nil
}

func (mock *TorrentClientMock) Torrents() []torrent.Download {
	return mock.torrents
}

type MargeletMock struct {
	messages []tgbotapi.Chattable
}

func (mock *MargeletMock) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	mock.messages = append(mock.messages, c)
	return tgbotapi.Message{}, nil
}

func (mock *MargeletMock) QuickSend(chatID int, message string) (tgbotapi.Message, error) {
	mock.messages = append(mock.messages, tgbotapi.NewMessage(chatID, message))
	return tgbotapi.Message{}, nil
}

func (mock MargeletMock) QuickReply(chatID, messageID int, message string) (tgbotapi.Message, error) {
	mock.messages = append(mock.messages, tgbotapi.NewMessage(chatID, message))
	return tgbotapi.Message{}, nil
}

func (mock MargeletMock) GetFileDirectURL(fileID string) (string, error) {
	return "example.com", nil
}

func (mock MargeletMock) IsMessageToMe(message tgbotapi.Message) bool {
	return true
}

func (mock MargeletMock) GetConfigRepository() *margelet.ChatConfigRepository {
	return nil
}

func (mock MargeletMock) GetSessionRepository() *margelet.SessionRepository {
	return nil
}

func (mock MargeletMock) GetRedis() *redis.Client {
	return nil
}

func (mock MargeletMock) HandleSession(message tgbotapi.Message, handler margelet.SessionHandler) {

}

func newMargeletMock() *MargeletMock {
	return &MargeletMock{}
}

type DownloadMock struct {
	infoHash       torrent.InfoHash
	info           *metainfo.Info
	bytesCompleted int64
	seeding        bool
	client         *torrent.Client
}

func (mock DownloadMock) InfoHash() torrent.InfoHash {
	return mock.infoHash
}

func (mock DownloadMock) GotInfo() <-chan struct{} {
	return make(<-chan struct{})
}

func (mock DownloadMock) Info() *metainfo.Info {
	return mock.info
}

func (mock DownloadMock) NewReader() (ret *torrent.Reader) {
	return nil
}

func (mock DownloadMock) PieceStateRuns() []torrent.PieceStateRun {
	return []torrent.PieceStateRun{}
}

func (mock DownloadMock) NumPieces() int {
	return 100
}

func (mock DownloadMock) Drop() {

}

func (mock DownloadMock) BytesCompleted() int64 {
	return mock.bytesCompleted
}

func (mock DownloadMock) SubscribePieceStateChanges() *pubsub.Subscription {
	return nil
}

func (mock DownloadMock) Seeding() bool {
	return mock.seeding
}

func (mock DownloadMock) SetDisplayName(dn string) {

}

func (mock DownloadMock) Client() *torrent.Client {
	return mock.client
}

func (mock DownloadMock) AddPeers(pp []torrent.Peer) error {
	return nil
}

func (mock DownloadMock) DownloadAll() {

}

func (mock DownloadMock) Trackers() [][]tracker.Client {
	return [][]tracker.Client{}
}

func (mock DownloadMock) Files() (ret []torrent.File) {
	return []torrent.File{}
}

func (mock DownloadMock) Peers() map[torrent.PeersKey]torrent.Peer {
	return map[torrent.PeersKey]torrent.Peer{}
}

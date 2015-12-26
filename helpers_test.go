package main

import (
	"github.com/Syfaro/telegram-bot-api"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/zhulik/margelet"
	"gopkg.in/redis.v3"
)

type TorrentClientMock struct {
	torrents []torrent.Torrent
}

func newTorrentClientMock() *TorrentClientMock {
	return &TorrentClientMock{}
}

func (mock TorrentClientMock) AddMagnet(string) (torrent.Torrent, error) {
	return mock.torrents[0], nil
}

func (mock TorrentClientMock) AddTorrent(info *metainfo.MetaInfo) (torrent.Torrent, error) {
	return mock.torrents[0], nil
}

func (mock *TorrentClientMock) Torrents() []torrent.Torrent {
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

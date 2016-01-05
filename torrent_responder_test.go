package main

import (
	//	"github.com/anacrolix/torrent/metainfo"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/redis.v3"
	"testing"
)

func TestTorrentResponder(t *testing.T) {
	Convey("With torrent responder", t, func() {
		redis := redis.NewClient(&redis.Options{
			Addr:     "127.0.0.1:6379",
			Password: "",
			DB:       8,
		})
		defer redis.FlushDb()
		bot := newMargeletMock()
		repo := newTorrentsRepository("torrent_bot_torrents", redis)
		client := newTorrentClientMock()
		downloadHandler, _ := newDownloadHandler("test", client, repo)
		responder, _ := newTorrentResponder("test", client, repo, downloadHandler, downloadMock)

		Convey("calling Response", func() {
			Convey("with unauthorized user's message", func() {
				from := tgbotapi.User{UserName: "another"}
				msg := tgbotapi.Message{From: from}
				responder.Response(bot, msg)

				Convey("sent message should contains authorization error", func() {
					So(bot.messages[0].(tgbotapi.MessageConfig).Text, ShouldEqual, "Sorry, you are not allowed to control me!")
				})

				Convey("only one message should be sent", func() {
					So(len(bot.messages), ShouldEqual, 1)
				})
			})

			Convey("with authorized user", func() {
				from := tgbotapi.User{UserName: "test"}

				Convey("with no files", func() {
					msg := tgbotapi.Message{From: from, Text: "test"}
					responder.Response(bot, msg)
					Convey("no messages should sent", func() {
						So(bot.messages, ShouldBeEmpty)
					})
				})

				Convey("with wrong file type", func() {
					document := tgbotapi.Document{FileID: "test", MimeType: "test/test"}
					msg := tgbotapi.Message{From: from, Text: "test", Document: document}
					responder.Response(bot, msg)
					Convey("no messages should sent", func() {
						So(bot.messages, ShouldBeEmpty)
					})
				})

				Convey("with torrent file", func() {
					document := tgbotapi.Document{FileID: "test", MimeType: "application/x-bittorrent"}
					msg := tgbotapi.Message{From: from, Text: "test", Document: document}
					responder.Response(bot, msg)

					Convey("should send message with torrent info", func() {
						So(bot.messages[1].(tgbotapi.MessageConfig).Text, ShouldEqual, "Name: bbb_sunflower_1080p_60fps_stereo_abl.mp4, Size: 515 MB")
					})

					Convey("2 message should be sent", func() {
						So(len(bot.messages), ShouldEqual, 2)
					})
				})
			})
		})
	})
}

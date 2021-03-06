package main

import (
	"github.com/anacrolix/torrent/metainfo"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/redis.v3"
	"testing"
)

func TestMagnetResponder(t *testing.T) {
	Convey("With magnet responder", t, func() {
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
		responder, _ := newMagnetResponder("test", client, repo, downloadHandler)

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

				Convey("with non-magnet text", func() {
					msg := tgbotapi.Message{From: from, Text: "test"}
					responder.Response(bot, msg)
					Convey("no messages should sent", func() {
						So(bot.messages, ShouldBeEmpty)
					})
				})

				Convey("with magnet test", func() {
					torr := &DownloadMock{info: &metainfo.Info{Length: 500, Name: "test"}, bytesCompleted: 100}
					torr.infoChan = make(chan struct{})
					go func() {
						torr.infoChan <- struct{}{}
					}()
					client.torrents = append(client.torrents, torr)

					msg := tgbotapi.Message{From: from, Text: "magnet:?xt=urn:btih:674D163D2184353CE21F3DE5196B0A6D7C2F9FC2&dn=bbb_sunflower_1080p_60fps_stereo_abl.mp4&tr=udp%3a%2f%2ftracker.openbittorrent.com%3a80%2fannounce&tr=udp%3a%2f%2ftracker.publicbt.com%3a80%2fannounce&ws=http%3a%2f%2fdistribution.bbb3d.renderfarming.net%2fvideo%2fmp4%2fbbb_sunflower_1080p_60fps_stereo_abl.mp4"}
					responder.Response(bot, msg)

					Convey("sent message should containse torrent info", func() {
						So(bot.messages[1].(tgbotapi.MessageConfig).Text, ShouldEqual, "Name: test, Size: 500 B")
					})

					Convey("2 messages should sent", func() {
						So(len(bot.messages), ShouldEqual, 2)
					})
				})

			})
		})
	})
}

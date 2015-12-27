package main

import (
	"github.com/Syfaro/telegram-bot-api"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"fmt"
)

func TestStatusHandler(t *testing.T) {
	Convey("When status responder", t, func() {
		bot := newMargeletMock()
		client := newTorrentClientMock()
		handler := newStatusHandler("test", "~/", client, findTorrentByMessage)

		Convey("when calling HelpMessage", func() {
			message := handler.HelpMessage()
			Convey("returning value not empty", func() {
				So(message, ShouldNotBeEmpty)
			})
		})

		Convey("when calling Response", func() {
			Convey("with unauthorized user's message", func() {
				from := tgbotapi.User{UserName: "another"}
				msg := tgbotapi.Message{From: from}
				handler.Response(bot, msg)

				Convey("sent message should contains authorization error", func() {
					So(bot.messages[0].(tgbotapi.MessageConfig).Text, ShouldEqual, "Sorry, you are not allowed to control me!")
				})

				Convey("only one message should be sent", func() {
					So(len(bot.messages), ShouldEqual, 1)
				})
			})

			Convey("with authorized user", func() {
				from := tgbotapi.User{UserName: "test"}
				msg := tgbotapi.Message{From: from}

				Convey("without downloads", func() {
					handler.Response(bot, msg)

					Convey("Sent message should countains information about empty download list", func() {
						So(bot.messages[0].(tgbotapi.MessageConfig).Text, ShouldEqual, "There is no downloads")
					})

					Convey("only one message should be sent", func() {
						So(len(bot.messages), ShouldEqual, 1)
					})
				})

				Convey("with existing downloads", func() {
					torr := DownloadMock{info: &metainfo.Info{Length: 500, Name: "test"}, bytesCompleted: 100}
					client.torrents = append(client.torrents, torr)
					handler.finder = func(client torrentClient, message tgbotapi.Message) (torrent.Download, error) {
						return torr, nil
					}

					Convey("with reply", func() {
						msg.ReplyToMessage = &tgbotapi.Message{}
						handler.Response(bot, msg)

						Convey("sent info about download", func() {
							So(bot.messages[0].(tgbotapi.MessageConfig).Text, ShouldEqual, `0000000000000000000000000000000000000000
Name: test
Size: 500 B
Progress: 20.00%
Seeding: false
Peers: 0
Location: ~/test`)
						})

						Convey("only one message should be sent", func() {
							So(len(bot.messages), ShouldEqual, 1)
						})
					})

					Convey("with hash argument", func() {
						msg.Text = "0000000000000000000000000000000000000000"
						handler.Response(bot, msg)

						Convey("sent info about download", func() {
							So(bot.messages[0].(tgbotapi.MessageConfig).Text, ShouldEqual, `0000000000000000000000000000000000000000
Name: test
Size: 500 B
Progress: 20.00%
Seeding: false
Peers: 0
Location: ~/test`)
						})

						Convey("only one message should be sent", func() {
							So(len(bot.messages), ShouldEqual, 1)
						})
					})

					Convey("without arguments", func() {
						torr1 := DownloadMock{info: &metainfo.Info{Length: 1000, Name: "test again"}, bytesCompleted: 100}
						client.torrents = append(client.torrents, torr1)
						handler.finder = func(client torrentClient, message tgbotapi.Message) (torrent.Download, error) {
							return torr, fmt.Errorf("")
						}
						handler.Response(bot, msg)

						Convey("sent info about all all downloads", func() {
							So(bot.messages[0].(tgbotapi.MessageConfig).Text, ShouldEqual, `0000000000000000000000000000000000000000
Name: test, Size: 500 B, Progress: 20.00%`)

							So(bot.messages[1].(tgbotapi.MessageConfig).Text, ShouldEqual, `0000000000000000000000000000000000000000
Name: test again, Size: 1.0 kB, Progress: 10.00%`)
						})

						Convey("only one message should be sent", func() {
							So(len(bot.messages), ShouldEqual, 2)
						})
					})
				})
			})
		})
	})
}

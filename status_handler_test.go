package main

import (
	"github.com/anacrolix/torrent/metainfo"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestStatusHandler(t *testing.T) {
	Convey("With status responder", t, func() {
		bot := newMargeletMock()
		client := newTorrentClientMock()
		handler := newStatusHandler("test", "~/", client)

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

					Convey("sent message should countains information about empty download list", func() {
						So(bot.messages[0].(tgbotapi.MessageConfig).Text, ShouldEqual, "There is no downloads")
					})

					Convey("only one message should be sent", func() {
						So(len(bot.messages), ShouldEqual, 1)
					})
				})

				Convey("with existing downloads", func() {
					torr := DownloadMock{info: &metainfo.Info{Length: 500, Name: "test"}, bytesCompleted: 100}
					client.torrents = append(client.torrents, torr)

					Convey("with reply", func() {
						Convey("with existing hasn", func() {
							msg.ReplyToMessage = &tgbotapi.Message{Text: "0000000000000000000000000000000000000000\nTest"}
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

						Convey("with unknown hash", func() {
							msg.ReplyToMessage = &tgbotapi.Message{Text: "test\nTest"}

							handler.Response(bot, msg)

							Convey("sent message should contains information abount unknown download", func() {
								So(bot.messages[0].(tgbotapi.MessageConfig).Text, ShouldEqual, "Cannot find download with hash test")
							})

							Convey("only one message should be sent", func() {
								So(len(bot.messages), ShouldEqual, 1)
							})
						})
					})

					Convey("with hash argument", func() {
						msg.Text = "/delete 0000000000000000000000000000000000000000"
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

					Convey("with unknown hash angument", func() {
						msg.Text = "/delete test"
						handler.Response(bot, msg)

						Convey("sent message that download not found", func() {
							So(bot.messages[0].(tgbotapi.MessageConfig).Text, ShouldEqual, "Cannot find download with hash test")
						})

						Convey("only one message should be sent", func() {
							So(len(bot.messages), ShouldEqual, 1)
						})
					})
				})
			})
		})
	})
}

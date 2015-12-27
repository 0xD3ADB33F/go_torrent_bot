package main

import (
	"github.com/Syfaro/telegram-bot-api"
	"github.com/anacrolix/torrent/metainfo"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestDeleteHandler(t *testing.T) {
	Convey("With status responder", t, func() {
		bot := newMargeletMock()
		client := newTorrentClientMock()
		handler := newDeleteHandler("test", "/tmp/", client)

		Convey("when calling HelpMessage", func() {
			message := handler.HelpMessage()
			Convey("returning value not empty", func() {
				So(message, ShouldNotBeEmpty)
			})
		})

		Convey("when calling handleDeleteCommand", func() {
			Convey("with unauthorized user's message", func() {
				from := tgbotapi.User{UserName: "another"}
				msg := tgbotapi.Message{From: from}
				handler.HandleResponse(bot, msg, []tgbotapi.Message{})

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

				Convey("without passing hash", func() {
					handler.HandleResponse(bot, msg, []tgbotapi.Message{})

					Convey("sent message shoult contains usage info", func() {
						So(bot.messages[0].(tgbotapi.MessageConfig).Text, ShouldEqual, "usage: /delete <download hash>")
					})
				})

				Convey("with reply", func() {
					torr := DownloadMock{info: &metainfo.Info{Length: 500, Name: "test"}, bytesCompleted: 100}
					client.torrents = append(client.torrents, torr)
					Convey("with existing hash", func() {
						msg.ReplyToMessage = &tgbotapi.Message{Text: "0000000000000000000000000000000000000000\nTest"}
						handler.HandleResponse(bot, msg, []tgbotapi.Message{})

						Convey("sent message should contains delete question", func() {
							So(bot.messages[0].(tgbotapi.MessageConfig).Text, ShouldEqual, `You trying to delete Name: test, Size: 500 B.
Would you like to remove downloaded files?`)
						})

						Convey("only one message should be sent", func() {
							So(len(bot.messages), ShouldEqual, 1)
						})
					})

					Convey("with unknown hash", func() {
						msg.ReplyToMessage = &tgbotapi.Message{Text: "test\nTest"}

						handler.HandleResponse(bot, msg, []tgbotapi.Message{})

						Convey("sent message should contains information abount unknown download", func() {
							So(bot.messages[0].(tgbotapi.MessageConfig).Text, ShouldEqual, "Cannot find download with hash test")
						})

						Convey("only one message should be sent", func() {
							So(len(bot.messages), ShouldEqual, 1)
						})
					})
				})

				Convey("with hash argument", func() {
					torr := DownloadMock{info: &metainfo.Info{Length: 500, Name: "test"}, bytesCompleted: 100}
					client.torrents = append(client.torrents, torr)

					Convey("without hash argument", func() {
						msg.Text = "/delete"
						handler.HandleResponse(bot, msg, []tgbotapi.Message{})
						Convey("sent message shoult contains usage info", func() {
							So(bot.messages[0].(tgbotapi.MessageConfig).Text, ShouldEqual, "usage: /delete <download hash>")
						})

						Convey("only one message should be sent", func() {
							So(len(bot.messages), ShouldEqual, 1)
						})
					})

					Convey("with existing hash", func() {
						msg.Text = "0000000000000000000000000000000000000000\nTest"
						handler.HandleResponse(bot, msg, []tgbotapi.Message{})

						Convey("sent message should contains delete question", func() {
							So(bot.messages[0].(tgbotapi.MessageConfig).Text, ShouldEqual, `You trying to delete Name: test, Size: 500 B.
Would you like to remove downloaded files?`)
						})

						Convey("only one message should be sent", func() {
							So(len(bot.messages), ShouldEqual, 1)
						})
					})

					Convey("with unknown hash", func() {
						msg.Text = "test\nTest"

						handler.HandleResponse(bot, msg, []tgbotapi.Message{})

						Convey("sent message should contains information abount unknown download", func() {
							So(bot.messages[0].(tgbotapi.MessageConfig).Text, ShouldEqual, "Cannot find download with hash test")
						})

						Convey("only one message should be sent", func() {
							So(len(bot.messages), ShouldEqual, 1)
						})
					})
				})
				Convey("with answer on delete question", func() {
					torr := DownloadMock{info: &metainfo.Info{Length: 500, Name: "test"}, bytesCompleted: 100}
					client.torrents = append(client.torrents, torr)
					Convey("with reply in previous message", func() {
						prevMessage := tgbotapi.Message{}
						prevMessage.ReplyToMessage = &tgbotapi.Message{Text: "0000000000000000000000000000000000000000\nTest"}

						Convey("with unknown unswer", func() {
							msg := tgbotapi.Message{Text: "test", From: from}
							handler.HandleResponse(bot, msg, []tgbotapi.Message{prevMessage})

							Convey("sent message should contains information about bot's misunderstanding", func() {
								So(bot.messages[0].(tgbotapi.MessageConfig).Text, ShouldEqual, `Sorry, i don't understand.`)
							})

							Convey("only one message should be sent", func() {
								So(len(bot.messages), ShouldEqual, 1)
							})
						})

						Convey("with cancel unswer", func() {
							msg := tgbotapi.Message{Text: "cancel", From: from}
							handler.HandleResponse(bot, msg, []tgbotapi.Message{prevMessage})

							Convey("sent message should contains information about canceling", func() {
								So(bot.messages[0].(tgbotapi.MessageConfig).Text, ShouldEqual, `Delete canceled!`)
							})

							Convey("only one message should be sent", func() {
								So(len(bot.messages), ShouldEqual, 1)
							})
						})

						Convey("with no unswer", func() {
							msg := tgbotapi.Message{Text: "no", From: from}
							handler.HandleResponse(bot, msg, []tgbotapi.Message{prevMessage})

							Convey("sent message should contains information abount removing torrent from download list", func() {
								So(bot.messages[0].(tgbotapi.MessageConfig).Text, ShouldEqual, `Downloading of Name: test, Size: 500 B is canceled!`)
							})

							Convey("only one message should be sent", func() {
								So(len(bot.messages), ShouldEqual, 1)
							})
						})

						Convey("with yes unswer", func() {
							msg := tgbotapi.Message{Text: "yes", From: from}
							handler.HandleResponse(bot, msg, []tgbotapi.Message{prevMessage})
							time.Sleep(100 * time.Millisecond)

							Convey("sent message should contains information abount removing torrent from download list", func() {
								So(bot.messages[0].(tgbotapi.MessageConfig).Text, ShouldEqual, `Downloading of Name: test, Size: 500 B is canceled, files removed!`)
							})

							Convey("only one message should be sent", func() {
								So(len(bot.messages), ShouldEqual, 1)
							})
						})
					})

				})
			})
		})
	})
}

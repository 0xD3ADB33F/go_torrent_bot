package main

import (
	"github.com/anacrolix/torrent/metainfo"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/redis.v3"
	"io/ioutil"
	"testing"
)

func TestDownloadHandler(t *testing.T) {
	Convey("With download handler", t, func() {
		bot := newMargeletMock()
		redis := redis.NewClient(&redis.Options{
			Addr:     "127.0.0.1:6379",
			Password: "",
			DB:       8,
		})
		defer redis.FlushDb()

		repo := newTorrentsRepository("torrent_bot_torrents", redis)
		client := newTorrentClientMock()
		handler, _ := newDownloadHandler("test", client, repo)

		Convey("when calling HelpMessage()", func() {
			result := handler.HelpMessage()
			Convey("response should not be empty", func() {
				So(result, ShouldNotBeEmpty)
			})
		})

		Convey("When calling HandleResponse()", func() {
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

			Convey("with authorized user's message", func() {
				chat := tgbotapi.Chat{ID: 100}
				from := tgbotapi.User{UserName: "test", ID: 500}

				Convey("with no existing torrents in repo", func() {
					msg := tgbotapi.Message{From: from, Chat: chat}
					handler.HandleResponse(bot, msg, []tgbotapi.Message{})

					Convey("sent help message", func() {
						So(bot.messages[0].(tgbotapi.MessageConfig).Text, ShouldEqual, "Please send me torrent file, do not use /download directly")
					})

					Convey("only one message should be sent", func() {
						So(len(bot.messages), ShouldEqual, 1)
					})
				})

				Convey("with magnet in repo", func() {
					repo.Add(100, 500, []byte("magnet:?xt=urn:btih:674D163D2184353CE21F3DE5196B0A6D7C2F9FC2&dn=bbb_sunflower_1080p_60fps_stereo_abl.mp4&tr=udp%3a%2f%2ftracker.openbittorrent.com%3a80%2fannounce&tr=udp%3a%2f%2ftracker.publicbt.com%3a80%2fannounce&ws=http%3a%2f%2fdistribution.bbb3d.renderfarming.net%2fvideo%2fmp4%2fbbb_sunflower_1080p_60fps_stereo_abl.mp4"))
					Convey("without responses", func() {
						msg := tgbotapi.Message{From: from, Chat: chat}
						handler.HandleResponse(bot, msg, []tgbotapi.Message{})

						Convey("sent help message", func() {
							So(bot.messages[0].(tgbotapi.MessageConfig).Text, ShouldEqual, "Would you like to download it?")
						})

						Convey("only one message should be sent", func() {
							So(len(bot.messages), ShouldEqual, 1)
						})
					})

					Convey("with responses", func() {
						Convey("with unknown answer", func() {
							msg := tgbotapi.Message{From: from, Chat: chat, Text: "test"}
							handler.HandleResponse(bot, msg, []tgbotapi.Message{msg})

							Convey("sent error message", func() {
								So(bot.messages[0].(tgbotapi.MessageConfig).Text, ShouldEqual, "Sorry, i don't understand.")
							})

							Convey("only one message should be sent", func() {
								So(len(bot.messages), ShouldEqual, 1)
							})
						})

						Convey("with yes answer", func() {
							msg := tgbotapi.Message{From: from, Chat: chat, Text: "yes"}

							Convey("with existing torrent", func() {
								torr := &DownloadMock{info: &metainfo.Info{Length: 500, Name: "test"}, bytesCompleted: 100}
								client.torrents = append(client.torrents, torr)
								handler.HandleResponse(bot, msg, []tgbotapi.Message{msg})

								Convey("sent downloading message", func() {
									So(bot.messages[0].(tgbotapi.MessageConfig).Text, ShouldEqual, "test is downloading...")
								})

								Convey("only one message should be sent", func() {
									So(len(bot.messages), ShouldEqual, 1)
								})
							})

							Convey("with missing torrent", func() {
								repo.Delete(100, 500)
								handler.HandleResponse(bot, msg, []tgbotapi.Message{msg})

								Convey("sent error message", func() {
									So(bot.messages[0].(tgbotapi.MessageConfig).Text, ShouldEqual, "Sorry, something is wrong. Please, try again")
								})

								Convey("only one message should be sent", func() {
									So(len(bot.messages), ShouldEqual, 1)
								})
							})
						})

						Convey("with no answer", func() {
							msg := tgbotapi.Message{From: from, Chat: chat, Text: "no"}
							handler.HandleResponse(bot, msg, []tgbotapi.Message{msg})

							Convey("sent downloading message", func() {
								So(bot.messages[0].(tgbotapi.MessageConfig).Text, ShouldEqual, "Ok, i will forgive it")
							})

							Convey("only one message should be sent", func() {
								So(len(bot.messages), ShouldEqual, 1)
							})
						})

					})
				})

				Convey("with torrent file in repo", func() {
					data, _ := ioutil.ReadFile("testdata/bbb.torrent")
					repo.Add(100, 500, data)
					Convey("without responses", func() {
						msg := tgbotapi.Message{From: from, Chat: chat}
						handler.HandleResponse(bot, msg, []tgbotapi.Message{})

						Convey("sent help message", func() {
							So(bot.messages[0].(tgbotapi.MessageConfig).Text, ShouldEqual, "Would you like to download it?")
						})

						Convey("only one message should be sent", func() {
							So(len(bot.messages), ShouldEqual, 1)
						})
					})

					Convey("with responses", func() {
						Convey("with unknown answer", func() {
							msg := tgbotapi.Message{From: from, Chat: chat, Text: "test"}
							handler.HandleResponse(bot, msg, []tgbotapi.Message{msg})

							Convey("sent error message", func() {
								So(bot.messages[0].(tgbotapi.MessageConfig).Text, ShouldEqual, "Sorry, i don't understand.")
							})

							Convey("only one message should be sent", func() {
								So(len(bot.messages), ShouldEqual, 1)
							})
						})

						Convey("with yes answer", func() {
							msg := tgbotapi.Message{From: from, Chat: chat, Text: "yes"}

							Convey("with existing torrent", func() {
								torr := &DownloadMock{info: &metainfo.Info{Length: 500, Name: "test"}, bytesCompleted: 100}
								client.torrents = append(client.torrents, torr)
								handler.HandleResponse(bot, msg, []tgbotapi.Message{msg})

								Convey("sent downloading message", func() {
									So(bot.messages[0].(tgbotapi.MessageConfig).Text, ShouldEqual, "test is downloading...")
								})

								Convey("only one message should be sent", func() {
									So(len(bot.messages), ShouldEqual, 1)
								})
							})

							Convey("with missing torrent", func() {
								repo.Delete(100, 500)
								handler.HandleResponse(bot, msg, []tgbotapi.Message{msg})

								Convey("sent error message", func() {
									So(bot.messages[0].(tgbotapi.MessageConfig).Text, ShouldEqual, "Sorry, something is wrong. Please, try again")
								})

								Convey("only one message should be sent", func() {
									So(len(bot.messages), ShouldEqual, 1)
								})
							})
						})

						Convey("with no answer", func() {
							msg := tgbotapi.Message{From: from, Chat: chat, Text: "no"}
							handler.HandleResponse(bot, msg, []tgbotapi.Message{msg})

							Convey("sent downloading message", func() {
								So(bot.messages[0].(tgbotapi.MessageConfig).Text, ShouldEqual, "Ok, i will forgive it")
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

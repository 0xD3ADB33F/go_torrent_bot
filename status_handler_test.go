package main

import (
	"github.com/Syfaro/telegram-bot-api"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"github.com/anacrolix/torrent"
)

func TestStatusHandler(t *testing.T) {
	Convey("When status responder", t, func() {
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

			Convey("with authorized user", func(){
				from := tgbotapi.User{UserName: "test"}
				msg := tgbotapi.Message{From: from}
				handler.Response(bot, msg)

				Convey("without downloads", func(){
					Convey("Sent message should countains information about empty download list", func() {
						So(bot.messages[0].(tgbotapi.MessageConfig).Text, ShouldEqual, "There is no downloads")
					})

					Convey("only one message should be sent", func() {
						So(len(bot.messages), ShouldEqual, 1)
					})
				})

				Convey("with existing downloads", func() {
					t := torrent.Torrent{}
					client.torrents = append(client.torrents, t)
				})
			})
		})
	})
}

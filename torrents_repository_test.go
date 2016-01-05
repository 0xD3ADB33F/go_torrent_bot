package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/redis.v3"
	"testing"
)

func TestTorrentRepository(t *testing.T) {
	Convey("With torrents repository", t, func() {
		redis := redis.NewClient(&redis.Options{
			Addr:     "127.0.0.1:6379",
			Password: "",
			DB:       8,
		})
		defer redis.FlushDb()

		repo := newTorrentsRepository("torrent_bot_torrents", redis)

		Convey("after calling Add", func() {
			repo.Add(100, 500, []byte("test"))

			Convey("Exists for this chat and user", func() {
				result := repo.Exists(100, 500)
				Convey("should return true", func() {
					So(result, ShouldBeTrue)
				})
			})

			Convey("Exists for another chat and user", func() {
				result := repo.Exists(100, 501)
				Convey("should return false", func() {
					So(result, ShouldBeFalse)
				})
			})

			Convey("Get for this chat and user", func() {
				result, _ := repo.Get(100, 500)
				Convey("should return stored content", func() {
					So(result, ShouldResemble, []byte("test"))
				})
			})

			Convey("Get for another chat and user", func() {
				_, err := repo.Get(100, 501)
				Convey("should return error", func() {
					So(err, ShouldNotBeNil)
				})
			})

			Convey("Delete for this chat and user", func() {
				repo.Delete(100, 500)
				result := repo.Exists(100, 500)
				Convey("should deletes stored data", func() {
					So(result, ShouldBeFalse)
				})
			})

			Convey("Delete for another chat and user", func() {
				repo.Delete(100, 501)
				result := repo.Exists(100, 500)
				Convey("should not delete stored data", func() {
					So(result, ShouldBeTrue)
				})
			})
		})
	})
}

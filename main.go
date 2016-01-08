package main

import (
	"fmt"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/data/mmap"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/yvasiyarov/gorelic"
	"github.com/zhulik/margelet"
	"gopkg.in/alecthomas/kingpin.v2"
	"syscall"
)

// Not portable, only Linux and BSD. Windows sucks.
func checkDownloadPath(path string) {
	if syscall.Access(path, 2) != nil {
		panic(fmt.Errorf("%s is not writeable or not exists", path))
	}
}

func main() {
	token := kingpin.Flag("token", "Telegram Bot token").Required().Short('t').String()
	redisURL := kingpin.Flag("redis_url", "Redis url").Short('r').Default("127.0.0.1:6379").String()
	redisPassword := kingpin.Flag("redis_password", "Redis password").Default("").Short('p').String()
	redisDB := kingpin.Flag("redis_db", "Redis password").Default("0").Short('d').Int64()
	downloadPath := kingpin.Flag("path", "Download path").Required().Short('o').String()
	authorizedUsername := kingpin.Flag("username", "Username of user, thich can control bot").Required().Short('u').String()
	newrelicLicence := kingpin.Flag("newrelic", "Newrelic licence for monitoring").Short('n').String()
	kingpin.Parse()

	checkDownloadPath(*downloadPath)

	bot, err := margelet.NewMargelet("torrent_bot", *redisURL, *redisPassword, *redisDB, *token, false)

	if err != nil {
		panic(err)
	}

	config := torrent.Config{
		TorrentDataOpener: func(info *metainfo.Info) torrent.Data {
			ret, _ := mmap.TorrentData(info, *downloadPath)
			return ret
		},
	}

	tClient, err := torrent.NewClient(&config)
	if err != nil {
		panic(err)
	}
	defer tClient.Close()

	client := &torrentClientProxy{tClient}

	repo := newTorrentsRepository("torrent_bot_torrents", bot.GetRedis())

	downHandler, _ := newDownloadHandler(*authorizedUsername, client, repo)

	torrentResponder, err := newTorrentResponder(*authorizedUsername, client, repo, downHandler, download)
	if err != nil {
		panic(err)
	}

	magnetResponder, err := newMagnetResponder(*authorizedUsername, client, repo, downHandler)
	if err != nil {
		panic(err)
	}

	bot.AddMessageResponder(torrentResponder)
	bot.AddMessageResponder(magnetResponder)
	bot.AddSessionHandler("/download", downHandler)
	bot.AddCommandHandler("/status", newStatusHandler(*authorizedUsername, *downloadPath, client))
	bot.AddSessionHandler("/delete", newDeleteHandler(*authorizedUsername, *downloadPath, client))

	if len(*newrelicLicence) > 0 {
		agent := gorelic.NewAgent()
		agent.NewrelicName = "go_torrent_bot"
		agent.Verbose = false
		agent.NewrelicLicense = *newrelicLicence
		err := agent.Run()
		if err != nil {
			panic(err)
		}
	}

	bot.Run()
}

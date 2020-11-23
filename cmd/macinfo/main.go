package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/jgengo/macinfo-client/internal/config"
	"github.com/jgengo/macinfo-client/internal/gatherer"
	"github.com/jgengo/macinfo-client/internal/sender"
	"github.com/jgengo/macinfo-client/internal/utils"

	"github.com/getsentry/sentry-go"
)

const appVersion = "0.8"

func doEvery(d time.Duration) {
	localHub := sentry.CurrentHub().Clone()
	for range time.Tick(d) {

		if err := config.ConnectOSQ(); err != nil {
			log.Printf("osquery (error) while creating a new client: %v\n", err)
			localHub.CaptureException(err)
			continue
		}
		system, err := gatherer.GetInfo()
		if err != nil {
			log.Printf("gatherer (error) getinfo: %v\n", err)
			localHub.CaptureException(err)
			continue
		}
		if err := sender.Process(system); err != nil {
			log.Printf("sender (error) process: %v\n", err)
			localHub.CaptureException(err)
			continue
		}
		utils.OsQ.Client.Close()
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cfgPtr := flag.String("cfg", "/var/macinfo/macinfo.yml", "specify another config path")
	versionPtr := flag.Bool("version", false, "display app version")
	flag.Parse()

	if *versionPtr {
		fmt.Println("MacInfo version", appVersion)
		return
	}

	if err := config.Initiate(*cfgPtr); err != nil {
		log.Fatalf("error while initializing the config file: %v\n", err)
	}
	defer utils.OsQ.Client.Close()

	err := sentry.Init(sentry.ClientOptions{
		Dsn:     utils.Cfg.SentryDSN,
		Debug:   true,
		Release: appVersion,
	})
	if err != nil {
		log.Fatalf("sentry (error) failed to init sentry: %v\n", err)
	}
	defer sentry.Flush(2 * time.Second)

	doEvery(utils.Cfg.SyncInterval * time.Second)

	for {

	}
}

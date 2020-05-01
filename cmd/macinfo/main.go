package main

import (
	"flag"
	"log"
	"time"

	"github.com/jgengo/macinfo-client/internal/config"
	"github.com/jgengo/macinfo-client/internal/gatherer"
	"github.com/jgengo/macinfo-client/internal/sender"
	"github.com/jgengo/macinfo-client/internal/utils"
)

func doEvery(d time.Duration) {
	for range time.Tick(d) {
		system := gatherer.GetInfo()
		sender.Process(system)
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cfgPtr := flag.String("cfg", "/var/macinfo/macinfo.yml", "specify another config path")
	flag.Parse()

	config.Initiate(*cfgPtr)
	defer utils.OsQ.Client.Close()

	doEvery(utils.Cfg.SyncInterval * time.Minute)
	for {

	}
}

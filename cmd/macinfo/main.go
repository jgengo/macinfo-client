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
)

const appVersion = "0.4"

func doEvery(d time.Duration) {
	for range time.Tick(d) {
		system := gatherer.GetInfo()
		sender.Process(system)
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

	config.Initiate(*cfgPtr)
	defer utils.OsQ.Client.Close()

	doEvery(utils.Cfg.SyncInterval * time.Minute)
	for {

	}
}

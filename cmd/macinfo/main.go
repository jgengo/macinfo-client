package main

import (
	"flag"

	"github.com/jgengo/macinfo-client/internal/config"
	"github.com/jgengo/macinfo-client/internal/gatherer"
	"github.com/jgengo/macinfo-client/internal/sender"
	"github.com/jgengo/macinfo-client/internal/utils"
)

func main() {
	cfgPtr := flag.String("cfg", "/etc/macinfo.yml", "specify another config path")
	flag.Parse()

	config.Initiate(*cfgPtr)

	system := gatherer.GetInfo()
	sender.Process(system)

	utils.OsQ.Client.Close()
}

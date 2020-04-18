package main

import (
	"flag"

	"github.com/jgengo/macinfo-client/internal/config"
	"github.com/jgengo/macinfo-client/internal/gatherer"
	"github.com/jgengo/macinfo-client/internal/sender"
)

func main() {
	cfgPtr := flag.String("cfg", "/etc/macinfo.yml", "specify another config path")
	flag.Parse()

	config.Initiate(*cfgPtr)

	system := gatherer.GetInfo()
	sender.Process(system)
}

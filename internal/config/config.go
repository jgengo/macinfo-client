package config

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"time"

	"github.com/jgengo/macinfo-client/internal/utils"
	"github.com/kolide/osquery-go"
	"gopkg.in/yaml.v2"
)

func readConfig(cfgPath string) ([]byte, error) {
	fd, err := filepath.Abs(cfgPath)
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadFile(fd)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func ConnectOSQ() {
	c, err := osquery.NewClient(utils.Cfg.OsqSock, 10*time.Second)
	if err != nil {
		log.Fatalf("osquery (error) while creating a new client: %v\n", err)
	}
	utils.OsQ.Client = c
}

// Initiate populates the global variable Cfg and OsQ with the information in the yml
func Initiate(cfgPath string) {
	content, err := readConfig(cfgPath)
	if err != nil {
		log.Fatalf("error while reading the config file: %v\n", err)
	}
	if err := yaml.Unmarshal(content, &utils.Cfg); err != nil {
		log.Fatalf("cannot unmarshall the config file: %v\n", err)
	}

	utils.Cfg.CfgPath = cfgPath
}

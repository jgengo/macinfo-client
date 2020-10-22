package config

import (
	"io/ioutil"
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

func ConnectOSQ() error {
	c, err := osquery.NewClient(utils.Cfg.OsqSock, 10*time.Second)
	if err != nil {
		return err
	}
	utils.OsQ.Client = c
	return nil
}

// Initiate populates the global variable Cfg and OsQ with the information in the yml
func Initiate(cfgPath string) error {
	content, err := readConfig(cfgPath)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(content, &utils.Cfg); err != nil {
		return err
	}

	utils.Cfg.CfgPath = cfgPath
	return nil
}

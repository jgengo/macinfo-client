package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"time"

	"github.com/kolide/osquery-go"
)

// Cfg is a global variable that stores all the useful Config information
var Cfg Config

// OsQ is a globbal variable that Store OsQuery Client
var OsQ OsQuery

// Config stores the config yaml information
type Config struct {
	OsqSock      string        `yaml:"osquery_sock"`
	APIURL       string        `yaml:"api_url"`
	APIToken     string        `yaml:"api_token"`
	SyncInterval time.Duration `yaml:"sync_interval"`
	CfgPath      string
}

// OsQuery stores osquery information
type OsQuery struct {
	Client *osquery.ExtensionManagerClient
}

// ChangeToken changes the token and save it into the config file
func ChangeToken(token string) {
	Cfg.APIToken = token

	b, err := ioutil.ReadFile(Cfg.CfgPath)
	if err != nil {
		log.Printf("failed to read the config file: %v\n", err)
		return
	}

	bString := string(b)
	r, _ := regexp.Compile("token: (.*)\\n")
	bString = r.ReplaceAllString(bString, fmt.Sprintf("token: %s\n", token))

	if err = ioutil.WriteFile(Cfg.CfgPath, []byte(bString), 0600); err != nil {
		log.Printf("failed to save new config file: %v\n", err)
		return
	}
}

package utils

import (
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

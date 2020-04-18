package utils

import "github.com/kolide/osquery-go"

var Cfg Config
var OsQ OsQuery

// Config stores the config yaml information
type Config struct {
	OsqSock  string `yaml:"osquery_sock"`
	APIURL   string `yaml:"api_url"`
	APIToken string `yaml:"api_token"`
}

// OsQuery stores osquery information
type OsQuery struct {
	Client *osquery.ExtensionManagerClient
}

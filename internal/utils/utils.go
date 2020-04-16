package utils

import "github.com/kolide/osquery-go"

// OsQuery stores osquery information
type OsQuery struct {
	Client *osquery.ExtensionManagerClient
}

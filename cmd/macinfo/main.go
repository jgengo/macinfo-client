package main

import (
	"log"
	"time"

	"github.com/jgengo/macinfo-client/internal/hardware"
	"github.com/jgengo/macinfo-client/internal/utils"
	"github.com/kolide/osquery-go"
)

const (
	osqSock = "/var/osquery/osquery.em"
)

func main() {
	osq := &utils.OsQuery{}
	c, err := osquery.NewClient(osqSock, 10*time.Second)
	if err != nil {
		log.Fatalf("osquery (error) while creating a new client: %v\n", err)
	}
	osq.Client = c
	defer c.Close()

	hardware.GetInfo(osq)

	// resp, err := c.Query(os.Args[1])
	// if err != nil {
	// 	log.Fatalf("Error communicating with osqueryd: %v", err)
	// }
	// if resp.Status.Code != 0 {
	// 	log.Fatalf("osqueryd returned error: %s", resp.Status.Message)
	// }

	// fmt.Printf("Got results:\n%#v\n", resp.Response)
}

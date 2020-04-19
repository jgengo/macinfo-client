package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/jgengo/macinfo-client/internal/gatherer"
	"github.com/jgengo/macinfo-client/internal/utils"
)

// Process is the entrypoint function of the sender package
func Process(s *gatherer.System) {
	json, err := json.Marshal(s)
	if err != nil {
		fmt.Printf("Error while marshaling system struct: %v\n", err)
	}

	url := fmt.Sprintf("%s/sync", utils.Cfg.APIURL)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(json))
	if err != nil {
		log.Fatalf("error while sending the info: %v", err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

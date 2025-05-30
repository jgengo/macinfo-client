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

func sync(s *gatherer.System) (*http.Response, error) {
	sJSON, err := json.Marshal(s)
	if err != nil {
		log.Printf("failed to marshal: %v\n", err)
		return nil, err
	}

	url := fmt.Sprintf("%s/sync", utils.Cfg.APIURL)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(sJSON))
	if err != nil {
		log.Printf("failed to Post: %v", err)
		return nil, err
	}

	return resp, nil
}

// Process is the entrypoint function of the sender package
func Process(s *gatherer.System) error {
	resp, err := sync(s)
	if err != nil {
		log.Printf("failed to sync: %v", err)
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("failed to read body: %v", err)
		return err
	}

	if resp.StatusCode == 201 {
		var respCreated map[string]interface{}
		json.Unmarshal([]byte(body), &respCreated)
		s.Token = respCreated["token"].(string)
		utils.ChangeToken(respCreated["token"].(string))
		_, err := sync(s)
		if err != nil {
			log.Printf("error while re-sync: %v", err)
			return err
		}
	}

	log.Printf("Synced - server response: %s", resp.Status)
	return nil
}

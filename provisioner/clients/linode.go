package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/alancuriel/game-hosting-sass/provisioner/models"
)

const (
	LINODE_URL  = "https://api.linode.com/v4"
)

type Linode struct {
	HttpClient *http.Client
	ApiKey     string
}

func (l *Linode) CreateLinode(req *models.CreateLinodeRequest) (*models.CreateLinodeResponse, error) {
	if l.ApiKey == "" {
		return nil, fmt.Errorf("No api key provided")
	}

	resp, err := l.postJson("/linode/instances", req)

	if err != nil {
		log.Printf("Error sending POST createLinode %s\n", err.Error())
		return nil, fmt.Errorf("Error sending POST request for creating linode server")
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		log.Printf("Status code sending POST createLinode %d\n", resp.StatusCode)
		return nil, fmt.Errorf("Error sending POST request for creating linode server")
	}

	linodeResp, err := models.JsonFromBody[models.CreateLinodeResponse](resp.Body)

	if err != nil {
		return nil, fmt.Errorf("Error parsing response for creating linode server")
	}

	return linodeResp, nil
}

func (l *Linode) postJson(path string, body interface{}) (*http.Response, error) {
	jsonBody, err := json.Marshal(body)

	if err != nil {
		return nil, fmt.Errorf("Error creating json for POST %s", path)
	}

	req, err := http.NewRequest("POST", LINODE_URL+path, bytes.NewReader(jsonBody))

	req.Header.Add("Authorization", "Bearer "+l.ApiKey)
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		return nil, fmt.Errorf("Error creating request for POST %s", path)
	}

	return l.HttpClient.Do(req)
}

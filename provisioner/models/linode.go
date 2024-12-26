package models

import (
	"encoding/json"
	"fmt"
	"io"
)

type CreateLinodeRequest struct {
	Image        string            `json:"image,omitempty"`
	Region       string            `json:"region,omitempty"`
	Label        string            `json:"label,omitempty"`
	InstanceType string            `json:"type,omitempty"`
	PrivateIp    bool              `json:"private_ip"`
	RootPass     string            `json:"root_pass,omitempty"`
	FirewallId   uint32            `json:"firewall_id,omitempty"`
	Metadata     map[string]string `json:"metadata,omitempty"`
}

type CreateLinodeResponse struct {
	Id   int64    `json:"id"`
	Ipv4 []string `json:"ipv4,omitempty"`
}

func JsonFromBody[T any](readCloser io.ReadCloser) (*T, error) {
	if readCloser == nil {
		return nil, fmt.Errorf("Invalid readCloser provided")
	}

	linodeResponse := new(T)
	bytes, err := io.ReadAll(readCloser)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, linodeResponse)

	if err != nil {
		return nil, err
	}
	return linodeResponse, nil
}

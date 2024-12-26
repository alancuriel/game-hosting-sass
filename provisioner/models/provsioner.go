package models

type ProvisionMcServerRequest struct {
	Instance MinecraftInstance `json:"instance,omitempty"`
	Region   Region            `json:"region,omitempty"`
	Username string            `json:"username,omitempty"`
}

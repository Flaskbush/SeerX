package models

import "time"

// Connected Agents
type Agent struct {
	ID       string    `json:"id"`
	IP       string    `json:"ip"`
	LastSeen time.Time `json:"last_seen"`
}

package model

type Action struct {
	Action string            `json:"action"`
	Args   map[string]string `json:"args"`
}

package model

type ContainerStats struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	State    string            `json:"state"`
	LastUsed string            `json:"last_used"`
	Image    string            `json:"image"`
	Labels   map[string]string `json:"labels"`
	Ports    string            `json:"ports"`
	CPU      string            `json:"cpu"`
	Memory   string            `json:"memory"`
	NetIO    string            `json:"net_io"`
	DiskIO   string            `json:"disk_io"`
}

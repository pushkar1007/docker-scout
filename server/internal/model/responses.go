package model

type CreateVolumeRequest struct {
	Name    string            `json:"name"`
	Driver  string            `json:"driver"`
	Labels  map[string]string `json:"labels"`
	Options map[string]string `json:"options"`
}

type BuildImageRequest struct {
	ContextPath string `json:"context_path"`
	Dockerfile  string `json:"dockerfile"`
	Tag         string `json:"tag"`
	NoCache     bool   `json:"no_cache"`
}

type CreateNetworkRequest struct {
	Name    string            `json:"name"`
	Driver  string            `json:"driver"`
	Options map[string]string `json:"options"`
}

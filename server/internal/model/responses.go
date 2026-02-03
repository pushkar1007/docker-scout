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

type CreateContainerRequest struct {
	Name          string            `json:"name"`
	Image         string            `json:"image"`
	Entrypoint    []string          `json:"entrypoint"`
	Cmd           []string          `json:"cmd"`
	Env           []string          `json:"env"`
	Labels        map[string]string `json:"labels"`
	ExposedPorts  map[string]string `json:"exposed_ports"`
	PortBindings  map[string]string `json:"port_bindings"`
	Volumes       []string          `json:"volumes"`
	NetworkMode   string            `json:"network_mode"`
	RestartPolicy string            `json:"restart_policy"`
}

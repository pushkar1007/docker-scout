package model

type VolumeSummary struct {
	Name       string            `json:"name"`
	Driver     string            `json:"driver"`
	Mountpoint string            `json:"mountpoint"`
	Scope      string            `json:"scope"`
	Labels     map[string]string `json:"labels"`
	CreatedAt  string            `json:"created_at"`
	InUse      bool              `json:"in_use"`
}

type VolumeListResponse struct {
	Volumes []VolumeSummary `json:"volumes"`
}

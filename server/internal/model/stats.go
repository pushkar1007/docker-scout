package model

type DashboardData struct {
	Summary    SystemSummary   `json:"summary"`
	Containers []ContainerStats `json:"containers"`
}

type SystemSummary struct {
	ActiveContainers int    `json:"active_containers"`
	AvgCPU           string `json:"avg_cpu"`
	AvgMemory        string `json:"avg_memory"`
	AvgNetIO         string `json:"avg_net_io"`
	AvgDiskIO        string `json:"avg_disk_io"`
}

type StatsSnapshot struct {
	CPUPercent   float64
	MemoryBytes  uint64
	MemoryLimit  uint64
	NetRxBps     uint64
	NetTxBps     uint64
	DiskReadBps  uint64
	DiskWriteBps uint64
}

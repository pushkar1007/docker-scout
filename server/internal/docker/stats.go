package docker

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"docker-scout/internal/model"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
)

func ReadContainerStats(ctx context.Context, cli *client.Client, id string) (model.StatsSnapshot, error) {
	statsCtx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	stats, err := cli.ContainerStats(statsCtx, id, client.ContainerStatsOptions{Stream: true})
	if err != nil {
		return model.StatsSnapshot{}, err
	}
	defer stats.Body.Close()

	dec := json.NewDecoder(stats.Body)

	var prev container.StatsResponse
	if err := dec.Decode(&prev); err != nil {
		return model.StatsSnapshot{}, err
	}

	var cur container.StatsResponse
	if err := dec.Decode(&cur); err != nil {
		return model.StatsSnapshot{}, err
	}

	interval := cur.Read.Sub(prev.Read).Seconds()
	if interval <= 0 {
		interval = 1
	}

	return model.StatsSnapshot{
		CPUPercent:   computeCPUPercent(prev, cur),
		MemoryBytes:  computeMemoryBytes(cur),
		MemoryLimit:  cur.MemoryStats.Limit,
		NetRxBps:     computeNetRxBps(prev, cur, interval),
		NetTxBps:     computeNetTxBps(prev, cur, interval),
		DiskReadBps:  computeDiskReadBps(prev, cur, interval),
		DiskWriteBps: computeDiskWriteBps(prev, cur, interval),
	}, nil
}

func computeCPUPercent(prev, cur container.StatsResponse) float64 {
	cpuDelta := float64(cur.CPUStats.CPUUsage.TotalUsage - prev.CPUStats.CPUUsage.TotalUsage)
	sysDelta := float64(cur.CPUStats.SystemUsage - prev.CPUStats.SystemUsage)
	if cpuDelta <= 0 || sysDelta <= 0 {
		return 0
	}

	online := int(cur.CPUStats.OnlineCPUs)
	if online == 0 {
		online = len(cur.CPUStats.CPUUsage.PercpuUsage)
	}
	if online == 0 {
		online = 1
	}

	return (cpuDelta / sysDelta) * float64(online) * 100
}

func computeMemoryBytes(cur container.StatsResponse) uint64 {
	usage := cur.MemoryStats.Usage
	if cur.MemoryStats.Stats != nil {
		if cache, ok := cur.MemoryStats.Stats["cache"]; ok && usage > cache {
			usage -= cache
		}
	}
	return usage
}

func computeNetRxBps(prev, cur container.StatsResponse, interval float64) uint64 {
	prevRx, _ := sumNetworkBytes(prev.Networks)
	curRx, _ := sumNetworkBytes(cur.Networks)
	if interval <= 0 {
		return 0
	}
	return uint64(float64(curRx-prevRx) / interval)
}

func computeNetTxBps(prev, cur container.StatsResponse, interval float64) uint64 {
	_, prevTx := sumNetworkBytes(prev.Networks)
	_, curTx := sumNetworkBytes(cur.Networks)
	if interval <= 0 {
		return 0
	}
	return uint64(float64(curTx-prevTx) / interval)
}

func sumNetworkBytes(networks map[string]container.NetworkStats) (uint64, uint64) {
	var rx uint64
	var tx uint64
	for _, n := range networks {
		rx += n.RxBytes
		tx += n.TxBytes
	}
	return rx, tx
}

func computeDiskReadBps(prev, cur container.StatsResponse, interval float64) uint64 {
	prevRead, _ := sumDiskBytes(prev)
	curRead, _ := sumDiskBytes(cur)
	if interval <= 0 {
		return 0
	}
	return uint64(float64(curRead-prevRead) / interval)
}

func computeDiskWriteBps(prev, cur container.StatsResponse, interval float64) uint64 {
	_, prevWrite := sumDiskBytes(prev)
	_, curWrite := sumDiskBytes(cur)
	if interval <= 0 {
		return 0
	}
	return uint64(float64(curWrite-prevWrite) / interval)
}

func sumDiskBytes(stats container.StatsResponse) (uint64, uint64) {
	var read uint64
	var write uint64

	if stats.StorageStats.ReadSizeBytes > 0 || stats.StorageStats.WriteSizeBytes > 0 {
		return stats.StorageStats.ReadSizeBytes, stats.StorageStats.WriteSizeBytes
	}

	for _, entry := range stats.BlkioStats.IoServiceBytesRecursive {
		switch {
		case strings.EqualFold(entry.Op, "read"):
			read += entry.Value
		case strings.EqualFold(entry.Op, "write"):
			write += entry.Value
		}
	}

	return read, write
}

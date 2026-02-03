package state

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"docker-scout/internal/docker"
	"docker-scout/internal/model"
	"docker-scout/internal/util"

	"github.com/moby/moby/client"
)

func UpdateCache(cli *client.Client, cache *Cache) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	containers, err := docker.ListContainers(ctx, cli, true)
	if err != nil {
		return
	}

	var (
		tempStats = make([]model.ContainerStats, 0, len(containers.Items))
		totalCPU  float64
		totalMem  uint64
		totalNet  uint64
		totalDisk uint64
		count     int
	)

	for _, c := range containers.Items {
		if c.State != "running" && c.State != "paused" {
			continue
		}

		stats, statsErr := docker.ReadContainerStats(ctx, cli, c.ID)
		inspect, err := docker.InspectContainer(ctx, cli, c.ID)
		if err != nil {
			continue
		}

		name := ""
		if len(c.Names) > 0 {
			name = strings.TrimPrefix(c.Names[0], "/")
		}

		cs := model.ContainerStats{
			ID:       c.ID[:12],
			Name:     name,
			State:    string(c.State),
			LastUsed: util.ResolveLastUsed(inspect),
			Image:    c.Image,
			Labels:   docker.ResolveLabels(c.Labels, inspect),
			Ports:    util.FormatPorts(c.Ports),
			CPU:      util.FormatCPU(stats.CPUPercent, statsErr == nil),
			Memory:   util.FormatMemory(stats.MemoryBytes, stats.MemoryLimit, statsErr == nil),
			NetIO:    util.FormatNetIO(stats.NetRxBps, stats.NetTxBps, statsErr == nil),
			DiskIO:   util.FormatDiskIO(stats.DiskReadBps, stats.DiskWriteBps, statsErr == nil),
		}

		tempStats = append(tempStats, cs)
		if statsErr == nil {
			totalCPU += stats.CPUPercent
			totalMem += stats.MemoryBytes
			totalNet += (stats.NetRxBps + stats.NetTxBps)
			totalDisk += (stats.DiskReadBps + stats.DiskWriteBps)
			count++
		}
	}

	sort.Slice(tempStats, func(i, j int) bool {
		return tempStats[i].Name < tempStats[j].Name
	})

	summary := model.SystemSummary{
		ActiveContainers: count,
		AvgCPU:           "-",
		AvgMemory:        "-",
		AvgNetIO:         "-",
		AvgDiskIO:        "-",
	}
	if count > 0 {
		summary.AvgCPU = fmt.Sprintf("%.2f%%", totalCPU/float64(count))
		summary.AvgMemory = util.FormatBytes(uint64(float64(totalMem) / float64(count)))
		summary.AvgNetIO = util.FormatRate(uint64(float64(totalNet) / float64(count)))
		summary.AvgDiskIO = util.FormatRate(uint64(float64(totalDisk) / float64(count)))
	}

	cache.Set(model.DashboardData{Summary: summary, Containers: tempStats})
}

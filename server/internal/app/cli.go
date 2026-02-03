package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"docker-scout/internal/docker"
	"docker-scout/internal/util"
)

func RunCLI() {
	ctx := context.Background()

	cli, err := docker.NewClient()
	if err != nil {
		log.Fatalf("docker client init failed: %v", err)
	}

	containers, err := docker.ListContainers(ctx, cli, true)
	if err != nil {
		log.Fatalf("container list failed: %v", err)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "CONTAINER ID\tNAME\tSTATE\tLAST USED\tIMAGE\tLABELS\tPORTS\tCPU\tMEM\tNET I/O\tDISK R/W")

	for _, c := range containers.Items {
		if c.State != "running" && c.State != "paused" {
			continue
		}

		inspect, err := docker.InspectContainer(ctx, cli, c.ID)
		if err != nil {
			log.Printf("inspect failed for %s: %v", c.ID[:12], err)
			continue
		}

		lastUsed := util.ResolveLastUsed(inspect)
		stats, statsErr := docker.ReadContainerStats(ctx, cli, c.ID)
		if statsErr != nil {
			log.Printf("stats failed for %s: %v", c.ID[:12], statsErr)
		}

		name := ""
		if len(c.Names) > 0 {
			name = c.Names[0]
		}

		labels := util.FormatLabels(docker.ResolveLabels(c.Labels, inspect))
		ports := util.FormatPorts(c.Ports)
		cpu := util.FormatCPU(stats.CPUPercent, statsErr == nil)
		mem := util.FormatMemory(stats.MemoryBytes, stats.MemoryLimit, statsErr == nil)
		netIO := util.FormatNetIO(stats.NetRxBps, stats.NetTxBps, statsErr == nil)
		diskIO := util.FormatDiskIO(stats.DiskReadBps, stats.DiskWriteBps, statsErr == nil)

		fmt.Fprintf(
			w,
			"%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			c.ID[:12],
			name,
			c.State,
			lastUsed,
			c.Image,
			labels,
			ports,
			cpu,
			mem,
			netIO,
			diskIO,
		)
	}

	if err := w.Flush(); err != nil {
		log.Fatalf("output flush failed: %v", err)
	}
}

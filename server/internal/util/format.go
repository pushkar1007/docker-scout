package util

import (
	"fmt"
	"sort"
	"strings"

	"github.com/moby/moby/api/types/container"
)

func FormatLabels(labels map[string]string) string {
	if len(labels) == 0 {
		return "-"
	}

	keys := make([]string, 0, len(labels))
	for k := range labels {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%s", k, labels[k]))
	}

	return strings.Join(parts, ", ")
}

func FormatPorts(ports []container.PortSummary) string {
	if len(ports) == 0 {
		return "-"
	}

	parts := make([]string, 0, len(ports))
	for _, p := range ports {
		if p.PublicPort > 0 {
			parts = append(parts, fmt.Sprintf("%d", p.PublicPort))
			continue
		}

		parts = append(parts, fmt.Sprintf("%d", p.PrivatePort))
	}

	return strings.Join(parts, ", ")
}

func FormatCPU(cpu float64, ok bool) string {
	if !ok {
		return "-"
	}
	return fmt.Sprintf("%.2f%%", cpu)
}

func FormatMemory(usage, limit uint64, ok bool) string {
	if !ok {
		return "-"
	}
	if limit == 0 {
		return FormatBytes(usage)
	}
	percent := (float64(usage) / float64(limit)) * 100
	return fmt.Sprintf("%s (%.0f%%)", FormatBytes(usage), percent)
}

func FormatNetIO(rxBps, txBps uint64, ok bool) string {
	if !ok {
		return "-"
	}
	return fmt.Sprintf("%s/%s", FormatRate(rxBps), FormatRate(txBps))
}

func FormatDiskIO(readBps, writeBps uint64, ok bool) string {
	if !ok {
		return "-"
	}
	return fmt.Sprintf("%s/%s", FormatRate(readBps), FormatRate(writeBps))
}

func FormatRate(bps uint64) string {
	return fmt.Sprintf("%s/s", FormatBytes(bps))
}

func FormatBytes(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%dB", b)
	}
	div := float64(unit)
	suffix := "KiB"
	val := float64(b) / div
	if val >= unit {
		div *= unit
		suffix = "MiB"
		val = float64(b) / div
	}
	if val >= unit {
		div *= unit
		suffix = "GiB"
		val = float64(b) / div
	}
	return fmt.Sprintf("%.2f%s", val, suffix)
}

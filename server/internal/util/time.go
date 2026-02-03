package util

import (
	"fmt"
	"time"

	"github.com/moby/moby/client"
)

// ResolveLastUsed maps Docker's limited timestamps to a "last used" label.
func ResolveLastUsed(inspect client.ContainerInspectResult) string {
	if inspect.Container.State == nil {
		return "unknown"
	}

	if inspect.Container.State.StartedAt == "" {
		return "unknown"
	}

	t, err := time.Parse(time.RFC3339Nano, inspect.Container.State.StartedAt)
	if err != nil {
		return "unknown"
	}

	age := time.Since(t)
	if age < 0 {
		age = 0
	}

	if age < 24*time.Hour {
		hours := int(age.Hours())
		if hours < 1 {
			hours = 1
		}
		return fmt.Sprintf("%dh", hours)
	}

	days := int(age.Hours() / 24)
	if days < 1 {
		days = 1
	}
	if days == 1 {
		return "1 day"
	}
	return fmt.Sprintf("%d days", days)
}

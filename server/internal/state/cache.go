package state

import (
	"sync"
	"time"

	"docker-scout/internal/model"
)

// Cache holds the latest dashboard data and stats snapshot.
// It is updated by a background updater and read by /stats API handler.
// RWMutex allows concurrent reads; version field enables change detection.
type Cache struct {
	mu         sync.RWMutex
	data       model.DashboardData
	lastUpdate time.Time
	version    uint64
	cond       *sync.Cond
}

func NewCache() *Cache {
	c := &Cache{
		data: model.DashboardData{
			Summary: model.SystemSummary{
				AvgCPU:    "-",
				AvgMemory: "-",
				AvgNetIO:  "-",
				AvgDiskIO: "-",
			},
			Containers: []model.ContainerStats{},
		},
	}
	c.cond = sync.NewCond(&c.mu)
	return c
}

// Get returns a snapshot of the current dashboard data.
// Reads hold only RLock, allowing concurrent GET requests to proceed.
func (c *Cache) Get() model.DashboardData {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.data
}

// Set updates the cache atomically and signals all waiters.
// The version counter increments on every update; future extensions can use this for change detection.
func (c *Cache) Set(data model.DashboardData) {
	c.mu.Lock()
	c.data = data
	c.lastUpdate = time.Now()
	c.version++
	c.cond.Broadcast()
	c.mu.Unlock()
}

package state

import (
	"sync"
	"time"

	"docker-scout/internal/model"
)

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

func (c *Cache) Get() model.DashboardData {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.data
}

func (c *Cache) Set(data model.DashboardData) {
	c.mu.Lock()
	c.data = data
	c.lastUpdate = time.Now()
	c.version++
	c.cond.Broadcast()
	c.mu.Unlock()
}

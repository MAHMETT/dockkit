package perf

import (
	"runtime"
	"sync"
	"time"
)

type Stats struct {
	Uptime         time.Duration
	ServicesTotal  int
	ServicesActive int
	DockerAPICalls int64
	HubAPICalls    int64
	CacheHits      int64
	CacheMisses    int64
	Goroutines     int
	MemoryAlloc    uint64
	LastPollTime   time.Duration
	mu             sync.RWMutex
}

var startTime = time.Now()

func GetStats() *Stats {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return &Stats{
		Uptime:      time.Since(startTime),
		Goroutines:  runtime.NumGoroutine(),
		MemoryAlloc: m.Alloc,
	}
}

func (s *Stats) RecordDockerCall() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.DockerAPICalls++
}

func (s *Stats) RecordHubCall() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.HubAPICalls++
}

func (s *Stats) RecordCacheHit() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.CacheHits++
}

func (s *Stats) RecordCacheMiss() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.CacheMisses++
}

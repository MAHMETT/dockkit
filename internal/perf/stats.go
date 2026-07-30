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

var processStart = time.Now()

// GetStats returns current runtime stats.
func GetStats() *Stats {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return &Stats{
		Uptime:      time.Since(processStart),
		Goroutines:  runtime.NumGoroutine(),
		MemoryAlloc: m.Alloc,
	}
}

func (s *Stats) IncDockerCalls() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.DockerAPICalls++
}

func (s *Stats) IncHubCalls() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.HubAPICalls++
}

func (s *Stats) IncCacheHits() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.CacheHits++
}

func (s *Stats) IncCacheMisses() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.CacheMisses++
}

// Deprecated: Use IncDockerCalls instead.
func (s *Stats) RecordDockerCall() { s.IncDockerCalls() }

// Deprecated: Use IncHubCalls instead.
func (s *Stats) RecordHubCall() { s.IncHubCalls() }

// Deprecated: Use IncCacheHits instead.
func (s *Stats) RecordCacheHit() { s.IncCacheHits() }

// Deprecated: Use IncCacheMisses instead.
func (s *Stats) RecordCacheMiss() { s.IncCacheMisses() }

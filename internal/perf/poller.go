package perf

import (
	"time"
)

type PollState int

const (
	PollActive PollState = iota
	PollIdle
	PollPaused
)

type Poller struct {
	state         PollState
	activeInterval time.Duration
	idleInterval   time.Duration
	pausedInterval  time.Duration
	changeCount    int
	lastChange     time.Time
	onPoll         func()
}

func NewPoller(onPoll func()) *Poller {
	return &Poller{
		state:          PollActive,
		activeInterval: 5 * time.Second,
		idleInterval:   30 * time.Second,
		pausedInterval: 60 * time.Second,
		lastChange:     time.Now(),
		onPoll:         onPoll,
	}
}

func (p *Poller) Interval() time.Duration {
	switch p.state {
	case PollActive:
		return p.activeInterval
	case PollIdle:
		return p.idleInterval
	default:
		return p.pausedInterval
	}
}

func (p *Poller) RecordChange() {
	p.changeCount = 0
	p.state = PollActive
	p.lastChange = time.Now()
}

func (p *Poller) RecordNoChange() {
	p.changeCount++
	if p.changeCount >= 3 {
		p.state = PollIdle
	}
}

func (p *Poller) Pause() {
	p.state = PollPaused
}

func (p *Poller) Resume() {
	p.state = PollActive
}

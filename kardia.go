package kardia

import (
	"fmt"
	"time"
)

type Heart interface {
	Beat()
}

type heart struct {
	done      <-chan any
	Heartbeat chan any
	pulse     <-chan time.Time
	Results   chan time.Time
	workGen   <-chan time.Time
}

func New(done <-chan any, pulseInterval time.Duration) *heart {
	pulse := time.Tick(pulseInterval)
	workGen := time.Tick(2 * pulseInterval)
	heartbeat := make(chan any)
	results := make(chan time.Time)
	return &heart{
		done:      done,
		Heartbeat: heartbeat,
		pulse:     pulse,
		Results:   results,
		workGen:   workGen,
	}
}

func (h *heart) Beat() {
	defer close(h.Heartbeat)
	for {
		select {
		case <-h.done:
			fmt.Println("DONE!")
			return
		case <-h.pulse:
			h.sendPulse()
		}
	}
}

func (h *heart) sendPulse() {
	select {
	case h.Heartbeat <- struct{}{}:
	default:
	}
}

func (h *heart) sendResult(r time.Time) {
	for {
		select {
		case <-h.done:
			return
		case <-h.pulse:
			h.sendPulse()
		case h.Results <- r:
			return
		}
	}
}

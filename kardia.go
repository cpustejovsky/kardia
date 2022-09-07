package kardia

import (
	"time"
)

// Result contains results that may be useful to have from a heartbeat pattern
// TODO: Determine if Result struct is needed; if so, what else should be added other than current time?
type Result struct {
	Time time.Time
}

type resultsChan chan Result

type heart struct {
	done      <-chan any
	Heartbeat chan any
	pulse     <-chan time.Time
	Results   resultsChan
	workGen   <-chan time.Time
}

// New is a constructor function for heart; takes a done channel and pulse interval
func New(done <-chan any, pulseInterval time.Duration) *heart {
	pulse := time.Tick(pulseInterval)
	workGen := time.Tick(2 * pulseInterval)
	heartbeat := make(chan any)
	results := make(resultsChan)
	return &heart{
		done:      done,
		Heartbeat: heartbeat,
		pulse:     pulse,
		Results:   results,
		workGen:   workGen,
	}
}

// Beat is a method that sends Results and pulses until it reads from the done channel
func (h *heart) Beat() {
	defer close(h.Heartbeat)
	defer close(h.Results)
	for {
		select {
		case <-h.done:
			return
		case r := <-h.workGen:
			h.sendResult(r)
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
		case h.Results <- Result{Time: r}:
			return
		}
	}
}

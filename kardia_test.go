package kardia_test

import (
	"github.com/cpustejovsky/kardia"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestHeart_Beat(t *testing.T) {
	done := make(chan any)
	timeout := 20 * time.Millisecond
	time.AfterFunc(10*time.Millisecond, func() { close(done) })
	const interval = 1 * time.Millisecond
	heart := kardia.New(done, interval)
	go heart.Beat()
	counter := 0
	for {
		select {
		case _, ok := <-heart.Heartbeat:
			if ok == false {
				return
			}
			counter++
		case <-done:
		case <-time.After(timeout):
		}
	}
	assert.Equal(t, counter, 10)
}

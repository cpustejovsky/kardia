package kardia_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"kardia"
	"testing"
	"time"
)

func TestHeart_Beat(t *testing.T) {
	done := make(chan any)
	timeout := 11 * time.Second
	time.AfterFunc(10*time.Second, func() { close(done) })
	const interval = 1 * time.Second
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
			fmt.Println("pulse")
		case <-time.After(timeout):
		}
	}
	assert.Equal(t, counter, 10)
}

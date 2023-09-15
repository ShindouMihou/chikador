package chikador

import (
	"math"
	"sync"
	"time"
)

func withoutDedupe(chismis *Chismis) {
	go func() {
		for {
			select {
			case event, ok := <-chismis.watcher.Events:
				if !ok {
					return
				}
				chismis.queue.append(&Message{
					Event:    event.Op,
					Filename: event.Name,
					Error:    nil,
				})
			case err, ok := <-chismis.watcher.Errors:
				if !ok {
					return
				}
				chismis.queue.append(&Message{
					Error: err,
				})
			}
		}
	}()
}

func withDedupe(chismis *Chismis) {
	var (
		waitFor = 100 * time.Millisecond
		mu      sync.Mutex
		timers  = make(map[string]*time.Timer)
	)
	go func() {
		for {
			select {
			case event, ok := <-chismis.watcher.Events:
				if !ok {
					return
				}
				mu.Lock()
				timer, ok := timers[event.Name]
				mu.Unlock()
				if !ok {
					timer = time.AfterFunc(math.MaxInt64, func() {
						chismis.queue.append(&Message{
							Event:    event.Op,
							Filename: event.Name,
							Error:    nil,
						})

						mu.Lock()
						delete(timers, event.Name)
						mu.Unlock()
					})
					timer.Stop()

					mu.Lock()
					timers[event.Name] = timer
					mu.Unlock()
				}

				timer.Reset(waitFor)
			case err, ok := <-chismis.watcher.Errors:
				if !ok {
					return
				}
				chismis.queue.append(&Message{
					Error: err,
				})
			}
		}
	}()
}

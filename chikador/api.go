package chikador

import (
	"github.com/fsnotify/fsnotify"
)

type Chismis struct {
	queue   eventQueue
	watcher *fsnotify.Watcher
	closed  bool
}

type Message struct {
	Event    fsnotify.Op
	Filename string
	Error    error
}

// Watch creates a new Chismis instance that is watching over the given path, and its subdirectories if Recursive is
// enabled. You can add the following options to give more power to the file watching capabilities:
//
// * Recursive: adds all the subdirectories under the directory and the subdirectories' subdirectories' subdirectories' sub.... you get it.
//
// * WithDedupe: (recommended) dedupes events which operating systems tend to have, this uses the implementation from fsnotify's own repository.
func Watch(path string, opts ...Option) (*Chismis, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	if err = watcher.Add(path); err != nil {
		_ = watcher.Close()
		return nil, err
	}
	options2 := &options{
		recursive: false,
		kind:      withoutDedupe,
	}
	for _, opt := range opts {
		opt(options2)
	}
	if options2.recursive {
		if err := add(path, watcher); err != nil {
			_ = watcher.Close()
			return nil, err
		}
	}
	chismis := &Chismis{
		queue:   eventQueue{},
		watcher: watcher,
		closed:  false,
	}
	options2.kind(chismis)
	return chismis, nil
}

// Poll tries to poll data from the event queue, if there is none then it will return a nil pointer.
// Recommended to use in a forever for-loop to continuously poll data. Alternatively, you can use Listen which is the
// same as a for-loop, but asynchronous and shorter.
//
// Do note that you shouldn't have this twice, nor should you also use Listen if you are already polling data with this
// method as the two methods will battle it out over who has the resource.
func (chismis *Chismis) Poll() *Message {
	return chismis.queue.poll()
}

// IsClosed checks whether this Chismis channel is already closed, important if you want to stop program immediately
// when closed.
func (chismis *Chismis) IsClosed() bool {
	return chismis.closed
}

// Listen listens to events from the event queue, this is running in another goroutine and uses the same polling method
// internally.
//
// Do note that you shouldn't listen to a Chismis instance twice, only one listener ever, not even Poll, as this will
// result in the resource being contended.
func (chismis *Chismis) Listen(fn func(msg *Message)) {
	go func() {
		for {
			if chismis.closed {
				break
			}
			msg := chismis.Poll()
			if msg == nil {
				continue
			}
			fn(msg)
		}
	}()
}

// Close removes all watches and closes all event channels.
func (chismis *Chismis) Close() error {
	chismis.closed = true
	return chismis.watcher.Close()
}

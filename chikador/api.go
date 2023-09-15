package chikador

import (
	"github.com/fsnotify/fsnotify"
)

type Chismis struct {
	queue   eventQueue
	watcher *fsnotify.Watcher
}

type Message struct {
	Event    fsnotify.Op
	Filename string
	Error    error
}

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
	}
	options2.kind(chismis)
	return chismis, nil
}

func (chismis *Chismis) Poll() *Message {
	return chismis.queue.poll()
}

func (chismis *Chismis) Listen(fn func(msg *Message)) {
	go func() {
		for {
			msg := chismis.Poll()
			if msg == nil {
				continue
			}
			fn(msg)
		}
	}()
}

func (chismis *Chismis) Close() error {
	return chismis.watcher.Close()
}

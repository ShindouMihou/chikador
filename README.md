# chikador

*because of chismis, i became a chikador.* Chikador is a simple overlay, or abstraction, over [`fsnotify`](https://github.com/fsnotify/fsnotify) 
that brings a simpler way of using [`fsnotify`](https://github.com/fsnotify/fsnotify), which by itself is already a cool tool that does many 
things very easily.

## Demo
```go
package main

import (
	"fmt"
	"github.com/ShindouMihou/chikador/chikador"
	"log"
)

func main() {
	chismis, err := chikador.Watch(".tests/", chikador.WithDedupe)
	if err != nil {
		log.Fatalln(err)
	}
	defer chismis.Close()
	chismis.Listen(func(msg *chikador.Message) {
		fmt.Println("received ", msg.Event, " from ", msg.Filename)
	})
	<-make(chan struct{})
}
```

## installation
```go
go get github.com/ShindouMihou/chikador
```

## file watcher

for an examplified-view of the different methods, please view our examples:
- [`async`](examples/async/main.go): asynchronously send events to a callback
- [`polling`](examples/polling/main.go): synchronously poll events from the queue

as explained from the above, we have two methods of polling events (`async` and `polling`). internally, `async` is simply 
`polling` but ran in another goroutine, but to demonstrate how we can create the two different methods:

### `async`

asynchronous is the easiest way, you can simply use the `chismis.Listen` method to create an async watcher.

```go
func main() {
    chismis, err := chikador.Watch(".tests/", chikador.WithDedupe)
    if err != nil {
        log.Fatalln(err)
    }
    defer chismis.Close()
	chismis.Listen(func(msg *chikador.Message) {
        fmt.Println("received ", msg.Event, " from ", msg.Filename)
    })
    <-make(chan struct{})
}
```

### `polling`

when you want to handle your way of asynchronous yourself, or simply want to use the for-loop interface yourself, then 
you can use the `polling` method by creating a `Chismis` instance using `chikador.Watch` method. unlike the `async` 
method, we have the `Poll` method which returns a `*Message` that we can use to poll:
```go
func main() {
    chismis, err := chikador.Watch(".tests/", chikador.WithDedupe)
    if err != nil {
        log.Fatalln(err)
    }
    defer chismis.Close()
    for {
        message := chismis.Poll()
        if message != nil {
            continue
        }
        fmt.Println("received ", message.Event, " from ", message.Filename)
    }
}
```

## `dedupe`?

As stated in [`fsnotify`](https://github.com/fsnotify/fsnotify) , the operating system may duplicate events, and in one of their examples, they've included a 
deduplicator which we've brought down to chikador that can be enabled by adding the `chikador.WithDedupe` option, as seen, 
in the examples, such as in:
- [`async`](examples/async/main.go)

## `recursive`?

[`fsnotify`](https://github.com/fsnotify/fsnotify) will listen to the directory's files, but it cannot listen onto the subdirectories, which is why we include a little 
utility to have `chikador` scan through all subdirectories and their corresponding subdirectories to listen to changes on them as well. You can enable this behavior by adding 
the `chikador.Recursive` option, see example in:
- [`recursive`](examples/recursive/main.go)
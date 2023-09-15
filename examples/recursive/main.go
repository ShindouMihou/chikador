package main

import (
	"fmt"
	"github.com/ShindouMihou/chikador/chikador"
	"log"
)

func main() {
	chismis, err := chikador.Watch(".tests/", chikador.Recursive, chikador.WithDedupe)
	if err != nil {
		log.Fatalln(err)
	}
	defer chismis.Close()
	chismis.Listen(func(msg *chikador.Message) {
		fmt.Println("received ", msg.Event, " from ", msg.Filename)
	})
	<-make(chan struct{})
}

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
	for {
		message := chismis.Poll()
		if message != nil {
			continue
		}
		fmt.Println("received ", message.Event, " from ", message.Filename)
	}
}

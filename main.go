package main

import (
	"log"
	"sync"
)

const bufferSize = 10

func main() {
	cfg, err := parseConfig("sunny_5_skiers/config.json")
	if err != nil {
		log.Panic("Cannot load config")
	}

	_ = cfg

	ch := make(chan Event, bufferSize)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		parseEvents("../sunny_5_skiers/events", ch)
	}()

	wg.Wait()
}

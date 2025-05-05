package main

import (
	"io"
	"os"
	"sync"
)

const bufferSize = 10

func getSortedCompetitors(cfg Config, from io.Reader) []Competitor {
	events := make(chan Event, bufferSize)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		parseEvents(from, events)
	}()

	competitors := make(map[int]*Competitor)
	wg.Add(1)
	go func() {
		defer wg.Done()
		processEvents(events, competitors, cfg)
	}()

	wg.Wait()

	calculateCompetitors(competitors)

	return getSortedByTime(competitors)
}

func main() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		panic("CONFIG_PATH not set")
	}

	cfg, err := parseConfig(configPath)
	if err != nil {
		panic(err)
	}

	eventsInput := io.Reader(os.Stdin)
	competitors := getSortedCompetitors(cfg, eventsInput)

	printFinalReport(competitors, cfg)
}

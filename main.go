package main

import (
	"fmt"
	"sync"
)

const bufferSize = 10

func getSortedCompetitors(cfg Config, eventsPath string) []Competitor {
	events := make(chan Event, bufferSize)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		parseEvents(eventsPath, events)
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
	configPath := "sunny_5_skiers/config.json"
	eventsPath := "sunny_5_skiers/events"

	cfg, err := parseConfig(configPath)
	if err != nil {
		panic(err)
	}

	competitors := getSortedCompetitors(cfg, eventsPath)

	fmt.Println("\nFinal Report:")
	printFinalReport(competitors, cfg)
}

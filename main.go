package main

import (
	"io"
	"os"
)

const bufferSize = 10

func getSortedCompetitors(cfg Config, from io.Reader) []Competitor {
	eventsCh := make(chan Event, bufferSize)

	go parseEvents(from, eventsCh)

	competitors := make(map[int]*Competitor)
	processEvents(eventsCh, competitors, cfg)

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

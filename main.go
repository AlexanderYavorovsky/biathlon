package main

import (
	"io"
	"log"
	"os"
)

const (
	bufferSize = 10
	configEnv  = "CONFIG_PATH"
)

func getSortedCompetitors(cfg Config, eventStream io.Reader) []Competitor {
	eventCh := parseEvents(eventStream)

	competitors := make(map[int]*Competitor)
	processEvents(eventCh, competitors, cfg)

	return getSortedByTime(getList(competitors))
}

func main() {
	configPath := os.Getenv(configEnv)
	if configPath == "" {
		log.Fatalf("%s is not set", configEnv)
	}

	cfg, err := parseConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}

	eventsInput := os.Stdin
	competitors := getSortedCompetitors(cfg, eventsInput)

	printFinalReport(competitors, cfg)
}

package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	Laps        int    `json:"laps"`
	LapLen      int    `json:"lapLen"`
	PenaltyLen  int    `json:"penaltyLen"`
	FiringLines int    `json:"firingLines"`
	Start       string `json:"start"`
	StartDelta  string `json:"startDelta"`
}

func parseConfig(path string) Config {
	f, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	cfg := Config{}
	err = json.Unmarshal(f, &cfg)
	if err != nil {
		panic(err)
	}

	return cfg
}

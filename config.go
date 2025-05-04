package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	timeLayout     = "15:04:05.000"
	durationLayout = "15:04:05"
)

type Config struct {
	Laps        int      `json:"laps"`
	LapLen      int      `json:"lapLen"`
	PenaltyLen  int      `json:"penaltyLen"`
	FiringLines int      `json:"firingLines"`
	Start       Time     `json:"start"`
	StartDelta  Duration `json:"startDelta"`
}

type Time struct {
	time.Time
}

type Duration struct {
	time.Duration
}

func (t *Time) UnmarshalJSON(data []byte) error {
	timeStr := strings.Trim(string(data), `"`)
	if timeStr == "" {
		return nil
	}

	parsed, err := time.Parse(timeLayout, timeStr)
	if err != nil {
		return fmt.Errorf("invalid time format: %w", err)
	}

	t.Time = parsed

	return nil
}

func (d *Duration) UnmarshalJSON(data []byte) error {
	timeStr := strings.Trim(string(data), `"`)
	if timeStr == "" {
		return nil
	}

	parsed, err := time.Parse(durationLayout, timeStr)
	if err != nil {
		return fmt.Errorf("invalid time format: %w", err)
	}

	d.Duration = time.Duration(parsed.Second())*time.Second +
		time.Duration(parsed.Minute())*time.Minute +
		time.Duration(parsed.Hour())*time.Hour

	return nil
}

func parseConfig(path string) (Config, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	cfg := Config{}
	err = json.Unmarshal(f, &cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}

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
	durationLayout = time.TimeOnly

	errInvalidTime           = "invalid time format for Time: %w"
	errInvalidDuration       = "invalid time format for Duration: %w"
	errCannotReadFile        = "cannot read file: %w"
	errCannotUnmarshalConfig = "cannot unmarshal config: %w"
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

	layouts := []string{timeLayout, durationLayout}
	var err error
	for _, l := range layouts {
		parsed, err := time.Parse(l, timeStr)
		if err == nil {
			t.Time = parsed
			return nil
		}
	}

	return fmt.Errorf(errInvalidTime, err)
}

func (d *Duration) UnmarshalJSON(data []byte) error {
	timeStr := strings.Trim(string(data), `"`)
	if timeStr == "" {
		return nil
	}

	parsed, err := time.Parse(durationLayout, timeStr)
	if err != nil {
		return fmt.Errorf(errInvalidDuration, err)
	}

	d.Duration = time.Duration(parsed.Second())*time.Second +
		time.Duration(parsed.Minute())*time.Minute +
		time.Duration(parsed.Hour())*time.Hour

	return nil
}

func parseConfig(path string) (Config, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf(errCannotReadFile, err)
	}

	cfg := Config{}
	err = json.Unmarshal(f, &cfg)
	if err != nil {
		return Config{}, fmt.Errorf(errCannotUnmarshalConfig, err)
	}

	return cfg, nil
}

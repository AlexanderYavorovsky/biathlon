package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestResultString(t *testing.T) {
	cfgTwoFiringLines := Config{
		Laps:        2,
		LapLen:      3500,
		PenaltyLen:  150,
		FiringLines: 2,
		Start:       Time{time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC)},
		StartDelta:  Duration{time.Minute + 30*time.Second},
	}

	comp := Competitor{
		ID:                  1,
		Status:              StatusFinished,
		PlannedStart:        time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC),
		TotalTime:           30 * time.Minute,
		CurrentLapStartTime: time.Date(0, 1, 1, 10, 30, 0, 0, time.UTC),
		Laps: []Lap{
			{
				Duration: 10 * time.Minute,
				Speed:    float64(3500) / float64(10*60),
			},
			{
				Duration: 20 * time.Minute,
				Speed:    float64(3500) / float64(20*60),
			},
		},
		Penalty: time.Minute,
		Hits:    8,
	}

	tests := []struct {
		name string
		comp Competitor
		cfg  Config
		want string
	}{
		{
			name: "finished",
			comp: comp,
			cfg:  cfgTwoFiringLines,
			want: "[00:30:00.000] 1 [{00:10:00.000, 5.833}, {00:20:00.000, 2.917}] {00:01:00.000, 5.000} 8/10",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			got := getResult(tt.comp, tt.cfg)
			assert.Equal(tt.want, got.String())
		})
	}
}

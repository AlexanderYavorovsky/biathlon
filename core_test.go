package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestProcessEvent(t *testing.T) {
	cfgOneLap := Config{
		Laps:        1,
		LapLen:      3500,
		PenaltyLen:  150,
		FiringLines: 1,
		Start:       Time{time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC)},
		StartDelta:  Duration{time.Minute + 30*time.Second},
	}
	tests := []struct {
		name        string
		event       Event
		competitors map[int]*Competitor
		cfg         Config
		want        map[int]*Competitor
	}{
		{
			name: "event start set",
			event: Event{
				ID:           EventStartSet,
				CompetitorID: 1,
				Time:         time.Date(0, 1, 1, 9, 30, 0, 0, time.UTC),
				ExtraParams:  []string{"10:00:00.000"},
			},
			competitors: map[int]*Competitor{},
			cfg:         cfgOneLap,
			want: map[int]*Competitor{
				1: {
					ID:           1,
					PlannedStart: time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name: "event started after deadline",
			event: Event{
				ID:           EventStarted,
				CompetitorID: 1,
				Time:         time.Date(0, 1, 1, 10, 2, 0, 0, time.UTC),
			},
			competitors: map[int]*Competitor{
				1: {
					ID:           1,
					PlannedStart: time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC),
				},
			},
			cfg: cfgOneLap,
			want: map[int]*Competitor{
				1: {
					ID:                  1,
					Status:              StatusNotStarted,
					PlannedStart:        time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC),
					CurrentLapStartTime: time.Date(0, 1, 1, 10, 2, 0, 0, time.UTC),
				},
			},
		},
		{
			name: "event ended main lap and finished",
			event: Event{
				ID:           EventEndedMainLap,
				CompetitorID: 1,
				Time:         time.Date(0, 1, 1, 10, 30, 0, 0, time.UTC),
			},
			competitors: map[int]*Competitor{
				1: {
					ID:                  1,
					PlannedStart:        time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC),
					CurrentLapStartTime: time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC),
				},
			},
			cfg: cfgOneLap,
			want: map[int]*Competitor{
				1: {
					ID:                  1,
					Status:              StatusFinished,
					PlannedStart:        time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC),
					TotalTime:           30 * time.Minute,
					CurrentLapStartTime: time.Date(0, 1, 1, 10, 30, 0, 0, time.UTC),
					Laps:                []Lap{{Duration: 30 * time.Minute, Speed: float64(3500) / float64(30*60)}},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			processEvent(tt.event, tt.competitors, cfgOneLap)
			assert.Equal(tt.want, tt.competitors)
		})
	}
}

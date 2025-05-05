package main

import (
	"fmt"
	"time"
)

const (
	StatusNotStarted     = "NotStarted"
	StatusNotFinished    = "NotFinished"
	StatusFinished       = "Finished"
	TargetsPerFiringLine = 5
)

type Lap struct {
	Duration time.Duration
	Speed    float64
}

type Competitor struct {
	ID                   int
	Status               string
	PlannedStart         time.Time
	StartTime            time.Time
	TotalTime            time.Duration
	LastLapStartTime     time.Time
	LastPenaltyStartTime time.Time
	Laps                 []Lap
	Penalty              time.Duration
	Hits                 int
}

func processEvents(events <-chan Event, competitors map[int]*Competitor, cfg Config) {
	for e := range events {
		fmt.Println(e)
		processEvent(e, competitors, cfg)
	}
}

func processEvent(e Event, competitors map[int]*Competitor, cfg Config) {
	competitor, ok := competitors[e.CompetitorID]
	if !ok {
		competitor = &Competitor{ID: e.CompetitorID}
		competitors[e.CompetitorID] = competitor
	}

	switch e.ID {
	case EventStartSet:
		processStartSet(e, competitor)
	case EventStarted:
		processStarted(cfg, e, competitor)
	case EventHit:
		competitor.Hits++
	case EventEnteredPenalty:
		competitor.LastPenaltyStartTime = e.Time
	case EventLeftPenalty:
		duration := e.Time.Sub(competitor.LastPenaltyStartTime)
		competitor.Penalty += duration
	case EventEndedMainLap:
		processEndedMainLap(cfg, e, competitor)
	case EventCantContinue:
		competitor.Status = StatusNotFinished
	}
}

func processStartSet(e Event, comp *Competitor) {
	t, err := time.Parse(timeLayout, e.ExtraParams[0])
	if err != nil {
		return
	}
	comp.PlannedStart = t
}

func processStarted(cfg Config, e Event, comp *Competitor) {
	comp.StartTime = e.Time
	comp.LastLapStartTime = e.Time
	if comp.StartTime.After(comp.PlannedStart.Add(cfg.StartDelta.Duration)) {
		comp.Status = StatusNotStarted
		event := Event{
			ID:           EventDisqualified,
			CompetitorID: e.CompetitorID,
			Time:         e.Time,
		}
		fmt.Println(event)
	}
}

func processEndedMainLap(cfg Config, e Event, comp *Competitor) {
	duration := e.Time.Sub(comp.LastLapStartTime)
	if len(comp.Laps) == 0 { // add start difference to first lap time
		duration += comp.StartTime.Sub(comp.PlannedStart)
	}

	speed := float64(cfg.LapLen) / duration.Seconds()
	lap := Lap{Duration: duration, Speed: speed}
	comp.Laps = append(comp.Laps, lap)
	comp.LastLapStartTime = e.Time

	if len(comp.Laps) == cfg.Laps { // final lap finished
		comp.Status = StatusFinished
		event := Event{
			ID:           EventFinished,
			CompetitorID: e.CompetitorID,
			Time:         e.Time,
		}
		fmt.Println(event)
	}
}

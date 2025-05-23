package main

import (
	"fmt"
	"log"
	"time"
)

const (
	StatusNotStarted  = "NotStarted"
	StatusNotFinished = "NotFinished"
	StatusFinished    = "Finished"

	TargetsPerFiringLine = 5
)

type Lap struct {
	Duration time.Duration
	Speed    float64
}

type Competitor struct {
	ID                      int
	Status                  string
	PlannedStart            time.Time
	TotalTime               time.Duration
	CurrentLapStartTime     time.Time
	CurrentPenaltyStartTime time.Time
	Laps                    []Lap
	Penalty                 time.Duration
	Hits                    int
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

	if competitor.Status == StatusNotFinished || competitor.Status == StatusNotStarted {
		return
	}

	switch e.ID {
	case EventStartSet:
		processStartSet(e, competitor)
	case EventStarted:
		processStarted(e, competitor, cfg)
	case EventHit:
		competitor.Hits++
	case EventEnteredPenalty:
		competitor.CurrentPenaltyStartTime = e.Time
	case EventLeftPenalty:
		duration := e.Time.Sub(competitor.CurrentPenaltyStartTime)
		competitor.Penalty += duration
	case EventEndedMainLap:
		processEndedMainLap(e, competitor, cfg)
	case EventCantContinue:
		competitor.Status = StatusNotFinished
	}
}

func processStartSet(e Event, comp *Competitor) {
	if len(e.ExtraParams) == 0 {
		log.Printf("Missing ExtraParams for event %d", e.ID)
		return
	}

	t, err := time.Parse(timeLayout, e.ExtraParams[0])
	if err != nil {
		log.Printf("Cannot parse start time for competitor %d: %v", comp.ID, err)
		return
	}
	comp.PlannedStart = t
}

func processStarted(e Event, comp *Competitor, cfg Config) {
	comp.CurrentLapStartTime = e.Time
	deadline := comp.PlannedStart.Add(cfg.StartDelta.Duration)
	if e.Time.Before(comp.PlannedStart) || e.Time.After(deadline) {
		comp.Status = StatusNotStarted
		event := createEvent(EventDisqualified, e.CompetitorID, e.Time)
		fmt.Println(event)
	}
}

func processEndedMainLap(e Event, comp *Competitor, cfg Config) {
	duration := e.Time.Sub(comp.CurrentLapStartTime)
	if len(comp.Laps) == 0 { // add start difference to first lap time
		duration += comp.CurrentLapStartTime.Sub(comp.PlannedStart)
	}

	speed := float64(cfg.LapLen) / duration.Seconds()
	comp.Laps = append(comp.Laps, Lap{
		Duration: duration,
		Speed:    speed,
	})
	comp.CurrentLapStartTime = e.Time
	comp.TotalTime += duration

	if len(comp.Laps) == cfg.Laps { // final lap finished
		comp.Status = StatusFinished
		event := createEvent(EventFinished, e.CompetitorID, e.Time)
		fmt.Println(event)
	}
}

func createEvent(eventID int, competitorID int, time time.Time) Event {
	return Event{
		ID:           eventID,
		CompetitorID: competitorID,
		Time:         time,
	}
}

package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	EventRegistered = 1 + iota
	EventStartSet
	EventOnStart
	EventStarted
	EventOnRange
	EventHit
	EventLeftRange
	EventEnteredPenalty
	EventLeftPenalty
	EventEndedMainLap
	EventCantContinue
)

const (
	EventDisqualified = 32 + iota
	EventFinished
)

const (
	errInvalidLineFormat = "invalid line format: %s"
	errCannotParseTime   = "cannot parse time: %w"
	errCannotConvert     = "cannot convert %s: %w"
)

type Event struct {
	ID           int
	CompetitorID int
	Time         time.Time
	ExtraParams  []string
}

func (e Event) String() string {
	eventTime := e.Time.Format(timeLayout)
	switch e.ID {
	case EventRegistered:
		return fmt.Sprintf("[%s] The competitor(%d) registered", eventTime, e.CompetitorID)
	case EventStartSet:
		return fmt.Sprintf(
			"[%s] The start time for the competitor(%d) was set by a draw to %s",
			eventTime,
			e.CompetitorID,
			e.ExtraParams[0],
		)
	case EventOnStart:
		return fmt.Sprintf(
			"[%s] The competitor(%d) is on the start line",
			eventTime,
			e.CompetitorID,
		)
	case EventStarted:
		return fmt.Sprintf("[%s] The competitor(%d) has started", eventTime, e.CompetitorID)
	case EventOnRange:
		return fmt.Sprintf(
			"[%s] The competitor(%d) is on the firing range(%s)",
			eventTime,
			e.CompetitorID,
			e.ExtraParams[0],
		)
	case EventHit:
		return fmt.Sprintf(
			"[%s] The target(%s) has been hit by competitor(%d)",
			eventTime,
			e.ExtraParams[0],
			e.CompetitorID,
		)
	case EventLeftRange:
		return fmt.Sprintf(
			"[%s] The competitor(%d) left the firing range",
			eventTime,
			e.CompetitorID,
		)
	case EventEnteredPenalty:
		return fmt.Sprintf(
			"[%s] The competitor(%d) entered the penalty laps",
			eventTime,
			e.CompetitorID,
		)
	case EventLeftPenalty:
		return fmt.Sprintf(
			"[%s] The competitor(%d) left the penalty laps",
			eventTime,
			e.CompetitorID,
		)
	case EventEndedMainLap:
		return fmt.Sprintf("[%s] The competitor(%d) ended the main lap", eventTime, e.CompetitorID)
	case EventCantContinue:
		return fmt.Sprintf(
			"[%s] The competitor(%d) can`t continue: %s",
			eventTime,
			e.CompetitorID,
			strings.Join(e.ExtraParams, " "),
		)
	case EventDisqualified:
		return fmt.Sprintf("[%s] The competitor(%d) is disqualified", eventTime, e.CompetitorID)
	case EventFinished:
		return fmt.Sprintf("[%s] The competitor(%d) finished the race", eventTime, e.CompetitorID)
	default:
		return fmt.Sprintf("Unknown eventID: %d", e.ID)
	}
}

func parseEvents(eventStream io.Reader) <-chan Event {
	eventCh := make(chan Event, bufferSize)
	go func() {
		defer close(eventCh)
		scanner := bufio.NewScanner(eventStream)
		for scanner.Scan() {
			event, err := parseEvent(scanner.Text())
			if err != nil {
				log.Printf("error parsing event: %s", err)
				continue
			}
			eventCh <- event
		}
		if err := scanner.Err(); err != nil {
			log.Printf("error reading event stream: %s\n", err)
		}
	}()

	return eventCh
}

func parseEvent(line string) (Event, error) {
	split := strings.Split(line, " ")
	if len(split) < 3 {
		return Event{}, fmt.Errorf(errInvalidLineFormat, line)
	}
	timeStr := split[0][1 : len(split[0])-1]

	eventTime, err := time.Parse(timeLayout, timeStr)
	if err != nil {
		return Event{}, fmt.Errorf(errCannotParseTime, err)
	}

	eventID, err := strconv.Atoi(split[1])
	if err != nil {
		return Event{}, fmt.Errorf(errCannotConvert, "eventID", err)
	}

	competitorID, err := strconv.Atoi(split[2])
	if err != nil {
		return Event{}, fmt.Errorf(errCannotConvert, "competitorID", err)
	}

	return Event{
		ID:           eventID,
		CompetitorID: competitorID,
		Time:         eventTime,
		ExtraParams:  split[3:],
	}, nil
}

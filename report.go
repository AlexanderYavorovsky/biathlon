package main

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type ResultInfo struct {
	Cfg         Config
	Comp        Competitor
	Status      string
	TotalShots  int
	MissedShots int
}

func getResult(comp Competitor, cfg Config) ResultInfo {
	totalShots := cfg.FiringLines * TargetsPerFiringLine

	var status string
	if comp.Status == StatusFinished {
		status = fmtDuration(comp.TotalTime)
	} else {
		status = comp.Status
	}

	return ResultInfo{
		Cfg:         cfg,
		Comp:        comp,
		Status:      status,
		TotalShots:  totalShots,
		MissedShots: totalShots - comp.Hits,
	}
}

func (r ResultInfo) String() string {
	return fmt.Sprintf("[%s] %d [%s] {%s} %d/%d",
		r.Status,
		r.Comp.ID,
		getLapsStr(r),
		getPenaltyStr(r),
		r.Comp.Hits,
		r.TotalShots,
	)
}

func getLapsStr(r ResultInfo) string {
	var lapsSB strings.Builder
	for i := 0; i < r.Cfg.Laps; i++ {
		if i >= len(r.Comp.Laps) {
			lapsSB.WriteString("{,}")
			continue
		}
		l := r.Comp.Laps[i]
		lapsSB.WriteString(fmt.Sprintf("{%s, %.3f}", fmtDuration(l.Duration), l.Speed))
		if i < r.Cfg.Laps-1 {
			lapsSB.WriteString(", ")
		}
	}

	return lapsSB.String()
}

func getPenaltyStr(r ResultInfo) string {
	if r.Comp.Penalty.Seconds() == 0 {
		return ","
	}

	totalPenaltyLen := float64(r.Cfg.PenaltyLen) * float64(r.MissedShots)
	penaltySpeed := totalPenaltyLen / float64(r.Comp.Penalty.Seconds())
	return fmt.Sprintf("%s, %.3f", fmtDuration(r.Comp.Penalty), penaltySpeed)
}

func fmtDuration(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	ms := d.Milliseconds() % 1000

	return fmt.Sprintf("%02d:%02d:%02d.%03d", h, m, s, ms)
}

func getList(competitors map[int]*Competitor) []Competitor {
	s := make([]Competitor, 0, len(competitors))
	for _, c := range competitors {
		s = append(s, *c)
	}
	return s
}

// Sort in following order:
// 1. List of finished, by ascending TotalTime;
// 2. List of NotFinished and NotStarted, by ascending ID.
func getSortedByTime(c []Competitor) []Competitor {
	sort.Slice(c, func(i, j int) bool {
		ci := c[i]
		cj := c[j]

		if ci.Status == StatusFinished {
			if cj.Status != ci.Status {
				return true
			}
			return ci.TotalTime < cj.TotalTime
		}

		if cj.Status == StatusFinished {
			return false
		}

		return ci.ID < cj.ID
	})

	return c
}

func printFinalReport(competitors []Competitor, cfg Config) {
	report := generateFinalReport(competitors, cfg)
	fmt.Println("\nFinal Report:")
	fmt.Println(report)
}

func generateFinalReport(competitors []Competitor, cfg Config) string {
	var sb strings.Builder
	for _, c := range competitors {
		sb.WriteString(getResult(c, cfg).String())
		sb.WriteString("\n")
	}
	return sb.String()
}

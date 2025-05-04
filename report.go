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
	TotalShots  int
	MissedShots int
}

func getResult(cfg Config, comp Competitor) ResultInfo {
	totalShots := cfg.FiringLines * TargetsPerFiringLine
	return ResultInfo{
		Cfg:         cfg,
		Comp:        comp,
		TotalShots:  totalShots,
		MissedShots: totalShots - comp.Hits,
	}
}

func (r ResultInfo) String() string {
	return fmt.Sprintf("[%s] %d [%s] {%s} %d/%d",
		r.Comp.Status,
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
	var penaltySB strings.Builder
	if r.Comp.Penalty.Seconds() == 0 {
		penaltySB.WriteString(",")
	} else {
		totalPenaltyLen := float64(r.Cfg.PenaltyLen) * float64(r.MissedShots)
		penaltySpeed := totalPenaltyLen / float64(r.Comp.Penalty.Seconds())
		penaltySB.WriteString(fmt.Sprintf("%s, %.3f", fmtDuration(r.Comp.Penalty), penaltySpeed))
	}

	return penaltySB.String()
}

func fmtDuration(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	ms := d.Milliseconds() % 1000

	return fmt.Sprintf("%02d:%02d:%02d.%03d", h, m, s, ms)
}

func calculateCompetitor(c *Competitor) {
	if c.Status == StatusNotStarted || c.Status == StatusNotFinished {
		return
	}

	for _, l := range c.Laps {
		c.TotalTime += l.Duration
	}

	c.Status = fmtDuration(c.TotalTime)
}

func calculateCompetitors(competitors map[int]*Competitor) {
	for _, c := range competitors {
		// TODO: use goroutines?
		calculateCompetitor(c)
	}
}

func getSortedByTime(competitors map[int]*Competitor) []Competitor {
	s := make([]Competitor, 0)
	//TODO: use counter instead to prevent reallocation?
	for _, c := range competitors {
		s = append(s, *c)
	}

	sort.Slice(s, func(i, j int) bool {
		return s[i].TotalTime < s[j].TotalTime
	})

	return s
}

func printFinalReport(competitors []Competitor, cfg Config) {
	//TODO: handle cases NotFinished, NotStarted -- when to print?
	for _, c := range competitors {
		r := getResult(cfg, c)
		fmt.Println(r)
	}
}

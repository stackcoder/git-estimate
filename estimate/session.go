package estimate

import (
	"time"
)

type WorkingSession struct {
	Baseline float64
}

// sums up the hours in a given day, assuming the beginning of the session 2 hours earlier than the first commit
// and the end of the job in the last commit of the day
func (ws WorkingSession) Estimate(byAuthors map[string][]time.Time) []Result {
	results := make([]Result, len(byAuthors))
	c := 0
	for k, _ := range byAuthors {
		r := &results[c]
		r.Author = k
		next := time.Time{}
		v := byAuthors[k]
		for _, t := range v {
			if next.IsZero() {
				next = t
				continue
			}

			diff := next.Sub(t).Hours()
			if diff < 8 {
				r.Hours += diff
			} else {
				r.Hours += (time.Duration(ws.Baseline) * time.Hour).Hours()
			}
			next = t
		}

		// we have a single commit from the author: add the default blunt estimate
		if len(v) == 1 {
			r.Hours += (time.Duration(ws.Baseline) * time.Hour).Hours()
		}

		r.Days = r.Hours / 8.0
		c++
	}
	return results
}
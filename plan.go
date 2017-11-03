package main

import (
	"fmt"
	"strings"
	"time"
)

type plan []*stop

func (p plan) copy() plan {
	res := plan{}
	for _, s := range p {
		res = append(res, s.copy())
	}
	return res
}

func (p *plan) add(start int, startTime time.Time, t *trip) {
	if start == t.start {
		// Can the trip start here? Yes, if the car is going in the right
		// direction and hasn't left yet.
		if p.dir(start, 0) == t.dir() && startTime.After(time.Now()) {
			(*p)[0].trips = append((*p)[0].trips, t)
			p.addEnd(t, 0)
			return
		}
	}

	prev := start
	for i, s := range *p {
		if s.floor == t.start {
			// Can the trip start here? Yes, if the car is going in the
			// right direction, OR if it's finished going in the wrong
			// direction.
			if p.dir(start, i) == t.dir() || p.dir(start, i+1) != -t.dir() {
				s.trips = append(s.trips, t)
				p.addEnd(t, i)
				return
			}
		}

		// Can we insert a stop before the current one? Yes, if the car is
		// going in the right direction and hasn't left the previous stop
		// yet. TODO: it's OK if it's left the previous stop as long as
		// it's still "far enough" ahead of the proposed new stop.
		if ordered(prev, t.start, s.floor, t.dir()) && (i > 0 || startTime.After(time.Now())) {
			pre := (*p)[:i]
			post := (*p)[i:]
			*p = append(plan{}, pre...)
			*p = append(*p, &stop{t.start, []*trip{t}})
			*p = append(*p, post...)
			p.addEnd(t, i)
			return
		}

		prev = s.floor
	}

	if t.start != start {
		*p = append(*p, &stop{t.start, []*trip{t}})
	}
	*p = append(*p, &stop{t.end, []*trip{t}})
}

func (p *plan) addEnd(t *trip, startIndex int) {
	prev := (*p)[startIndex].floor
	for i := startIndex + 1; i < len(*p); i++ {
		s := (*p)[i]
		if s.floor == t.end {
			s.trips = append(s.trips, t)
			return
		}
		if ordered(prev, t.end, s.floor, t.dir()) || p.dir(0, i+1) != t.dir() {
			pre := (*p)[:i]
			post := (*p)[i:]
			*p = append(plan{}, pre...)
			*p = append(*p, &stop{t.end, []*trip{t}})
			*p = append(*p, post...)
			return
		}
		prev = s.floor
	}
	*p = append(*p, &stop{t.end, []*trip{t}})
}

func (p *plan) dir(start, index int) int {
	if index >= len(*p) {
		return 0
	}
	if index == 0 {
		return dir(start, (*p)[0].floor)
	}
	return dir((*p)[index-1].floor, (*p)[index].floor)
}

func (p *plan) eval(start int, startTime time.Time) (avg, max time.Duration) {
	var (
		trips    = make(map[*trip]time.Duration)
		prev     = start
		elapsed  time.Duration
		totalDur time.Duration
		maxDur   time.Duration
	)
	for _, s := range *p {
		floors := abs(s.floor - prev)
		elapsed += time.Duration(int64(floors) * int64(floorDur))
		arriveTime := startTime.Add(elapsed)
		for _, t := range s.trips {
			if s.floor == t.end {
				tripDur := arriveTime.Sub(t.reqTime)
				trips[t] = tripDur
				totalDur += tripDur
				if tripDur > maxDur {
					maxDur = tripDur
				}
			}
		}
		prev = s.floor
	}

	return time.Duration(int64(totalDur) / int64(len(trips))), maxDur
}

func (p plan) String() string {
	var floorStrs []string
	for _, s := range p {
		floorStrs = append(floorStrs, fmt.Sprintf("%d", s.floor))
	}
	return fmt.Sprintf("[%s]", strings.Join(floorStrs, " "))
}

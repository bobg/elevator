package main

import (
	"fmt"
	"log"
	"time"
)

type car struct {
	letter byte

	plan plan

	trips      chan *trip // new-trip commands arrive on this channel
	eventTimer *time.Timer
	nextEvent  event

	lastStop     int       // current or latest floor stopped at
	lastStopTime time.Time // time at which the car left, or will leave, the laststop
}

func newCar(letter byte) *car {
	c := &car{
		letter:       letter,
		trips:        make(chan *trip),
		eventTimer:   time.NewTimer(time.Hour),
		lastStop:     1,
	}
	c.eventTimer.Stop()
	return c
}

func (c *car) run() {
	for {
		select {
		case t := <-c.trips:
			c.logf("adding trip %s", t)
			c.add(t)

		case <-c.eventTimer.C:
			switch c.nextEvent {
			case doorClose:
				c.logf("door closed")
				floors := abs(c.plan[0].floor - c.lastStop)
				c.nextEvent = arrive
				c.eventTimer.Reset(time.Duration(int64(floors) * int64(floorDur)))

			case arrive:
				c.logf("arrived at %d", c.plan[0].floor)
				c.lastStop = c.plan[0].floor
				c.plan = c.plan[1:]
				if len(c.plan) > 0 {
					c.nextEvent = doorClose
					c.eventTimer.Reset(doorDur)
					c.lastStopTime = time.Now().Add(doorDur)
				}

			default:
				// xxx error!
			}
		}
	}
}

func (c *car) add(t *trip) {
	var (
		wasIdle  = len(c.plan) == 0
		lastStopTime = c.lastStopTime
		nextStop int
	)
	if wasIdle {
		lastStopTime = time.Now().Add(doorDur)
	} else {
		nextStop = c.plan[0].floor
	}
	c.plan.add(c.lastStop, lastStopTime, t)
	if wasIdle {
		c.lastStopTime = lastStopTime
		c.nextEvent = doorClose
		c.eventTimer.Reset(doorDur)
	} else if nextStop != c.plan[0].floor {
		// there's a new nextStop
		floors := abs(c.plan[0].floor - c.lastStop)
		arriveTime := lastStopTime.Add(time.Duration(int64(floors) * int64(floorDur)))
		c.eventTimer.Reset(arriveTime.Sub(time.Now()))
	}
}

func (c *car) logf(format string, args ...interface{}) {
	format = fmt.Sprintf("car %c: %s %s", c.letter, format, c.plan)
	log.Printf(format, args...)
}

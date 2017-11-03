package main

import (
	"fmt"
	"time"
)

type trip struct {
	start, end int
	reqTime    time.Time
}

func (t *trip) dir() int {
	return dir(t.start, t.end)
}

func (t *trip) String() string {
	return fmt.Sprintf("<%d-%d>", t.start, t.end)
}

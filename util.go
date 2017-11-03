package main

import "time"

type event int

const (
	doorClose event = 1 + iota
	arrive
)

var (
	floorDur time.Duration = 2 * time.Second
	doorDur  time.Duration = 10 * time.Second
)

// tells whether a, b, c are ordered in the direction d
func ordered(a, b, c, d int) bool {
	if d > 0 {
		return c > b && b > a
	}
	if d < 0 {
		return a > b && b > c
	}
	return false
}

func dir(start, end int) int {
	switch {
	case start < end:
		return 1
	case start > end:
		return -1
	default:
		return 0
	}
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

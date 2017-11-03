package main

type stop struct {
	floor int
	trips []*trip
}

func (s *stop) copy() *stop {
	return &stop{s.floor, append([]*trip{}, s.trips...)}
}

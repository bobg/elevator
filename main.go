package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	var cars [5]*car
	for i := 0; i < len(cars); i++ {
		cars[i] = newCar('A' + uint8(i))
		go cars[i].run()
	}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		scanner2 := bufio.NewScanner(strings.NewReader(line))
		scanner2.Split(bufio.ScanWords)
		if !scanner2.Scan() {
			continue // xxx log error
		}
		aStr := scanner2.Text()
		a, err := strconv.Atoi(aStr)
		if err != nil {
			continue // xxx log error
		}
		if !scanner2.Scan() {
			continue // xxx log error
		}
		bStr := scanner2.Text()
		b, err := strconv.Atoi(bStr)
		if err != nil {
			continue // xxx log error
		}
		// xxx range-check a and b

		t := &trip{a, b, time.Now()}

		var (
			bestCar *car
			bestAvg time.Duration
		)
		for _, c := range cars {
			var (
				p            = c.plan.copy()
				lastStopTime = c.lastStopTime
			)
			if lastStopTime == (time.Time{}) {
				lastStopTime = time.Now().Add(doorDur)
			}
			p.add(c.lastStop, c.lastStopTime, t)
			avg, _ := p.eval(c.lastStop, c.lastStopTime)
			if bestCar == nil || avg < bestAvg {
				bestCar = c
				bestAvg = avg
			}
		}
		bestCar.trips <- t
	}
}

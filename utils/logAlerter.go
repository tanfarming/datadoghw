package utils

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

//LogAlerter receives logging events and fires alerts as per configured
type LogAlerter struct {
	params LogAlerterParams

	hitsQueue []int
	isFiring  bool
}

//LogAlerterParams contains configs for LogAlerter
type LogAlerterParams struct {
	TrackInterval        int
	HighTrafficThreshold int
}

//NewLogAlerter constructs a new LogAlerter
func NewLogAlerter(params LogAlerterParams) *LogAlerter {
	return &LogAlerter{
		params:   params,
		isFiring: false}
}

//hit responses to new logging event
func (la *LogAlerter) hit(unixTime int) {

	//enqueue new hit
	now := unixTime
	la.hitsQueue = append(la.hitsQueue, now)
	hits := len(la.hitsQueue)
	if hits == 1 {
		return
	}

	//calculate hitRate for events cached in queue
	timeSpan := now - la.hitsQueue[0]
	if timeSpan == 0 {
		timeSpan = 1
	}
	hitRate := hits / timeSpan

	//high traffic alert
	if !la.isFiring && hitRate > la.params.HighTrafficThreshold {
		la.fire(hitRate, now)
	}

	//low traffic, issue recover alert if currently firing
	if la.isFiring && int(hitRate) < la.params.HighTrafficThreshold {
		la.recover(now)
	}

	//drop expired events
	i := 0
	for timeSpan > la.params.TrackInterval {
		timeSpan = now - la.hitsQueue[i]
		i++
	}
	la.hitsQueue = la.hitsQueue[i:]

}

func (la *LogAlerter) fire(hitRate, timeStamp int) {
	msg := fmt.Sprintf(
		"+ High traffic generated an alert with avg. hit rate = {%d}, triggered at {%d} +",
		hitRate, timeStamp)
	hr := strings.Repeat("+", utf8.RuneCountInString(msg))
	fmt.Println(fmt.Sprintf("\n%s\n%s\n%s\n", hr, msg, hr))

	la.isFiring = true
}

func (la *LogAlerter) recover(timeStamp int) {
	msg := fmt.Sprintf("- High traffic alert recovered at %d -", timeStamp)
	hr := strings.Repeat("-", utf8.RuneCountInString(msg))
	fmt.Println(fmt.Sprintf("\n%v\n%v\n%v\n", hr, msg, hr))

	la.isFiring = false
}

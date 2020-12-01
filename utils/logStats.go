package utils

import (
	"fmt"
)

//logStats caches and manages data for logWorker
type logStats struct {
	httpRes               map[int]int
	hostHits              map[string]*httpHits
	sectionHits           map[string]*httpHits
	sectionWithMostHits   string
	throughput            *httpThroughput
	firstEntryUnixTime    int
	previousEntryUnixTime int
	timeSpanSec           int
}

func newLogStats() logStats {
	logStats := logStats{
		httpRes:               make(map[int]int),
		hostHits:              make(map[string]*httpHits),
		sectionHits:           make(map[string]*httpHits),
		sectionWithMostHits:   "",
		throughput:            newHttpThroughput(),
		firstEntryUnixTime:    -1,
		previousEntryUnixTime: -1,
		timeSpanSec:           0}
	return logStats
}

func (ls *logStats) ToString() string {

	httpResStr := "\noverall http responses: "
	for k, v := range ls.httpRes {
		httpResStr += fmt.Sprintf("%d x http%d, ", v, k)
	}

	hostHitsStr := "\n"
	for k, v := range ls.hostHits {
		hostHitsStr += fmt.Sprintf("\nhost: %s: total hits: %d \n", k, v.total)
		hostHitsStr += "\t methods: "
		for k, v := range v.byMethod {
			hostHitsStr += fmt.Sprintf("%d x %s, ", v, k)
		}
		hostHitsStr += "\n\t responses: "
		for k, v := range v.byHttpRes {
			hostHitsStr += fmt.Sprintf("%d x http%d, ", v, k)
		}
	}

	sectionHitsStr := "\n"
	for k, v := range ls.sectionHits {
		sectionHitsStr += fmt.Sprintf("\nsection: %s: total hits: %d \n", k, v.total)
		sectionHitsStr += "\t methods: "
		for k, v := range v.byMethod {
			sectionHitsStr += fmt.Sprintf("%d x %s, ", v, k)
		}
		sectionHitsStr += "\n\t responses: "
		for k, v := range v.byHttpRes {
			sectionHitsStr += fmt.Sprintf("%d x http%d, ", v, k)
		}
	}

	throughputStr := "\n"
	throughputStr += fmt.Sprintf("\ntotal throughput: %d bytes\n", ls.throughput.total)
	for k, v := range ls.throughput.byHost {
		throughputStr += fmt.Sprintf("\t%s: %d bytes\n", k, v)
	}

	return httpResStr + hostHitsStr + sectionHitsStr + throughputStr
}

//logStats.add update itself with new logData
func (ls *logStats) add(logData logData) {
	//needed for timeSpan calculation
	if ls.firstEntryUnixTime == -1 {
		ls.firstEntryUnixTime = logData.date
	}
	//cheaper than sort / heap / bst for tracking only the max
	if ls.sectionWithMostHits == "" {
		ls.sectionWithMostHits = logData.section
	}
	//update httpRes
	if _, ok := ls.httpRes[logData.resCode]; ok {
		ls.httpRes[logData.resCode]++
	} else {
		ls.httpRes[logData.resCode] = 1
	}
	//update hostHits
	if _, ok := ls.hostHits[logData.host]; ok {
		ls.hostHits[logData.host].add(logData)
	} else {
		ls.hostHits[logData.host] = newHttpHits()
		ls.hostHits[logData.host].add(logData)
	}
	//update sectionHits
	if _, ok := ls.sectionHits[logData.section]; ok {
		ls.sectionHits[logData.section].add(logData)
	} else {
		ls.sectionHits[logData.section] = newHttpHits()
		ls.sectionHits[logData.section].add(logData)
		if ls.sectionHits[logData.section].total > ls.sectionHits[ls.sectionWithMostHits].total {
			ls.sectionWithMostHits = logData.section
		}
	}
	//update throughput
	ls.throughput.add(logData)

	ls.timeSpanSec = logData.date - ls.firstEntryUnixTime
	ls.previousEntryUnixTime = logData.date
}

//httpThroughput manages throughput stats
type httpThroughput struct {
	total  int
	byHost map[string]int
}

func newHttpThroughput() *httpThroughput {
	return &httpThroughput{
		total:  0,
		byHost: make(map[string]int)}
}

func (tp *httpThroughput) add(logData logData) {
	tp.total += logData.bytes
	if _, ok := tp.byHost[logData.host]; ok {
		tp.byHost[logData.host] += logData.bytes
	} else {
		tp.byHost[logData.host] = logData.bytes
	}
}

//httpThroughput manages http hit stats
type httpHits struct {
	total     int
	byMethod  map[string]int
	byHttpRes map[int]int
}

func newHttpHits() *httpHits {
	return &httpHits{
		total:     0,
		byMethod:  make(map[string]int),
		byHttpRes: make(map[int]int)}
}

func (hh *httpHits) add(logData logData) {
	hh.total++

	if _, ok := hh.byMethod[logData.method]; ok {
		hh.byMethod[logData.method]++
	} else {
		hh.byMethod[logData.method] = 1
	}
	if _, ok := hh.byHttpRes[logData.resCode]; ok {
		hh.byHttpRes[logData.resCode]++
	} else {
		hh.byHttpRes[logData.resCode] = 1
	}
}

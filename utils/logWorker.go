package utils

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
)

//LogWorker process log, produce reports and update logAlerter
type LogWorker struct {
	params  LogWorkerParams
	alerter *LogAlerter

	LogReceiver chan string
	logStats    logStats
}

//LogWorkerParams contains configs for LogWorker
type LogWorkerParams struct {
	ReportInterval int
	FlushInterval  int
}

//NewLogWorker constructs a new LogWorker
func NewLogWorker(params LogWorkerParams, alerter *LogAlerter) *LogWorker {

	lw := LogWorker{
		params:      params,
		alerter:     alerter,
		LogReceiver: make(chan string),
		logStats:    newLogStats()}
	return &lw
}

//Go starts the LogWorker's listening loop, shipping incomming logs to processLogEntry() for as long as the channle's open
func (lw *LogWorker) Go(waitGroup *sync.WaitGroup) {
	for {
		log, open := <-lw.LogReceiver
		if !open {
			fmt.Println("---LogWorker: job done")
			waitGroup.Done()
			return
		}
		lw.processLogEntry(log)
	}
}

//processLogEntry does 4 tasks:
//1. parses the logEntry string into logData
//2. process the logData into logStats
//3. report the logStats at configured interval
//4. update the LogAlerter
func (lw *LogWorker) processLogEntry(logEntry string) {
	logData, err := newLogData(logEntry)
	if err != nil {
		fmt.Println("processLogEntry: skipping bad logEntry, err: " + err.Error())
	}

	lw.logStats.add(logData)

	if lw.logStats.timeSpanSec >= lw.params.ReportInterval {
		lw.report()
		lw.logStats = newLogStats()
	}

	lw.alerter.hit(logData.date)

}

func (lw *LogWorker) report() {
	msg := "\n---logWorker report---\n"
	msg += "most hit section:  " + lw.logStats.sectionWithMostHits +
		" with " + strconv.Itoa(lw.logStats.sectionHits[lw.logStats.sectionWithMostHits].total) + " hits\n"
	msg += "\nother might be useful stats: \n" + lw.logStats.ToString()
	fmt.Println(msg)
}

//logData produces logData for processLogEntry
type logData struct {
	host    string
	ident   string
	user    string
	date    int
	method  string
	section string
	resCode int
	bytes   int
}

func newLogData(logEntry string) (logData, error) {
	dataArr := strings.Split(logEntry, ",")

	host := strings.Trim(dataArr[0], `"`)
	ident := strings.Trim(dataArr[1], `"`)
	user := strings.Trim(dataArr[2], `"`)
	date, err := strconv.Atoi(dataArr[3])
	if err != nil {
		return logData{}, err
	}
	methodNsection := strings.Split(dataArr[4], `/`)
	method := strings.Trim(methodNsection[0], ` `)
	method = strings.Trim(method, `"`)

	resCode, err := strconv.Atoi(dataArr[5])
	if err != nil {
		return logData{}, err
	}
	bytes, err := strconv.Atoi(dataArr[6])
	if err != nil {
		return logData{}, err
	}
	return logData{
		host:    host,
		ident:   ident,
		user:    user,
		date:    date,
		method:  method,
		section: `/` + methodNsection[1],
		resCode: resCode,
		bytes:   bytes,
	}, nil
}

package processors

import (
	"bufio"
	"errors"
	"main/utils"
	"os"
	"strings"
	"sync"
)

func ProcessLogfile(file string) error {

	//read log file
	f, err := os.Open(file)
	if err != nil {
		return errors.New("failed to open file: " + err.Error())
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Scan()

	//validate log file
	err = isFileHeaderValid(scanner.Text())
	if err != nil {
		return err
	}

	//create logWorker
	alerter := utils.NewLogAlerter(
		utils.LogAlerterParams{
			TrackInterval:        120,
			HighTrafficThreshold: 10})
	logWorkerParams := utils.LogWorkerParams{
		ReportInterval: 10,
		FlushInterval:  -1}

	logWorker := utils.NewLogWorker(logWorkerParams, alerter)

	//start logWorker
	var wg sync.WaitGroup
	wg.Add(1)
	go logWorker.Go(&wg)

	//pushing logs to logWorker
	for scanner.Scan() {
		logWorker.LogReceiver <- scanner.Text()
	}

	close(logWorker.LogReceiver)
	wg.Wait()

	return nil
}

func isFileHeaderValid(logEntry string) error {
	dataArr := strings.Split(logEntry, ",")
	if len(dataArr) == 7 && dataArr[0] == `"remotehost"` {
		return nil
	}
	return errors.New("unexpected log file header")
}

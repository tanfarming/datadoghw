package main

import (
	"flag"
	"fmt"

	"main/processors"
)

func main() {

	logFile := flag.String("file", "", "path to the log file")
	flag.Parse()

	if *logFile != "" {
		err := processors.ProcessLogfile(*logFile)
		if err != nil {
			fmt.Println("ERR: " + err.Error())
		}
	} else {
		fmt.Println("ERR: missing cli argument:")
		flag.PrintDefaults()
	}

}

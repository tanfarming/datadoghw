package processors

import (
	"testing"
)

//TestProcessLogfile tests processors.ProcessLogfile
//requires valid log file ./testLog.txt
func TestProcessLogfile(t *testing.T) {
	err := ProcessLogfile("./testLog.txt")
	if err != nil {
		t.Errorf("ProcessLogfile failed, err" + err.Error())
	}
}

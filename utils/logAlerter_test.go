package utils

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

//TestLogAlerter tests processors.ProcessLogfile
func TestLogAlerter(t *testing.T) {

	r, w, _ := os.Pipe()
	os.Stdout = w

	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	alerter := NewLogAlerter(LogAlerterParams{TrackInterval: 1, HighTrafficThreshold: 10})

	//exceeded HighTrafficThreshold: should fire -- stdout log should contain "triggered at {1}"
	for i := 0; i < alerter.params.HighTrafficThreshold+1; i++ {
		alerter.hit(1)
	}
	fire := alerter.isFiring
	if !fire {
		t.Errorf("LogAlerter failed to fire ")
	}
	//should not fire duplicated alerts -- stdout log should NOT contain "triggered at {2}"
	alerter.hit(2)

	//should recover -- stdout should contain "recovered at 3"
	alerter.hit(3)
	fire = alerter.isFiring
	if fire {
		t.Errorf("LogAlerter failed to recover ")
	}
	w.Close()
	stdout := <-outC

	if !strings.Contains(stdout, "triggered at {1}") {
		t.Errorf("LogAlerter failed to output alert message ")
	}

	if strings.Contains(stdout, "triggered at {2}") {
		t.Errorf("LogAlerter failed to avoid duplicated alerts")
	}

	if !strings.Contains(stdout, "recovered at 3") {
		t.Errorf("LogAlerter failed to output recover message")
	}

}

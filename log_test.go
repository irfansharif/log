package log

import (
	"fmt"
	"os"
	"testing"
)

func TestSetGetPCMode(t *testing.T) {
	resetgstate()
	defer func() { resetgstate() }()

	f, l := getcaller()
	pc := fmt.Sprintf("%s:%d", f, l)
	setPCMode(pc, InfoMode)
	lmode, ok := getPCMode(pc)
	if !ok {
		t.Errorf("Expected PC mode for %s to be found", pc)
	}
	if lmode != InfoMode {
		t.Errorf("Expected log.InfoLevel, got %s", lmode)
	}
}

func TestGetPCMode(t *testing.T) {
	resetgstate()
	defer func() { resetgstate() }()

	f, l := getcaller()
	pc := fmt.Sprintf("%s:%d", f, l)
	_, ok := getPCMode(pc)
	if ok {
		t.Errorf("Didn't expected PC mode for %s to be found", pc)
	}
}

func TestLog(t *testing.T) {
	resetgstate()
	defer func() { resetgstate() }()

	logger := New(os.Stdout)
	logger.Info("\tlogger: TESTING LOGGER OUTPUT\n")
}

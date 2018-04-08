package log

import (
	"bytes"
	"fmt"
	"regexp"
	"testing"
)

func TestSetGetGlobalPCMode(t *testing.T) {
	resetgstate()
	defer func() { resetgstate() }()

	tp := fmt.Sprintf("%s:%d", "t.go", 42)
	EnableTracePoint(tp)
	enabled := CheckTracePoint(tp)
	if !enabled {
		t.Errorf("Expected tracepoint %s to be enabled", tp)
	}
}

func TestGetGlobalPCMode(t *testing.T) {
	resetgstate()
	defer func() { resetgstate() }()

	tp := fmt.Sprintf("%s:%d", "t.go", 42)
	enabled := CheckTracePoint(tp)
	if enabled {
		t.Errorf("Didn't expected tracepoint mode for %s to be enabled", tp)
	}
}

func TestInfoLog(t *testing.T) {
	resetgstate()
	defer func() { resetgstate() }()

	SetGlobalLogMode(InfoMode)

	buffer := new(bytes.Buffer)
	logger := New(buffer)
	{
		logger.Info("info")
		regex := "I [\\w]+.go:[\\d]+: info"
		match, err := regexp.Match(regex, buffer.Bytes())
		if err != nil {
			t.Error(err)
		}
		if !match {
			t.Errorf("expected pattern: \"%s\", got: %s", regex, buffer.String())
		}
		buffer.Reset()
	}
	{
		logger.Infof("infof")
		regex := "I [\\w]+.go:[\\d]+: infof"
		match, err := regexp.Match(regex, buffer.Bytes())
		if err != nil {
			t.Error(err)
		}
		if !match {
			t.Errorf("expected pattern: \"%s\", got: %s", regex, buffer.String())
		}
		buffer.Reset()
	}
	{
		logger.Infof("%t %d %s", true, 1, "infof")
		regex := "I [\\w]+.go:[0-9]+: true [\\d]+ infof"
		match, err := regexp.Match(regex, buffer.Bytes())
		if err != nil {
			t.Error(err)
		}
		if !match {
			t.Errorf("expected pattern: \"%s\", got: %s", regex, buffer.String())
		}
		buffer.Reset()
	}
}

func TestDebugModeEnableDisable(t *testing.T) {
	resetgstate()
	defer func() { resetgstate() }()

	SetGlobalLogMode(InfoMode | DebugMode)

	buffer := new(bytes.Buffer)
	logger := New(buffer)
	{
		logger.Debug("debug")
		logger.Debugf("%t %d %s", true, 1, "debugf")
		logger.Debugf("debugf")

		regex := ""
		match, err := regexp.Match(regex, buffer.Bytes())
		if err != nil {
			t.Error(err)
		}
		if !match {
			t.Errorf("expected pattern: \"%s\", got: %s", regex, buffer.String())
		}
		buffer.Reset()
	}
	{
		logger.Debug("debug")
		regex := "D [\\w]+.go:[\\d]+: debug"
		match, err := regexp.Match(regex, buffer.Bytes())
		if err != nil {
			t.Error(err)
		}
		if !match {
			t.Errorf("expected pattern: \"%s\", got: %s", regex, buffer.String())
		}
		buffer.Reset()
	}
}

func TestEnableTracePoint(t *testing.T) {
	resetgstate()
	defer func() { resetgstate() }()

	SetGlobalLogMode(DisabledMode)

	// XXX(irfansharif): This test depends on the exact difference in line
	// numbers between the call to callers and the logger.Info execution below.
	// The tracepoint is set to be the line exactly ten lines below it.

	file, line := caller(0)
	tp := fmt.Sprintf("%s:%d", file, line+10)
	EnableTracePoint(tp)
	if tpenabled := CheckTracePoint(tp); !tpenabled {
		t.Error("Expected tracepoint %s to be enabled; found disabled", tp)
	}

	buffer := new(bytes.Buffer)
	logger := New(buffer)
	{
		logger.Info()
		if buffer.Len() == 0 {
			t.Error("Expected stack trace to be populated, found empty buffer instead")
		}
		line, err := buffer.ReadString(byte('\n'))
		if err != nil {
			t.Error(err)
		}
		regex := "goroutine [\\d]+ \\[running\\]:"
		match, err := regexp.Match(regex, []byte(line))
		if err != nil {
			t.Error(err)
		}
		if !match {
			t.Errorf("expected pattern (first line of stack trace): \"%s\", got: %s", regex, line)
		}
	}
}

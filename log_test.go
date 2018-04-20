package log

import (
	"bytes"
	"fmt"
	"path/filepath"
	"regexp"
	"testing"
)

func TestSetGetGlobalPCMode(t *testing.T) {
	resetgstate()
	defer func() { resetgstate() }()

	tp := fmt.Sprintf("%s:%d", "t.go", 42)
	SetTracePoint(tp)
	enabled := GetTracePoint(tp)
	if !enabled {
		t.Errorf("Expected tracepoint %s to be enabled", tp)
	}
}

func TestGetGlobalPCMode(t *testing.T) {
	resetgstate()
	defer func() { resetgstate() }()

	tp := fmt.Sprintf("%s:%d", "t.go", 42)
	enabled := GetTracePoint(tp)
	if enabled {
		t.Errorf("Didn't expected tracepoint mode for %s to be enabled", tp)
	}
}

func TestInfoLog(t *testing.T) {
	resetgstate()
	defer func() { resetgstate() }()

	SetGlobalLogMode(InfoMode)

	buffer := new(bytes.Buffer)
	logger := New(Writer(buffer))
	{
		logger.Info("info")
		regex := "^I.*: info"
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
		regex := "^I.*: infof"
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
		regex := "^I.*: true 1 infof"
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
	logger := New(Writer(buffer))
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
		regex := "^D.*: debug"
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
	tp := fmt.Sprintf("%s:%d", filepath.Base(file), line+10)
	SetTracePoint(tp)
	if tpenabled := GetTracePoint(tp); !tpenabled {
		t.Errorf("Expected tracepoint %s to be enabled; found disabled", tp)
	}

	buffer := new(bytes.Buffer)
	logger := New(Writer(buffer))
	{
		logger.Info()
		if buffer.Len() == 0 {
			t.Error("Expected stack trace to be populated, found empty buffer instead")
		}

		line, err := buffer.ReadString(byte('\n'))
		if err != nil {
			t.Error(err)
		}

		goroutineRegex := "^goroutine [\\d]+ \\[running\\]:"
		match, err := regexp.Match(goroutineRegex, []byte(line))
		if err != nil {
			t.Error(err)
		}
		if !match {
			t.Errorf("expected pattern (first line): \"%s\", got: %s", goroutineRegex, line)
		}

		line, err = buffer.ReadString(byte('\n'))
		if err != nil {
			t.Error(err)
		}

		functionSignatureRegex := "^github.com/irfansharif/log.TestEnableTracePoint"
		match, err = regexp.Match(functionSignatureRegex, []byte(line))
		if err != nil {
			t.Error(err)
		}
		if !match {
			t.Errorf("expected pattern (second line): \"%s\", got: %s", functionSignatureRegex, line)
		}
	}
}

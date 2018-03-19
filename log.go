// Copyright 2018, Irfan Sharif.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package log

import (
	"bytes"
	"fmt"
	"io"
	"runtime"
	"strings"
)

// Rename "mode" as you can have multiple. Need something like mode, span.
type mode int

const (
	InfoMode mode = 1 << iota
	WarnMode
	ErrorMode
	FatalMode
	DebugMode

	// TODO(irfansharif): Is there a way to not check for this explicitly every
	// time? Somehow embed this in the representation?
	DisabledMode = 0
	DefaultMode  = InfoMode | WarnMode | ErrorMode | FatalMode
)

func (m mode) String() string {
	var buffer bytes.Buffer
	switch m {
	case InfoMode:
		buffer.WriteString("I")
		fallthrough
	case WarnMode:
		buffer.WriteString("W")
		fallthrough
	case ErrorMode:
		buffer.WriteString("E")
		fallthrough
	case FatalMode:
		buffer.WriteString("F")
		fallthrough
	case DebugMode:
		buffer.WriteString("D")
		fallthrough
	case DisabledMode:
		buffer.WriteString("X")
		fallthrough
	default:
		panic("Unidentified mode")
	}
	return buffer.String()
}

type Logger struct {
	w io.Writer
}

// TODO(irfansharif): Wrap configurations in variadic options.
// TODO(irfansharif): Make a rotated log, dir io.Writer.
func New(w io.Writer) *Logger {
	l := &Logger{
		w: w,
	}
	return l
}

// getcaller returns the file and line number of where the caller's caller's
// call site.
//
// f.go: 11 func f() {
// f.go: 12      g()
// f.go: 13 }
//
//
// g.go: 27 func g() {
// g.go: 28      file, line := getcaller()
// g.go: 29      fmt.Println(fmt.Sprintf("%s: %d", file, line)) // f.go: 12
// g.go: 29 }
//
func getcaller() (file string, line int) {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "[???]"
		line = -1
	} else {
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
	}
	return file, line
}

func (l *Logger) Info(v ...interface{}) {
	file, line := getcaller()
	pc := fmt.Sprintf("%s:%d", file, line)

	gmode := getLogMode()
	if (gmode & InfoMode) != 0 {
		for i := 0; i < len(v); i++ {
			switch t := v[i].(type) {
			case string:
				l.w.Write([]byte(fmt.Sprintf("%s: %s", pc, t)))
			}
		}
		return
	}

	lmode, ok := getPCMode(pc)
	if !ok {
		return
	}
	if (lmode & InfoMode) != 0 {
		for i := 0; i < len(v); i++ {
			switch t := v[i].(type) {
			case string:
				l.w.Write([]byte(fmt.Sprintf("%s: %s", pc, t)))
			}
		}
		return
	}
}

func (l *Logger) Infof(format string, v ...interface{}) {
	panic("unimplemented")
}

func (l *Logger) Warn(v ...interface{}) {
	panic("unimplemented")
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	panic("unimplemented")
}

func (l *Logger) Error(v ...interface{}) {
	panic("unimplemented")
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	panic("unimplemented")
}

func (l *Logger) Fatal(v ...interface{}) {
	panic("unimplemented")
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	panic("unimplemented")
}

func (l *Logger) Debug(v ...interface{}) {
	panic("unimplemented")
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	panic("unimplemented")
}

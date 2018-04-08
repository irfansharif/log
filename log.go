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
	"fmt"
	"io"
	"runtime"
	"runtime/debug"
	"strings"
)

// Logger. TODO(irfansharif): Comment.
type Logger struct {
	w io.Writer
}

// TODO(irfansharif): Wrap configurations in variadic options.
// TODO(irfansharif): Make a rotated log, dir io.Writer.
func New(w io.Writer) *Logger {
	return &Logger{
		w: w,
	}
}

func (l *Logger) Info(v ...interface{}) {
	l.log(InfoMode, fmt.Sprintf("%v", v...))
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.log(InfoMode, fmt.Sprintf(format, v...))
}

func (l *Logger) Warn(v ...interface{}) {
	l.log(WarnMode, fmt.Sprintf("%v", v...))
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.log(WarnMode, fmt.Sprintf(format, v...))
}

func (l *Logger) Error(v ...interface{}) {
	l.log(ErrorMode, fmt.Sprintf("%v", v...))
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.log(ErrorMode, fmt.Sprintf(format, v...))
}

func (l *Logger) Fatal(v ...interface{}) {
	l.log(FatalMode, fmt.Sprintf("%v", v...))
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.log(FatalMode, fmt.Sprintf(format, v...))
}

func (l *Logger) Debug(v ...interface{}) {
	l.log(DebugMode, fmt.Sprintf("%v", v...))
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.log(DebugMode, fmt.Sprintf(format, v...))
}

// TODO(irfansharif): Document that backtrace may be emitted despite not being
// in the right mode.
func (l *Logger) log(lmode mode, data string) {
	// Logger.log is called from Logger.{Info,Warn,Error,Fatal,Debug}{,f}. We
	// use a depth of two to retrieve the caller immediately preceding it.
	file, line := caller(2)
	tp := fmt.Sprintf("%s:%d", file, line)

	gmode := GetGlobalLogMode()
	fmode, ok := GetFileLogMode(file)
	shouldLog := (gmode&lmode) != DisabledMode || (ok && (fmode&lmode) != DisabledMode)
	if shouldLog {
		l.w.Write([]byte(fmt.Sprintf("%s %s: %s\n", lmode.String(), tp, data)))
	}

	tpenabled := CheckTracePoint(tp)
	if tpenabled {
		l.w.Write(stacktrace())
	}
}

// stacktrace returns the stack trace for the current goroutine.
func stacktrace() []byte {
	return debug.Stack()
}

// caller returns the file and line number of where the caller's caller's
// call site.
//
// e.go: 32 func e() {
// e.go: 33     f()
// e.go: 34 }
//
// f.go: 11 func f() {
// f.go: 12      g()
// f.go: 13 }
//
// g.go: 25 func g() {
// g.go: 26 	{
// g.go: 27         // Request caller one level above.
// g.go: 28 		file, line := caller(1)
// g.go: 29 		fmt.Println(fmt.Sprintf("%s: %d", file, line)) // f.go: 12
// g.go: 30 	}
// g.go: 31 	{
// g.go: 32         // Request caller two levels above.
// g.go: 33 		file, line := caller(2)
// g.go: 34 		fmt.Println(fmt.Sprintf("%s: %d", file, line)) // e.go: 33
// g.go: 35 	}
// g.go: 36 }
//
func caller(depth int) (file string, line int) {
	_, file, line, ok := runtime.Caller(depth + 1) // +1 to account for call to caller itself.
	if !ok {
		file = "[???]"
		line = -1
	} else {
		// TODO(irfansharif): Right now this isn't robust to shared filenames
		// across varied sub packages. This is a stand-in to allow for direct
		// file name specification without fully-specified paths (in the host
		// machine or relative to project root).
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
	}
	return file, line
}

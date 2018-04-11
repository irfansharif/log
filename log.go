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
	"runtime/debug"
	"strings"
)

// TODO(irfansharif): Provide a catchall global logger with warning?
// TODO(irfansharif): Allow for logging flags, configuring output.
// TODO(irfansharif): Allow for tags.
// TODO(irfansharif): Allow for custom leveling (verbosity, should work at the
// same filtered granularities).

// Logger. TODO(irfansharif): Comment.
type Logger struct {
	w io.Writer
}

func configure(l *Logger) {
}

// TODO(irfansharif): Wrap configurations in variadic options.
func New(w io.Writer, options ...option) *Logger {
	l := &Logger{
		w: w,
	}
	for _, option := range options {
		option(l)
	}
	configure(l)
	return l
}

func (l *Logger) Info(v ...interface{}) {
	l.log(InfoMode, fmt.Sprintln(v...))
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.log(InfoMode, fmt.Sprintf(format, v...))
}

func (l *Logger) Warn(v ...interface{}) {
	l.log(WarnMode, fmt.Sprint(v...))
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.log(WarnMode, fmt.Sprintf(format, v...))
}

func (l *Logger) Error(v ...interface{}) {
	l.log(ErrorMode, fmt.Sprint(v...))
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.log(ErrorMode, fmt.Sprintf(format, v...))
}

func (l *Logger) Fatal(v ...interface{}) {
	l.log(FatalMode, fmt.Sprint(v...))
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.log(FatalMode, fmt.Sprintf(format, v...))
}

func (l *Logger) Debug(v ...interface{}) {
	l.log(DebugMode, fmt.Sprint(v...))
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.log(DebugMode, fmt.Sprintf(format, v...))
}

// TODO(irfansharif): Document that backtrace may be emitted at statements from any mode.
func (l *Logger) log(lmode Mode, data string) {
	// Logger.log is only to be called from
	// Logger.{Info,Warn,Error,Fatal,Debug}{,f}. We use a depth of two to
	// retrieve the caller immediately preceding it.
	file, line := caller(2)
	tp := fmt.Sprintf("%s:%d", file, line)

	tpenabled := GetTracePoint(tp)
	if tpenabled {
		// Skip logger.log, and the invoking public wrapper
		// Logger.{Info,Warn,Error,Fatal,Debug}{,f}
		// TODO(irfansharif): Should this have the logger header as well?
		l.w.Write(stacktrace(2))
	}

	gmode := GetGlobalLogMode()
	fmode, ok := GetFileLogMode(file)
	shouldLog := (gmode&lmode) != DisabledMode || (ok && (fmode&lmode) != DisabledMode)
	if !shouldLog {
		return
	}

	l.w.Write([]byte(fmt.Sprintf("%s %s: %s", lmode.string(), tp, data)))
}

// TODO(irfansharif): Better comment. stacktrace returns the stack trace for
// the current goroutine, skipping n immediately preceding function traces
// (last being the caller, inclusive of the caller).
func stacktrace(skip int) []byte {
	skip *= 2 // Each function depth corresponds to to lines of stack trace output.
	skip += 2 // For debug.Stack()
	skip += 2 // For this function, log.stacktrace()

	b := debug.Stack()
	bs := bytes.Split(b, []byte("\n"))
	copy(bs[1:], bs[1+skip:]) // TODO(irfansharif): Is this a pointer copy? It could be.
	bs = bs[:len(bs)-skip]
	return bytes.Join(bs, []byte("\n"))
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

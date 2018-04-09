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
	"os"
	"os/user"
	"path/filepath"
	"time"
)

var (
	program  = "?"
	hostname = "?"
	username = "?"
	pid      = -1
)

func init() {
	program = filepath.Base(os.Args[0])

	host, err := os.Hostname()
	if err == nil {
		hostname = host
	}

	u, err := user.Current()
	if err == nil {
		username = u.Username
	}

	pid = os.Getpid()
}

// TODO(irfansharif): Provide default synchronized os.Stderr writer.

func filename(t time.Time) (fname, lname string) {
	fname = fmt.Sprintf("%s.%s.%s.%s.%d.log",
		program,
		hostname,
		username,
		t.Format("2006-01-02.15:04:05.999"),
		pid,
	)
	lname = fmt.Sprintf("%s.log", program)
	return fname, lname
}

// TODO(irfansharif): Document.
func RotatingDirectoryLogger(dirname string) io.Writer {
	os.MkdirAll(dirname, os.ModePerm)

	fname, lname := filename(time.Now())
	fname1, lname1 := filepath.Join(dirname, fname), filepath.Join(dirname, lname)
	_, _ = fname1, lname1

	f, err := os.Create(filepath.Join(dirname, fname))
	if err != nil {
		// TODO(irfansharif): Propagate linker error; can't afford to fail silently here.
		fmt.Println(err)
	}

	os.Remove(filepath.Join(dirname, lname))         // Remove symlink, if any, ignore error.
	os.Symlink(fname, filepath.Join(dirname, lname)) // Best effort symlinking, ignore error.

	return f
}

// TODO(irfansharif): Actually synchronize.
func SynchronizedWriter(w io.Writer) io.Writer {
	return w
}

type multiWriter struct {
	ws []io.Writer
}

func (m *multiWriter) Write(b []byte) (n int, err error) {
	for _, w := range m.ws {
		w.Write(b)
	}
	return len(b), nil
}

// TODO(irfansharif): Document.
func MultiWriter(w io.Writer, ws ...io.Writer) io.Writer {
	mw := &multiWriter{
		ws: nil,
	}
	mw.ws = append(mw.ws, w)
	mw.ws = append(mw.ws, ws...)
	return mw
}

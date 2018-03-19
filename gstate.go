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
	"sync"
	"sync/atomic"
)

type pcModeMap map[string]mode // Map from program counter fname.go:linenumber to mode.
type gstateT struct {
	cmode    atomic.Value
	pcModeMu struct {
		sync.Mutex
		m atomic.Value
	}
}

var gstate gstateT

// Need to initialize the atomics; to be used once during init time and for
// tests.
func resetgstate() {
	gstate.cmode.Store(DefaultMode)
	gstate.pcModeMu.m.Store(make(pcModeMap))
}

func init() {
	resetgstate()
}

func SetLogMode(m mode) {
	gstate.cmode.Store(m)
}

func getLogMode() mode {
	return gstate.cmode.Load().(mode)
}

func setPCMode(pc string, m mode) {
	gstate.pcModeMu.Lock()                     // Synchronize with other potential writers.
	ma := gstate.pcModeMu.m.Load().(pcModeMap) // Load current value of the map.
	mb := make(pcModeMap)                      // Create a new map.
	for pc, m := range ma {
		mb[pc] = m // Copy all data from the current object to the new one.
	}
	mb[pc] = m                  // Do the update that we need.
	gstate.pcModeMu.m.Store(mb) // Atomically replace the current object with the new one.
	// At this point all new readers start working with the new version.
	// The old version will be garbage collected once the existing readers
	// (if any) are done with it.
	gstate.pcModeMu.Unlock()
}

func getPCMode(pc string) (mode, bool) {
	pcmap := gstate.pcModeMu.m.Load().(pcModeMap)
	m, ok := pcmap[pc]
	return m, ok
}

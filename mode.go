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

type mode int

const (
	// TODO(irfansharif): FatalMode cannot be disabled. Ensure, perhaps there's
	// a way to enforce this implicitly in representation.
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
	switch m {
	case InfoMode:
		return "I"
	case WarnMode:
		return "W"
	case ErrorMode:
		return "E"
	case FatalMode:
		return "F"
	case DebugMode:
		return "D"
	case DisabledMode:
		return "X"
	default:
		return "?"
	}
}
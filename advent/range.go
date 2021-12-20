// Copyright 2018 Kyle Lemons
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package advent

// RangeTracker is a helper for tracking the min/max of integers.
//
// Its zero value is usable, but Min/Max should not be considered if !Valid.
type RangeTracker struct {
	Valid    bool // true if Track has been called
	Min, Max int
}

// Track updates the tracker based on the int, and returns it for easy chaining.
func (rt *RangeTracker) Track(n int) int {
	if !rt.Valid {
		rt.Min = n
		rt.Max = n
		rt.Valid = true
		return n
	}
	if n > rt.Max {
		rt.Max = n
	}
	if n < rt.Min {
		rt.Min = n
	}
	return n
}

// TrackAll is like Track but tracks all of the given numbers, not returning anything.
func (rt *RangeTracker) TrackAll(n ...int) {
	for _, n := range n {
		rt.Track(n)
	}
}

// Delta returns the size of the range.
func (rt *RangeTracker) Delta() int {
	if !rt.Valid {
		panic("Delta called on invalid RangeTracker")
	}
	return rt.Max - rt.Min
}

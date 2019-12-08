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

// Package acoday is the entrypoint for this AoC solution.
package acoday

import (
	"testing"

	"github.com/kylelemons/adventofcodesolutions/advent"
)

func part1(t *testing.T, in string) (ret int) {
	var items []int
	advent.Words(in).Scan(t, func(v int) {
		items = append(items, v)
	})

	pop := func() (v int) {
		items, v = items[1:], items[0]
		return
	}
	var process func()
	process = func() {
		childCount, metaCount := pop(), pop()
		for i := 0; i < childCount; i++ {
			process()
		}
		for i := 0; i < metaCount; i++ {
			ret += pop()
		}
	}
	process()
	return
}

func part2(t *testing.T, in string) (ret int) {
	var items []int
	advent.Words(in).Scan(t, func(v int) {
		items = append(items, v)
	})

	pop := func() (v int) {
		items, v = items[1:], items[0]
		return
	}
	var process func() int
	process = func() (ret int) {
		childCount, metaCount := pop(), pop()
		switch childCount {
		case 0:
			for i := 0; i < metaCount; i++ {
				ret += pop()
			}
		default:
			var sub []int
			for i := 0; i < childCount; i++ {
				sub = append(sub, process())
			}
			for i := 0; i < metaCount; i++ {
				idx := pop()
				if idx <= 0 || idx > len(sub) {
					continue
				}
				ret += sub[idx-1]
			}
		}
		return
	}
	return process()
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part1 example 0", "2 3 0 3 10 11 12 1 1 0 1 99 2 1 1 2", 138},
		{"part1 answer", advent.ReadFile(t, "input.txt"), 44338},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part1(t, test.in), test.want; got != want {
				t.Errorf("part1(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

func TestPart2(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part2 example 0", "2 3 0 3 10 11 12 1 1 0 1 99 2 1 1 2", 66},
		{"part2 answer", advent.ReadFile(t, "input.txt"), 37560},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(t, test.in), test.want; got != want {
				t.Errorf("part2(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

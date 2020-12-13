// Copyright 2019 Kyle Lemons
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
	"strconv"
	"testing"

	"github.com/kylelemons/adventofcodesolutions/advent"
)

func part1(t *testing.T, in string) (ret int) {
	var loI, hiI int
	advent.Scanner(in).Extract(t, `(\d+)-(\d+)`, &loI, &hiI)

nextNumber:
	for i := loI; i <= hiI; i++ {
		var double bool
		s := strconv.Itoa(i)
		for j := 0; j+1 < len(s); j++ {
			double = double || s[j] == s[j+1]
			if s[j] > s[j+1] {
				continue nextNumber
			}
		}
		if double {
			ret++
			// t.Log(i)
		}
	}
	return
}

func part2(t *testing.T, in string) (ret int) {
	var loI, hiI int
	advent.Scanner(in).Extract(t, `(\d+)-(\d+)`, &loI, &hiI)

nextNumber:
	for i := loI; i <= hiI; i++ {
		s := strconv.Itoa(i)
		for j := 0; j+1 < len(s); j++ {
			if s[j] > s[j+1] {
				continue nextNumber
			}
		}
		var counts [10]int
		for _, c := range s {
			counts[int(c-'0')]++
		}
		for _, c := range counts {
			if c == 2 {
				ret++
				continue nextNumber
			}
		}
	}
	return
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		// {"part1 example 0", "...", 0},
		{"part1 answer", "387638-919123", 466},
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
		// {"part2 example 0", "...", 0},
		{"part2 answer", "387638-919123", 292},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(t, test.in), test.want; got != want {
				t.Errorf("part2(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

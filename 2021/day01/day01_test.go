// Copyright 2021 Kyle Lemons
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
package aocday

import (
	"math"
	"testing"

	"github.com/kylelemons/adventofcodesolutions/advent"
)

type Input struct {
	depths []int
}

func parseInput(t *testing.T, in string) *Input {
	input := &Input{
		// ...
	}
	advent.Lines(in).Extract(t, `(.*)`, func(depth int) {
		input.depths = append(input.depths, depth)
	})
	return input
}

func part1(t *testing.T, in string) (ret int) {
	input := parseInput(t, in)

	prev := math.MaxInt
	for _, d := range input.depths {
		if d > prev {
			ret++
		}
		prev = d
	}

	return
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part1 example 0", `199
200
208
210
200
207
240
269
260
263`, 7},
		{"part1 answer", advent.ReadFile(t, "input.txt"), 1665},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part1(t, test.in), test.want; got != want {
				t.Errorf("part1(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

func part2(t *testing.T, in string) (ret int) {
	input := parseInput(t, in)

	for i := range input.depths {
		if i <= 2 {
			continue
		}
		if input.depths[i] > input.depths[i-3] {
			ret++
		}
	}

	return
}

func TestPart2(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part2 example 0", `199
200
208
210
200
207
240
269
260
263`, 5},
		{"part2 answer", advent.ReadFile(t, "input.txt"), 1702},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(t, test.in), test.want; got != want {
				t.Errorf("part2(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

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
	"sort"
	"testing"

	"github.com/kylelemons/adventofcodesolutions/advent"
)

type Input struct {
	Positions []int
}

func parseInput(t *testing.T, in string) *Input {
	input := &Input{
		// ...
	}
	advent.Split(in, ',').Scan(t, func(v int) {
		input.Positions = append(input.Positions, v)
	})
	sort.Ints(input.Positions)
	return input
}

func part1(t *testing.T, in string) (ret int) {
	input := parseInput(t, in)

	ret = math.MaxInt
	for loc := input.Positions[0]; loc <= input.Positions[len(input.Positions)-1]; loc++ {
		var fuel int
		for _, pos := range input.Positions {
			delta := pos - loc
			if delta < 0 {
				delta = -delta
			}
			fuel += delta
		}
		if fuel < ret {
			ret = fuel
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
		{"part1 example 0", "16,1,2,0,4,2,7,1,2,14", 37},
		{"part1 answer", advent.ReadFile(t, "input.txt"), 349812},
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

	var fuelFor func(delta int) int
	fuelFor = func(delta int) int {
		if delta <= 0 {
			return 0
		}
		return delta + fuelFor(delta-1)
	}
	cache := make(map[int]int)
	orig := fuelFor
	fuelFor = func(delta int) int {
		if ans, ok := cache[delta]; ok {
			return ans
		}
		ans := orig(delta)
		cache[delta] = ans
		return ans
	}

	ret = math.MaxInt
	for loc := input.Positions[0]; loc <= input.Positions[len(input.Positions)-1]; loc++ {
		var fuel int
		for _, pos := range input.Positions {
			delta := pos - loc
			if delta < 0 {
				delta = -delta
			}
			fuel += fuelFor(delta)
		}
		if fuel < ret {
			ret = fuel
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
		{"part2 example 0", "16,1,2,0,4,2,7,1,2,14", 168},
		{"part2 answer", advent.ReadFile(t, "input.txt"), 99763899},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(t, test.in), test.want; got != want {
				t.Errorf("part2(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

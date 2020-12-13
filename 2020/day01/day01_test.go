// Copyright 2020 Kyle Lemons
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
	"testing"

	"github.com/kylelemons/adventofcodesolutions/advent"
)

type Input struct {
	Expenses []int
}

func parseInput(t *testing.T, in string) *Input {
	input := &Input{
		// ...
	}
	advent.Lines(in).Extract(t, `(.*)`, func(exp int) {
		input.Expenses = append(input.Expenses, exp)
	})
	return input
}

func part1(t *testing.T, in string) (ret int) {
	input := parseInput(t, in)

	seen := make(map[int]bool)
	for _, exp := range input.Expenses {
		want := 2020 - exp
		if seen[want] {
			return want * exp
		}
		seen[exp] = true
	}

	return -1
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part1 example 0", `1721
979
366
299
675
1456`, 514579},
		{"part1 answer", advent.ReadFile(t, "input.txt"), 355875},
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

	for i1, v1 := range input.Expenses {
		for i2, v2 := range input.Expenses[i1:] {
			for _, v3 := range input.Expenses[i1+i2:] {
				if v1+v2+v3 == 2020 {
					return v1 * v2 * v3
				}
			}
		}
	}

	return -1
}

func TestPart2(t *testing.T) {

	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part2 example 0", `1721
979
366
299
675
1456`, 241861950},
		{"part2 answer", advent.ReadFile(t, "input.txt"), 140379120},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(t, test.in), test.want; got != want {
				t.Errorf("part2(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

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
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/kylelemons/adventofcodesolutions/advent"
)

type Input struct {
	seats []int
}

func parseInput(t *testing.T, in string) *Input {
	input := &Input{
		// ...
	}
	conv := strings.NewReplacer(
		"F", "0",
		"B", "1",
		"L", "0",
		"R", "1",
	)
	advent.Lines(in).Extract(t, `(.*)`, func(v string) {
		binrep := conv.Replace(v)
		numrep, err := strconv.ParseInt(binrep, 2, 64)
		if err != nil {
			t.Fatalf("%q is invalid binary int", binrep)
		}
		input.seats = append(input.seats, int(numrep))
	})
	return input
}

func part1(t *testing.T, in string) (ret int) {
	input := parseInput(t, in)

	sort.Ints(input.seats)
	return input.seats[len(input.seats)-1]
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part1 example 0", `BFFFBBFRRR`, 567},
		{"part1 example 1", `FFFBBBFRRR`, 119},
		{"part1 example 2", `BBFFBBFRLL`, 820},
		{"part1 answer", advent.ReadFile(t, "input.txt"), 930},
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

	sort.Ints(input.seats)
	for i := range input.seats[:len(input.seats)-1] {
		if input.seats[i]+1 != input.seats[i+1] {
			return input.seats[i] + 1
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
		{"part2 answer", advent.ReadFile(t, "input.txt"), 515},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(t, test.in), test.want; got != want {
				t.Errorf("part2(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

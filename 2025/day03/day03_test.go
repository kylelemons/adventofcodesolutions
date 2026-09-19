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
	"bufio"
	"cmp"
	"slices"
	"strings"
	"testing"

	"github.com/kylelemons/adventofcodesolutions/advent"
)

type Input struct {
	Banks []string
}

func parseInput(t *testing.T, in string) *Input {
	input := &Input{
		// ...
	}
	lines := bufio.NewScanner(strings.NewReader(in))
	for lines.Scan() {
		input.Banks = append(input.Banks, strings.TrimSpace(lines.Text()))
	}
	return input
}

func maxJoltage1(bank string) int {
	type pair struct {
		digit int
		pos   int
	}
	var pairs []pair
	for i, c := range bank {
		pairs = append(pairs, pair{
			digit: int(c - '0'),
			pos:   i,
		})
	}
	slices.SortFunc(pairs, func(a, b pair) int {
		return cmp.Or(
			cmp.Compare(b.digit, a.digit), // highest first
			cmp.Compare(a.pos, b.pos),     // lowest first
		)
	})

	for _, bestFirst := range pairs {
		bestSecond := pair{
			digit: -1,         // sentinel
			pos:   len(pairs), // worst case (out of range)
		}
		for _, p := range pairs {
			if p.pos <= bestFirst.pos {
				// can't pick this digit
				continue
			}
			if p.digit > bestSecond.digit {
				bestSecond = p
			}
		}
		if bestSecond.digit != -1 {
			return bestFirst.digit*10 + bestSecond.digit
		}
	}
	panic("not found")
}

func part1(t *testing.T, in string) (ret int) {
	input := parseInput(t, in)

	for _, bank := range input.Banks {
		ret += maxJoltage1(bank)
	}

	return
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part1 example 0", "987654321111111\n811111111111119\n234234234234278\n818181911112111", 357},
		{"part1 answer", advent.ReadFile(t, "input.txt"), 17766},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part1(t, test.in), test.want; got != want {
				t.Errorf("part1(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

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
	"fmt"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/kylelemons/adventofcodesolutions/advent"
)

type Input struct {
	FreshRanges [][2]int
	Ingredients []int
}

func parseInput(t *testing.T, in string) *Input {
	input := &Input{
		// ...
	}
	lines := bufio.NewScanner(strings.NewReader(in))
	for lines.Scan() {
		line := lines.Text()
		if len(line) == 0 {
			break
		}

		lo, hi, ok := strings.Cut(line, "-")
		if !ok {
			panic(fmt.Sprintf("invalid range %q", line))
		}
		loInt, err := strconv.Atoi(lo)
		if err != nil {
			panic(fmt.Sprintf("invalid range %q", line))
		}
		hiInt, err := strconv.Atoi(hi)
		if err != nil {
			panic(fmt.Sprintf("invalid range %q", line))
		}
		input.FreshRanges = append(input.FreshRanges, [2]int{loInt, hiInt})
	}
	for lines.Scan() {
		line := lines.Text()
		if len(line) == 0 {
			break
		}
		ingredient, err := strconv.Atoi(line)
		if err != nil {
			panic(fmt.Sprintf("invalid ingredient %q", line))
		}
		input.Ingredients = append(input.Ingredients, ingredient)
	}

	slices.SortFunc(input.FreshRanges, func(a, b [2]int) int {
		return cmp.Or(
			cmp.Compare(a[0], b[0]),
			cmp.Compare(a[1], b[1]),
		)
	})

	return input
}

func part1(t *testing.T, in string) (ret int) {
	input := parseInput(t, in)

	for _, ingr := range input.Ingredients {
		for _, rng := range input.FreshRanges {
			if ingr >= rng[0] && ingr <= rng[1] {
				ret++
				break
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
		{"part1 example 0", "3-5\n10-14\n16-20\n12-18\n\n1\n5\n8\n11\n17\n32", 3},
		{"part1 answer", advent.ReadFile(t, "input.txt"), 517},
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

	lastHi := 0
	for _, rng := range input.FreshRanges {
		lo, hi := rng[0], rng[1]
		if lo <= lastHi {
			lo = lastHi + 1
		}
		if lo > hi {
			continue
		}
		ret += hi - lo + 1 // inclusive
		lastHi = hi
	}

	return
}

func TestPart2(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part2 example 0", "3-5\n10-14\n16-20\n12-18\n\n1\n5\n8\n11\n17\n32", 14},
		{"part2 answer", advent.ReadFile(t, "input.txt"), 336173027056994},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(t, test.in), test.want; got != want {
				t.Errorf("part2(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

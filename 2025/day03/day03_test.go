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
	"math"
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

func maxJoltagePart2(bank string) int {
	// Recurrence relation:
	//   maxJoltage( [prefix] + [suffix], n ) = maxJoltage([prefix], 1) concat maxJoltage( [suffix], n-1 )

	concat := func(a, b, bDigits int) int {
		return a*int(math.Pow10(bDigits)) + b
	}

	var maxJoltage func(bank string, digits int) int
	maxJoltage = func(bank string, digits int) (_ret int) {
		// defer func() {
		// 	fmt.Printf("maxJoltage(%q, %d) = %d\n", bank, digits, _ret)
		// }()

		if len(bank) < digits {
			return -1
		}

		if digits == 0 {
			panic("can't find joltage of empty string")
		}
		if digits == 1 {
			maxChar := int(bank[0] - '0')
			for _, char := range bank[1:] {
				value := int(char - '0')
				if value > maxChar {
					maxChar = value
				}
			}
			return maxChar
		}

		maxValue := -1
		for prefixLen := range len(bank) {
			prefixJoltage := maxJoltage(bank[:prefixLen], 1)
			suffixJoltage := maxJoltage(bank[prefixLen:], digits-1)
			if prefixJoltage > 0 && suffixJoltage > 0 {
				value := concat(prefixJoltage, suffixJoltage, digits-1)
				if value > maxValue {
					maxValue = value
					// fmt.Printf("New max: joltage(%q):joltage(%q) = %v\n",
					// 	bank[:prefixLen], bank[prefixLen:],
					// 	value,
					// )
				}
			}
		}
		return maxValue
	}

	type cacheKey struct {
		bank   string
		digits int
	}
	cache := make(map[cacheKey]int)
	orig := maxJoltage
	maxJoltage = func(bank string, digits int) int {
		key := cacheKey{bank, digits}
		if v, ok := cache[key]; ok {
			return v
		}
		v := orig(bank, digits)
		cache[key] = v
		return v
	}

	return maxJoltage(bank, 12)
}

func part2(t *testing.T, in string) (ret int) {
	input := parseInput(t, in)

	for _, bank := range input.Banks {
		ret += maxJoltagePart2(bank)
	}

	return
}

func TestPart2(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part2 example 0", "987654321111111\n811111111111119\n234234234234278\n818181911112111", 3121910778619},
		{"part2 answer", advent.ReadFile(t, "input.txt"), 176582889354075},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(t, test.in), test.want; got != want {
				t.Errorf("part2(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

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
	"strings"
	"testing"

	"github.com/kylelemons/adventofcodesolutions/advent"
)

const CapsBit = 'A' ^ 'a'

func part1(t *testing.T, polymer string) (ret int) {
	for i := 1; i < len(polymer); i++ {
		j := i - 1
		// t.Logf("%q | ...%q | %q ~ %q", polymer, polymer[j:], polymer[i], polymer[j])
		for j >= 0 && i < len(polymer) && polymer[i]^CapsBit == polymer[j] {
			i++
			j--
		}
		if i-j > 1 {
			polymer = polymer[:j+1] + polymer[i:]
			i = j
		}
	}
	// t.Logf("%q = %d", polymer, len(polymer))
	return len(polymer)
}

func part2(t *testing.T, polymer string) (ret int) {
	ret = part1(t, strings.NewReplacer("a", "", "A", "").Replace(polymer))
	for low, hi := byte('b'), byte('B'); low <= 'z'; low, hi = low+1, hi+1 {
		v := part1(t, strings.NewReplacer(string([]byte{low}), "", string([]byte{hi}), "").Replace(polymer))
		if v < ret {
			ret = v
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
		{"part1 example 0", "dabAcCaCBAcCcaDA", 10},
		{"part1 answer", advent.ReadFile(t, "input.txt"), 10564},
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
		{"part2 example 0", "dabAcCaCBAcCcaDA", 4},
		{"part2 answer", advent.ReadFile(t, "input.txt"), 6336},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(t, test.in), test.want; got != want {
				t.Errorf("part2(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

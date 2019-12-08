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

func part1(t *testing.T, inState, inTrans string, generations int) (ret int) {
	var initial string
	advent.Scanner(inState).Extract(t, `initial state: ([.#]+)`, &initial)

	translate := make(map[string]byte)
	advent.Lines(inTrans).Extract(t, `([.#]{5}) => ([.#])`, func(from, to string) {
		translate[from] = to[0]
	})

	current, left := initial, 0
	for i := 0; i < generations; i++ {
		current = "...." + current + "...." // avoid out of bounds
		left -= 4

		next := new(strings.Builder)
		for j := range current {
			if j < 2 || j >= len(current)-2 {
				next.WriteByte('.')
				continue
			}
			b := byte('.')
			if n, ok := translate[current[j-2:][:5]]; ok {
				b = n
			}
			next.WriteByte(b)
		}
		current = strings.TrimRight(next.String(), ".")

		for current[0] == '.' {
			current = current[1:]
			left++
		}
	}

	for i, c := range current {
		if c == '#' {
			ret += i + left
		}
	}

	return
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name  string
		init  string
		trans string
		gens  int
		want  int
	}{
		{"part1 example 0", "initial state: #..#.#..##......###...###", strings.TrimSpace(`
...## => #
..#.. => #
.#... => #
.#.#. => #
.#.## => #
.##.. => #
.#### => #
#.#.# => #
#.### => #
##.#. => #
##.## => #
###.. => #
###.# => #
####. => #
`), 20, 325},
		{"part1 answer", advent.ReadFile(t, "state.txt"), advent.ReadFile(t, "transitions.txt"), 20, 3120},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part1(t, test.init, test.trans, test.gens), test.want; got != want {
				t.Errorf("part1(%#v, %#v)\n = %#v, want %#v", test.init, test.trans, got, want)
			}
		})
	}
}

func part2(t *testing.T, inState, inTrans string, generations int) (ret int) {
	var initial string
	advent.Scanner(inState).Extract(t, `initial state: ([.#]+)`, &initial)

	translate := make(map[string]byte)
	advent.Lines(inTrans).Extract(t, `([.#]{5}) => ([.#])`, func(from, to string) {
		translate[from] = to[0]
	})

	type state struct{ generation, left int }
	seen := make(map[string]state)

	current, left := initial, 0
	for i := 0; i < generations; i++ {
		if s, ok := seen[current]; ok {
			cycle := i - s.generation
			shift := left - s.left
			t.Logf("Cycle detected: gen=%d-%d, length=%d, shift=%d", s.generation, i, cycle, shift)
			t.Logf("  %s", current)

			remaining := generations - i - 1
			cycles := remaining / cycle
			i += cycles * cycle
			left += cycles * shift
			t.Logf("Fast forward: %d cycles to gen=%d, left=%d", cycles, i, left)
		} else {
			seen[current] = state{i, left}
		}

		current = "...." + current + "...." // avoid out of bounds
		left -= 4

		temp := new(strings.Builder)
		for j := range current {
			if j < 2 || j >= len(current)-2 {
				temp.WriteByte('.')
				continue
			}
			b := byte('.')
			if n, ok := translate[current[j-2:][:5]]; ok {
				b = n
			}
			temp.WriteByte(b)
		}

		next := strings.TrimRight(temp.String(), ".")
		for next[0] == '.' {
			next = next[1:]
			left++
		}

		current = next
	}

	for i, c := range current {
		if c == '#' {
			ret += i + left
		}
	}

	return
}

func TestPart2(t *testing.T) {
	tests := []struct {
		name  string
		init  string
		trans string
		gens  int
		want  int
	}{
		{"part2 example 0", "initial state: #..#.#..##......###...###", strings.TrimSpace(`
...## => #
..#.. => #
.#... => #
.#.#. => #
.#.## => #
.##.. => #
.#### => #
#.#.# => #
#.### => #
##.#. => #
##.## => #
###.. => #
###.# => #
####. => #
`), 20, 325},
		{"part2 answer", advent.ReadFile(t, "state.txt"), advent.ReadFile(t, "transitions.txt"), 50000000000, 2950000001598},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(t, test.init, test.trans, test.gens), test.want; got != want {
				t.Errorf("part2(%#v, %#v)\n = %#v, want %#v", test.init, test.trans, got, want)
			}
		})
	}
}

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
package aocday

import (
	"sort"
	"testing"

	"github.com/kylelemons/adventofcodesolutions/advent"
)

type Input struct {
	numbers []int
}

func parseInput(t *testing.T, in string) *Input {
	input := &Input{
		// ...
	}
	advent.Lines(in).Scan(t, func(n int) {
		input.numbers = append(input.numbers, n)
	})
	return input
}

func part1(t *testing.T, in string, winsz int) (ret int) {
	input := parseInput(t, in)

nextNumber:
	for i := winsz; i < len(input.numbers); i++ {
		n := input.numbers[i]
		win := input.numbers[i-winsz:][:winsz]

		needed := make(map[int]bool)
		for _, v0 := range win {
			if needed[v0] {
				continue nextNumber
			}

			v1 := n - v0
			needed[v1] = true
		}
		return n // no pair found
	}

	return -1
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name string
		n    int
		in   string
		want int
	}{
		{"part1 example 0", 5, `35
20
15
25
47
40
62
55
65
95
102
117
150
182
127
219
299
277
309
576`, 127},
		{"part1 answer", 25, advent.ReadFile(t, "input.txt"), 26134589},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part1(t, test.in, test.n), test.want; got != want {
				t.Errorf("part1(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

func part2(t *testing.T, in string, target int) (ret int) {
	input := parseInput(t, in)

	sum := 0
	win := []int{}
	for _, v := range input.numbers {
		sum += v
		win = append(win, v)
		switch {
		case sum < target:
			// need more numbers, keep going
		case sum > target:
			// too high, remove part of the window
			for len(win) > 0 && sum > target {
				remove := win[0]
				sum -= remove
				win = win[1:]
			}
			if sum != target {
				break
			}
			fallthrough
		default:
			// sum == target, answer found; sort to get min/max
			sort.Ints(win)
			return win[0] + win[len(win)-1]
		}

	}

	return -1
}

func TestPart2(t *testing.T) {
	tests := []struct {
		name string
		n    int
		in   string
		want int
	}{
		{"part2 example 0", 5, `35
20
15
25
47
40
62
55
65
95
102
117
150
182
127
219
299
277
309
576`, 62},
		{"part2 answer", 25, advent.ReadFile(t, "input.txt"), 3535124},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sum := part1(t, test.in, test.n)
			if got, want := part2(t, test.in, sum), test.want; got != want {
				t.Errorf("part2(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

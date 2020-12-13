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
	"testing"

	"github.com/kylelemons/adventofcodesolutions/advent"
)

type Input struct {
	groups []group
}

type group struct {
	answers []map[string]bool
}

func parseInput(t *testing.T, in string) *Input {
	input := &Input{
		// ...
	}
	advent.Records(in).Each(func(i int, lines advent.Scanner) {
		g := group{}
		advent.Lines(string(lines)).Scan(t, func(line string) {
			answer := make(map[string]bool)
			for i := range line {
				answer[line[i:][:1]] = true
			}
			g.answers = append(g.answers, answer)
		})
		input.groups = append(input.groups, g)
	})
	return input
}

func part1(t *testing.T, in string) (ret int) {
	input := parseInput(t, in)

	for _, g := range input.groups {
		answers := make(map[string]bool)
		for _, ans := range g.answers {
			for k := range ans {
				answers[k] = true
			}
		}
		ret += len(answers)
	}
	return
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part1 example 0", `abc

a
b
c

ab
ac

a
a
a
a

b`, 11},
		{"part1 answer", advent.ReadFile(t, "input.txt"), 6612},
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

	for _, g := range input.groups {
		answers := make(map[string]int)
		for _, ans := range g.answers {
			for k := range ans {
				answers[k]++
			}
		}
		for _, v := range answers {
			if v == len(g.answers) {
				ret++
			}
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
		{"part2 example 0", `abc

a
b
c

ab
ac

a
a
a
a

b`, 6},
		{"part2 answer", advent.ReadFile(t, "input.txt"), 3268},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(t, test.in), test.want; got != want {
				t.Errorf("part2(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

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
	"testing"

	"github.com/kylelemons/adventofcodesolutions/advent"
)

type Command struct {
	What  string
	Count int
}

type Input struct {
	Commands []Command
}

func parseInput(t *testing.T, in string) *Input {
	input := &Input{
		// ...
	}
	advent.Lines(in).Scan(t, func(what string, count int) {
		input.Commands = append(input.Commands, Command{what, count})
	})
	return input
}

func part1(t *testing.T, in string) (ret int) {
	input := parseInput(t, in)

	var pos, depth int
	for _, c := range input.Commands {
		switch c.What {
		case "forward":
			pos += c.Count
		case "down":
			depth += c.Count
		case "up":
			depth -= c.Count
		}
	}

	return pos * depth
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part1 example 0", `forward 5
down 5
forward 8
up 3
down 8
forward 2`, 150},
		{"part1 answer", advent.ReadFile(t, "input.txt"), 2322630},
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

	var pos, depth, aim int
	for _, c := range input.Commands {
		switch c.What {
		case "forward":
			pos += c.Count
			depth += aim * c.Count
		case "down":
			aim += c.Count
		case "up":
			aim -= c.Count
		}
	}

	return pos * depth
}

func TestPart2(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part2 example 0", `forward 5
down 5
forward 8
up 3
down 8
forward 2`, 900},
		{"part2 answer", advent.ReadFile(t, "input.txt"), 2105273490},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(t, test.in), test.want; got != want {
				t.Errorf("part2(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

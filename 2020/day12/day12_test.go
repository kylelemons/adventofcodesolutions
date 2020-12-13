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
	"github.com/kylelemons/adventofcodesolutions/advent/coords"
)

type Instruction struct {
	Op    string
	Count int
}

type Input struct {
	Instrs []Instruction
}

func parseInput(t *testing.T, in string) *Input {
	input := &Input{
		// ...
	}
	advent.Lines(in).Extract(t, `([NESWLRF])(\d+)`, func(c string, n int) {
		input.Instrs = append(input.Instrs, Instruction{
			Op:    c,
			Count: n,
		})
	})
	return input
}

func part1(t *testing.T, in string) (ret int) {
	input := parseInput(t, in)

	facing := coords.East
	loc := coords.RC(0, 0)
	for _, instr := range input.Instrs {
		switch instr.Op {
		case "N":
			loc = loc.Add(coords.North.Scale(instr.Count))
		case "E":
			loc = loc.Add(coords.East.Scale(instr.Count))
		case "S":
			loc = loc.Add(coords.South.Scale(instr.Count))
		case "W":
			loc = loc.Add(coords.West.Scale(instr.Count))
		case "L":
			for i := 0; i < instr.Count; i += 90 {
				facing = facing.Left()
			}
		case "R":
			for i := 0; i < instr.Count; i += 90 {
				facing = facing.Right()
			}
		case "F":
			loc = loc.Add(facing.Scale(instr.Count))
		}
	}

	r := loc.R()
	if r < 0 {
		r = -r
	}
	c := loc.C()
	if c < 0 {
		c = -c
	}
	return r + c
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part1 example 0", `F10
N3
F7
R90
F11`, 25},
		{"part1 answer", advent.ReadFile(t, "input.txt"), 1152},
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

	loc := coords.RC(0, 0)
	way := coords.East.Scale(10).Add(coords.North.Scale(1))
	for _, instr := range input.Instrs {
		switch instr.Op {
		case "N":
			way = way.Add(coords.North.Scale(instr.Count))
		case "E":
			way = way.Add(coords.East.Scale(instr.Count))
		case "S":
			way = way.Add(coords.South.Scale(instr.Count))
		case "W":
			way = way.Add(coords.West.Scale(instr.Count))
		case "L":
			for i := 0; i < instr.Count; i += 90 {
				way = way.Left()
			}
		case "R":
			for i := 0; i < instr.Count; i += 90 {
				way = way.Right()
			}
		case "F":
			loc = loc.Add(way.Scale(instr.Count))
		}
	}

	r := loc.R()
	if r < 0 {
		r = -r
	}
	c := loc.C()
	if c < 0 {
		c = -c
	}
	return r + c
}

func TestPart2(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part2 example 0", `F10
N3
F7
R90
F11`, 286},
		{"part2 answer", advent.ReadFile(t, "input.txt"), 58637},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(t, test.in), test.want; got != want {
				t.Errorf("part2(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

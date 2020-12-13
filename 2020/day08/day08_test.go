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

type Program struct {
	PC    int
	Instr []Instruction

	Accumulator int
}

type Instruction struct {
	Op   string
	Arg0 int
}

func (p *Program) Step(t *testing.T) {
	if p.PC < 0 || p.PC > len(p.Instr) {
		t.Fatalf("PC (%d) out of bounds (%d)", p.PC, len(p.Instr))
	}
	instr := p.Instr[p.PC]
	switch instr.Op {
	case "nop":
	case "acc":
		p.Accumulator += instr.Arg0
	case "jmp":
		p.PC += instr.Arg0 - 1 // offset the decr that will happen later
	}
	p.PC++
}

func (p *Program) Halted() bool {
	return p.PC == len(p.Instr)
}

func parseInput(t *testing.T, in string) *Program {
	input := &Program{
		// ...
	}
	advent.Lines(in).Extract(t, `(\w+) ([+-]?\d+)`, func(op string, arg0 int) {
		input.Instr = append(input.Instr, Instruction{
			Op:   op,
			Arg0: arg0,
		})
	})
	return input
}

func part1(t *testing.T, in string) (ret int) {
	input := parseInput(t, in)

	seen := make(map[int]bool)
	for {
		if seen[input.PC] {
			return input.Accumulator
		}
		seen[input.PC] = true
		input.Step(t)
	}
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part1 example 0", `nop +0
acc +1
jmp +4
acc +3
jmp -3
acc -99
acc +1
jmp -4
acc +6`, 5},
		{"part1 answer", advent.ReadFile(t, "input.txt"), 1832},
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
	try := func(input *Program) (acc int, halt bool) {
		seen := make(map[int]bool)
		for {
			if input.Halted() {
				return input.Accumulator, true
			}
			if seen[input.PC] {
				return input.Accumulator, false
			}
			seen[input.PC] = true
			input.Step(t)
		}
	}
	input := parseInput(t, in)
	for i := range input.Instr {
		prog := parseInput(t, in)
		switch prog.Instr[i].Op {
		case "nop":
			prog.Instr[i].Op = "jmp"
		case "jmp":
			prog.Instr[i].Op = "nop"
		default:
			continue
		}
		acc, halt := try(prog)
		if halt {
			return acc
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
		{"part2 example 0", `nop +0
acc +1
jmp +4
acc +3
jmp -3
acc -99
acc +1
jmp -4
acc +6`, 8},
		{"part2 answer", advent.ReadFile(t, "input.txt"), 662},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(t, test.in), test.want; got != want {
				t.Errorf("part2(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

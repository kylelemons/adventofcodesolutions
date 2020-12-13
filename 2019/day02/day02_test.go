// Copyright 2019 Kyle Lemons
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

func part1(t *testing.T, in string) (pos0 int) {
	var mem []int
	for _, s := range strings.Split(in, ",") {
		var i int
		advent.Scanner(s).Scan(t, &i)
		mem = append(mem, i)
	}

	mem[1] = 12
	mem[2] = 02

	pc := 0
	for {
		if pc < 0 || pc > len(mem) {
			t.Fatalf("PC %d out of bounds", pc)
		}
		op := mem[pc]
		switch op {
		case 1: // add
			src1, src2, dst := mem[pc+1], mem[pc+2], mem[pc+3]
			// t.Logf("ADD: prog[%d] = prog[%d] (%v) + prog[%d] (%v)", dst, src1, mem[src1], src2, mem[src2])
			mem[dst] = mem[src1] + mem[src2]
			pc += 4
		case 2: // mul
			src1, src2, dst := mem[pc+1], mem[pc+2], mem[pc+3]
			// t.Logf("MUL: prog[%d] = prog[%d] (%v) * prog[%d] (%v)", dst, src1, mem[src1], src2, mem[src2])
			mem[dst] = mem[src1] * mem[src2]
			pc += 4
		case 99:
			// t.Logf("END")
			return mem[0]
		default:
			t.Fatalf("UNKNOWN OPCODE %d", op)
		}
	}
}

func part2(t *testing.T, in string) (output int) {
	run := func(noun, verb int) (prog0 int) {
		var mem []int
		for _, s := range strings.Split(in, ",") {
			var i int
			advent.Scanner(s).Scan(t, &i)
			mem = append(mem, i)
		}

		mem[1] = noun
		mem[2] = verb

		pc := 0
		for {
			if pc < 0 || pc > len(mem) {
				t.Fatalf("PC %d out of bounds", pc)
			}
			op := mem[pc]
			switch op {
			case 1: // add
				src1, src2, dst := mem[pc+1], mem[pc+2], mem[pc+3]
				mem[dst] = mem[src1] + mem[src2]
				pc += 4
			case 2: // mul
				src1, src2, dst := mem[pc+1], mem[pc+2], mem[pc+3]
				mem[dst] = mem[src1] * mem[src2]
				pc += 4
			case 99:
				return mem[0]
			default:
				t.Fatalf("UNKNOWN OPCODE %d", op)
			}
		}
	}

	for noun := 0; noun < 100; noun++ {
		for verb := 0; verb < 100; verb++ {
			if run(noun, verb) == 19690720 {
				return 100*noun + verb
			}
		}
	}
	return -1
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part1 answer", advent.ReadFile(t, "input.txt"), 3706713},
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
		{"part2 answer", advent.ReadFile(t, "input.txt"), 8609},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(t, test.in), test.want; got != want {
				t.Errorf("part2(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

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
	"fmt"
	"strconv"
	"testing"

	"github.com/kylelemons/adventofcodesolutions/advent"
)

// A Program is an executable instruction with the 2019 instruction set.
type Program struct {
	Source string     // "source code" aka comma-separated instructions
	Input  func() int // Inputs as rquested by the program
	Output func(int)  // Outputs from the program
}

// Run a program and return the memory after it halts.
func (p *Program) Run(t *testing.T) (mem []int) {
	advent.Split(p.Source, ',').Scan(t, func(instr int) {
		mem = append(mem, instr)
	})

	pc := 0
	adv := func() (v int) {
		v, pc = mem[pc], pc+1
		return
	}

	for {
		if pc < 0 || pc > len(mem) {
			t.Fatalf("PC %d out of bounds", pc)
		}

		next := adv()
		op, flags := next%100, strconv.Itoa(next/100)

		type Param struct {
			mode  byte
			value int
		}
		get := func() Param {
			p := Param{'0', adv()}
			if n := len(flags); n > 0 {
				flags, p.mode = flags[:n-1], flags[n-1]
			}
			return p
		}

		mode := func(p Param) *int {
			switch p.mode {
			case '0': // positional
				return &mem[p.value]
			case '1': // immediate
				v := p.value // paranoia
				return &v
			default:
				t.Fatalf("Unrecognized flag %q", p.mode)
				return nil
			}
		}
		debug := func(p Param) string {
			switch p.mode {
			case '0':
				return fmt.Sprintf("(%v @ mem[%d])", mem[p.value], p.value)
			case '1':
				return fmt.Sprintf("(%v)", p.value)
			default:
				t.Fatalf("Unrecognized flag %q", p.mode)
				return ""
			}
		}
		switch op {
		case 1: // add
			src1, src2, dst := get(), get(), get()
			t.Logf("ADD %s + %s -> %s", debug(src1), debug(src2), debug(dst))
			*mode(dst) = *mode(src1) + *mode(src2)
		case 2: // mul
			src1, src2, dst := get(), get(), get()
			t.Logf("MUL %s * %s -> %s", debug(src1), debug(src2), debug(dst))
			*mode(dst) = *mode(src1) * *mode(src2)
		case 3: // input
			dst := get()
			t.Logf("INPUT -> %s", debug(dst))
			*mode(dst) = p.Input()
		case 4: // output
			src := get()
			t.Logf("OUTPUT <- %s", debug(src))
			p.Output(*mode(src))
		case 5: // jump-nonzero
			cond, to := get(), get()
			t.Logf("JNZ %s to %s", debug(cond), debug(to))
			if *mode(cond) != 0 {
				t.Logf("  ... branch taken")
				pc = *mode(to)
			}
		case 6: // jump-zero
			cond, to := get(), get()
			t.Logf("JZ %s to %s", debug(cond), debug(to))
			if *mode(cond) == 0 {
				t.Logf("  ... branch taken")
				pc = *mode(to)
			}
		case 7: // less-than
			a, b, dst := get(), get(), get()
			t.Logf("LT %s < %s -> %s", debug(a), debug(b), debug(dst))
			if *mode(a) < *mode(b) {
				*mode(dst) = 1
			} else {
				*mode(dst) = 0
			}
		case 8: // equals
			a, b, dst := get(), get(), get()
			t.Logf("EQ %s == %s -> %s", debug(a), debug(b), debug(dst))
			if *mode(a) == *mode(b) {
				*mode(dst) = 1
			} else {
				*mode(dst) = 0
			}
		case 99:
			t.Logf("HALT")
			return mem
		default:
			t.Fatalf("UNKNOWN OPCODE %d", op)
		}
	}
}

func part1(t *testing.T, in string) (ret int) {
	(&Program{
		Source: in,
		Input:  func() int { return 1 },
		Output: func(v int) {
			// Previous output wasn't terminal code, so check it:
			if ret != 0 {
				t.Errorf("CHECK failed (off by %d)", ret)
			}
			t.Logf("CHECK: %v", v)
			ret = v
		},
	}).Run(t)
	return
}

func part2(t *testing.T, in string, input int) (ret int) {
	(&Program{
		Source: in,
		Input:  func() int { return input },
		Output: func(v int) {
			// Previous output wasn't terminal code, so check it:
			if ret != 0 {
				t.Errorf("CHECK failed (off by %d)", ret)
			}
			t.Logf("CHECK: %v", v)
			ret = v
		},
	}).Run(t)
	return
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part1 example 0", "1002,4,3,4,33", 0},
		{"part1 example 1", "1101,100,-1,4,0", 0},
		{"part1 answer", advent.ReadFile(t, "input.txt"), 3122865},
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
		name  string
		in    string
		input int
		want  int
	}{
		{"part2 example 0 ne", "3,9,8,9,10,9,4,9,99,-1,8", 5, 0},
		{"part2 example 0 eq", "3,9,8,9,10,9,4,9,99,-1,8", 8, 1},
		{"part2 example 1 lt", "3,9,7,9,10,9,4,9,99,-1,8", 5, 1},
		{"part2 example 1 gt", "3,9,7,9,10,9,4,9,99,-1,8", 9, 0},
		{"part2 example 2 ine", "3,3,1108,-1,8,3,4,3,99", 5, 0},
		{"part2 example 2 ieq", "3,3,1108,-1,8,3,4,3,99", 8, 1},
		{"part2 example 3 ilt", "3,3,1107,-1,8,3,4,3,99", 5, 1},
		{"part2 example 3 igt", "3,3,1107,-1,8,3,4,3,99", 9, 0},
		{"part2 example 4 jpos", "3,12,6,12,15,1,13,14,13,4,13,99,-1,0,1,9", 0, 0},
		{"part2 example 4 njpos", "3,12,6,12,15,1,13,14,13,4,13,99,-1,0,1,9", 5, 1},
		{"part2 example 5 jim", "3,3,1105,-1,9,1101,0,0,12,4,12,99,1", 0, 0},
		{"part2 example 5 njim", "3,3,1105,-1,9,1101,0,0,12,4,12,99,1", 5, 1},
		{"part2 big example a", "3,21,1008,21,8,20,1005,20,22,107,8,21,20,1006,20,31,1106,0,36,98,0,0,1002,21,125,20,4,20,1105,1,46,104,999,1105,1,46,1101,1000,1,20,4,20,1105,1,46,98,99", 7, 999},
		{"part2 big example b", "3,21,1008,21,8,20,1005,20,22,107,8,21,20,1006,20,31,1106,0,36,98,0,0,1002,21,125,20,4,20,1105,1,46,104,999,1105,1,46,1101,1000,1,20,4,20,1105,1,46,98,99", 8, 1000},
		{"part2 big example c", "3,21,1008,21,8,20,1005,20,22,107,8,21,20,1006,20,31,1106,0,36,98,0,0,1002,21,125,20,4,20,1105,1,46,104,999,1105,1,46,1101,1000,1,20,4,20,1105,1,46,98,99", 9, 1001},
		{"part2 answer", advent.ReadFile(t, "input.txt"), 5, 773660},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(t, test.in, test.input), test.want; got != want {
				t.Errorf("part2(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

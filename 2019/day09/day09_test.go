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
	rel := 0
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

		brk := func(size int) int {
			for len(mem) < size {
				mem = append(mem, make([]int, 1024)...)
			}
			return size
		}
		mode := func(p Param) *int {
			switch p.mode {
			case '0': // positional
				return &mem[brk(p.value)]
			case '1': // immediate
				v := p.value // paranoia
				return &v
			case '2': // relative
				return &mem[brk(p.value+rel)]
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
		_ = debug
		switch op {
		case 1: // add
			src1, src2, dst := get(), get(), get()
			// t.Logf("ADD %s + %s -> %s", debug(src1), debug(src2), debug(dst))
			*mode(dst) = *mode(src1) + *mode(src2)
		case 2: // mul
			src1, src2, dst := get(), get(), get()
			// t.Logf("MUL %s * %s -> %s", debug(src1), debug(src2), debug(dst))
			*mode(dst) = *mode(src1) * *mode(src2)
		case 3: // input
			dst := get()
			// t.Logf("INPUT -> %s", debug(dst))
			*mode(dst) = p.Input()
		case 4: // output
			src := get()
			// t.Logf("OUTPUT <- %s", debug(src))
			p.Output(*mode(src))
		case 5: // jump-nonzero
			cond, to := get(), get()
			// t.Logf("JNZ %s to %s", debug(cond), debug(to))
			if *mode(cond) != 0 {
				// t.Logf("  ... branch taken")
				pc = *mode(to)
			}
		case 6: // jump-zero
			cond, to := get(), get()
			// t.Logf("JZ %s to %s", debug(cond), debug(to))
			if *mode(cond) == 0 {
				// t.Logf("  ... branch taken")
				pc = *mode(to)
			}
		case 7: // less-than
			a, b, dst := get(), get(), get()
			// t.Logf("LT %s < %s -> %s", debug(a), debug(b), debug(dst))
			if *mode(a) < *mode(b) {
				*mode(dst) = 1
			} else {
				*mode(dst) = 0
			}
		case 8: // equals
			a, b, dst := get(), get(), get()
			// t.Logf("EQ %s == %s -> %s", debug(a), debug(b), debug(dst))
			if *mode(a) == *mode(b) {
				*mode(dst) = 1
			} else {
				*mode(dst) = 0
			}
		case 9: // adjrel
			a := get()
			rel += *mode(a)
		case 99:
			// t.Logf("HALT")
			return mem
		default:
			t.Fatalf("UNKNOWN OPCODE %d", op)
		}
	}
}
func part1(t *testing.T, in string) (ret int) {
	prog := &Program{
		Source: in,
		Input:  func() int { return 1 },
		Output: func(v int) { ret = v },
	}
	prog.Run(t)

	return
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		// {"part1 example 0", "...", 0},
		{"part1 answer", advent.ReadFile(t, "input.txt"), 3742852857},
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
	prog := &Program{
		Source: in,
		Input:  func() int { return 2 },
		Output: func(v int) { ret = v },
	}
	prog.Run(t)

	return
}

func TestPart2(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part2 answer", advent.ReadFile(t, "input.txt"), 73439},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(t, test.in), test.want; got != want {
				t.Errorf("part2(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

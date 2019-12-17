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

// Package intcode implements the 2019 Advent of Code intcode computer.
package intcode

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/kylelemons/adventofcodesolutions/advent"
)

// A Program is an executable instruction with the 2019 instruction set.
type Program struct {
	Input  func() int // Inputs as rquested by the program
	Output func(int)  // Outputs from the program

	// Memory contains the (mutable) instruction memory for the program.
	Memory []int

	// Debugf is called with debugging information during execution.
	Debugf func(format string, args ...interface{})
}

// Compile compiles source into a program.
func Compile(t *testing.T, source string) *Program {
	var mem []int
	advent.Split(source, ',').Scan(t, func(instr int) {
		mem = append(mem, instr)
	})
	return &Program{
		Memory: mem,
		Input:  func() int { t.Fatalf("Input requested but no Input function provided to program"); return 0 },
		Output: func(v int) { t.Logf("Output(%v)", v) },
		Debugf: func(string, ...interface{}) {},
	}
}

// Run a program and return the memory after it halts.
func (p *Program) Run(t *testing.T) {
	mem, pc, rel := p.Memory, 0, 0

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
			for len(mem) <= size {
				mem = append(mem, make([]int, 1024)...)
			}
			return size
		}
		mode := func(param Param) *int {
			switch param.mode {
			case '0': // positional
				return &mem[brk(param.value)]
			case '1': // immediate
				v := param.value // paranoia
				return &v
			case '2': // relative
				return &mem[brk(param.value+rel)]
			default:
				t.Fatalf("Unrecognized flag %q", param.mode)
				return nil
			}
		}
		debug := func(param Param) fmt.Stringer {
			switch param.mode {
			case '0':
				return lazyf("(%v @ mem[%d])", mem[brk(param.value)], param.value)
			case '1':
				return lazyf("(%v)", param.value)
			case '2': // relative
				return lazyf("(%v @ mem[%d+%d])", mem[brk(param.value+rel)], param.value, rel)
			default:
				fmt.Printf("Unrecognized flag %q\n", param.mode)
				t.Fatalf("Unrecognized flag %q", param.mode)
				return nil
			}
		}
		_ = debug
		switch op {
		case 1: // add
			src1, src2, dst := get(), get(), get()
			p.Debugf("ADD %s + %s -> %s", debug(src1), debug(src2), debug(dst))
			*mode(dst) = *mode(src1) + *mode(src2)
		case 2: // mul
			src1, src2, dst := get(), get(), get()
			p.Debugf("MUL %s * %s -> %s", debug(src1), debug(src2), debug(dst))
			*mode(dst) = *mode(src1) * *mode(src2)
		case 3: // input
			dst := get()
			p.Debugf("INPUT -> %s", debug(dst))
			*mode(dst) = p.Input()
		case 4: // output
			src := get()
			p.Debugf("OUTPUT <- %s", debug(src))
			p.Output(*mode(src))
		case 5: // jump-nonzero
			cond, to := get(), get()
			p.Debugf("JNZ %s to %s", debug(cond), debug(to))
			if *mode(cond) != 0 {
				p.Debugf("  ... branch taken")
				pc = *mode(to)
			}
		case 6: // jump-zero
			cond, to := get(), get()
			p.Debugf("JZ %s to %s", debug(cond), debug(to))
			if *mode(cond) == 0 {
				p.Debugf("  ... branch taken")
				pc = *mode(to)
			}
		case 7: // less-than
			a, b, dst := get(), get(), get()
			p.Debugf("LT %s < %s -> %s", debug(a), debug(b), debug(dst))
			if *mode(a) < *mode(b) {
				*mode(dst) = 1
			} else {
				*mode(dst) = 0
			}
		case 8: // equals
			a, b, dst := get(), get(), get()
			p.Debugf("EQ %s == %s -> %s", debug(a), debug(b), debug(dst))
			if *mode(a) == *mode(b) {
				*mode(dst) = 1
			} else {
				*mode(dst) = 0
			}
		case 9: // adjrel
			a := get()
			rel += *mode(a)
		case 99:
			p.Debugf("HALT")
			return
		default:
			t.Fatalf("UNKNOWN OPCODE %d", op)
		}
	}
}

type lazyPrinter struct {
	format string
	args   []interface{}
}

func (p lazyPrinter) String() string { return fmt.Sprintf(p.format, p.args...) }

func lazyf(format string, args ...interface{}) lazyPrinter {
	return lazyPrinter{format, args}
}

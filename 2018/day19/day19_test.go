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
	"bytes"
	"flag"
	"fmt"
	"log"
	"math"
	"strings"
	"testing"
	"time"

	"github.com/kylelemons/adventofcodesolutions/advent"
)

var debug = flag.Bool("debug", false, "If true, print debugging information")

type Computer struct {
	regs []int
}

func (c *Computer) addr(dst, a, b int, regs []int) { regs[dst] = regs[a] + regs[b] }
func (c *Computer) addi(dst, a, b int, regs []int) { regs[dst] = regs[a] + b }
func (c *Computer) mulr(dst, a, b int, regs []int) { regs[dst] = regs[a] * regs[b] }
func (c *Computer) muli(dst, a, b int, regs []int) { regs[dst] = regs[a] * b }
func (c *Computer) banr(dst, a, b int, regs []int) { regs[dst] = regs[a] & regs[b] }
func (c *Computer) bani(dst, a, b int, regs []int) { regs[dst] = regs[a] & b }
func (c *Computer) borr(dst, a, b int, regs []int) { regs[dst] = regs[a] | regs[b] }
func (c *Computer) bori(dst, a, b int, regs []int) { regs[dst] = regs[a] | b }
func (c *Computer) setr(dst, a, b int, regs []int) { regs[dst] = regs[a] }
func (c *Computer) seti(dst, a, b int, regs []int) { regs[dst] = a }
func (c *Computer) gtir(dst, a, b int, regs []int) { regs[dst] = ibool(a > regs[b]) }
func (c *Computer) gtri(dst, a, b int, regs []int) { regs[dst] = ibool(regs[a] > b) }
func (c *Computer) gtrr(dst, a, b int, regs []int) { regs[dst] = ibool(regs[a] > regs[b]) }
func (c *Computer) eqir(dst, a, b int, regs []int) { regs[dst] = ibool(a == regs[b]) }
func (c *Computer) eqri(dst, a, b int, regs []int) { regs[dst] = ibool(regs[a] == b) }
func (c *Computer) eqrr(dst, a, b int, regs []int) { regs[dst] = ibool(regs[a] == regs[b]) }

var instr2fun = map[string]func(c *Computer, dst, a, b int, regs []int){
	"addr": (*Computer).addr,
	"addi": (*Computer).addi,
	"mulr": (*Computer).mulr,
	"muli": (*Computer).muli,
	"banr": (*Computer).banr,
	"bani": (*Computer).bani,
	"borr": (*Computer).borr,
	"bori": (*Computer).bori,
	"setr": (*Computer).setr,
	"seti": (*Computer).seti,
	"gtir": (*Computer).gtir,
	"gtri": (*Computer).gtri,
	"gtrr": (*Computer).gtrr,
	"eqir": (*Computer).eqir,
	"eqri": (*Computer).eqri,
	"eqrr": (*Computer).eqrr,
}

func (c *Computer) Call(t *testing.T, i Instr) {
	instr2fun[i.Opcode](c, i.Dst, i.A, i.B, c.regs)
}

func ibool(b bool) int {
	if b {
		return 1
	}
	return 0
}

type Instr struct {
	Opcode    string
	A, B, Dst int
}

type Input struct {
	Comp Computer

	Registers int
	PCReg     int
	Instrs    []Instr
}

func parseInput(t *testing.T, in string) *Input {
	input := &Input{
		Registers: 6,
		PCReg:     -1,
	}
	advent.Lines(in).Each(func(i int, line advent.Scanner) {
		var instr Instr
		switch {
		case line.CanExtract(t, `#ip (\d+)`, &input.PCReg):
		case line.CanExtract(t, `(\S+) (\d+) (\d+) (\d+)(?:.*#.*)?`, advent.Fields(&instr)...):
			input.Instrs = append(input.Instrs, instr)
		default:
			t.Fatalf("Unrecognized line: %q", line)
		}
	})
	if input.PCReg < 0 {
		t.Fatalf("No PC Register set in program")
	}
	input.Comp.regs = make([]int, input.Registers)
	return input
}

func (in *Input) Execute(t *testing.T) {
	status := time.NewTicker(1 * time.Second)
	defer status.Stop()

	touched := bytes.Repeat([]byte{'-'}, len(in.Instrs))
	before := make([]int, len(in.Comp.regs))
	var ops int
	for {
		pc := in.Comp.regs[in.PCReg]
		if pc >= len(in.Instrs) || pc < 0 {
			return
		}
		copy(before, in.Comp.regs)
		instr := in.Instrs[pc]
		touched[pc] = '#'
		in.Comp.Call(t, instr)
		in.Comp.regs[in.PCReg]++

		ops++
		select {
		case <-status.C:
			log.Printf("pc=%-3d | instr=%-15s | before=%10v | after=%10v | touched=%s | ops=%v",
				pc, fmt.Sprint(instr), before, in.Comp.regs, touched, ops)
			ops = 0
			touched = bytes.Repeat([]byte{'-'}, len(in.Instrs))
		default:
		}
	}
}

func part1(t *testing.T, in string) (ret int) {
	input := parseInput(t, in)
	input.Execute(t)
	return input.Comp.regs[0]
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part1 example 0", strings.Trim(`
#ip 0
seti 5 0 1
seti 6 0 2
addi 0 1 0
addr 1 2 3
setr 1 0 0
seti 8 0 4
seti 9 0 5
`, "\n"), 7},
		{"part1 answer", advent.ReadFile(t, "input.txt"), 2047},
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
	// input := parseInput(t, in)
	// input.Comp.regs[0] = 1
	// input.Execute(t)
	// return input.Comp.regs[0]

	target := 10551424
	max := int(math.Sqrt(float64(target)) + 1)
	for factor := 1; factor <= max; factor++ {
		quo, rem := target/factor, target%factor
		if rem == 0 {
			ret += quo + factor
		}
	}
	return
}

/* Notes about part 2:

Approximate translation of the program:

Registers:
	r0 : output
	r1 : factor
	r2 : target
	r3 : pc
	r4 : temp
	r5 : counter

Program:
	target = 10551424
	for factor = 1; factor <= target; factor++
		for counter = 1; counter <= target; counter++
			if factor*counter == target:
				output += factor

Answer:
	sum of (factors of 10551424) = 24033240
*/

func TestPart2(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part2 answer", advent.ReadFile(t, "input.txt"), 24033240},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(t, test.in), test.want; got != want {
				t.Errorf("part2(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

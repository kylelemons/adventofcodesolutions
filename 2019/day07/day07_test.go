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
	"sync"
	"testing"

	"github.com/kylelemons/adventofcodesolutions/2019/advent"
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
		case 99:
			// t.Logf("HALT")
			return mem
		default:
			t.Fatalf("UNKNOWN OPCODE %d", op)
		}
	}
}

func part1(t *testing.T, in string) (ret int) {
	for i := 0; i < 5; i++ {
		avail := make([]int, 0, 4)
		for c := 0; c < 5; c++ {
			switch c {
			case i:
			default:
				avail = append(avail, c)
			}
		}
		for _, j := range avail {
			avail := make([]int, 0, 3)
			for c := 0; c < 5; c++ {
				switch c {
				case i, j:
				default:
					avail = append(avail, c)
				}
			}
			for _, k := range avail {
				avail := make([]int, 0, 2)
				for c := 0; c < 5; c++ {
					switch c {
					case i, j, k:
					default:
						avail = append(avail, c)
					}
				}
				for _, l := range avail {
					phases := []int{i, j, k, l, 0}
					for c := 0; c < 5; c++ {
						switch c {
						case i, j, k, l:
						default:
							phases[4] = c
						}
					}

					var ampInput int
					for amp := 0; amp < 5; amp++ {
						inputs := []int{phases[amp], ampInput}
						(&Program{
							Source: in,
							Input:  func() (v int) { v, inputs = inputs[0], inputs[1:]; return },
							Output: func(v int) {
								ampInput = v
							},
						}).Run(t)
					}
					if ampInput > ret {
						ret = ampInput
						// t.Logf("New best input: %v = %v", phases, ret)
					}

				}
			}
		}
	}
	return
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part1 example 0", "3,15,3,16,1002,16,10,16,1,16,15,15,4,15,99,0,0", 43210},
		{"part1 example 0", "3,23,3,24,1002,24,10,24,1002,23,-1,23,101,5,23,23,1,24,23,23,4,23,99,0,0", 54321},
		{"part1 answer", advent.ReadFile(t, "input.txt"), 24405},
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
	for i := 5; i < 10; i++ {
		avail := make([]int, 0, 4)
		for c := 5; c < 10; c++ {
			switch c {
			case i:
			default:
				avail = append(avail, c)
			}
		}
		for _, j := range avail {
			avail := make([]int, 0, 3)
			for c := 5; c < 10; c++ {
				switch c {
				case i, j:
				default:
					avail = append(avail, c)
				}
			}
			for _, k := range avail {
				avail := make([]int, 0, 2)
				for c := 5; c < 10; c++ {
					switch c {
					case i, j, k:
					default:
						avail = append(avail, c)
					}
				}
				for _, l := range avail {
					phases := []int{i, j, k, l, 0}
					for c := 5; c < 10; c++ {
						switch c {
						case i, j, k, l:
						default:
							phases[4] = c
						}
					}

					var chans []chan int
					for c := 0; c < 5; c++ {
						chans = append(chans, make(chan int, 1))
					}

					var wg sync.WaitGroup
					for amp := 0; amp < 5; amp++ {
						amp := amp

						wg.Add(1)
						go func() {
							defer wg.Done()
							var sentPhase bool
							(&Program{
								Source: in,
								Input: func() (v int) {
									if !sentPhase {
										sentPhase = true
										return phases[amp]
									}
									return <-chans[amp]
								},
								Output: func(v int) {
									i := (amp + 1) % len(chans)
									chans[i] <- v
								},
							}).Run(t)
						}()
					}
					chans[0] <- 0
					wg.Wait()
					if got := <-chans[0]; got > ret {
						ret = got
						// t.Logf("New best input: %v = %v", phases, ret)
					}

				}
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
		{"part2 example 0", "3,26,1001,26,-4,26,3,27,1002,27,2,27,1,27,26,27,4,27,1001,28,-1,28,1005,28,6,99,0,0,5", 139629729},
		{"part2 example 0", "3,52,1001,52,-5,52,3,53,1,52,56,54,1007,54,5,55,1005,55,26,1001,54,-5,54,1105,1,12,1,53,54,53,1008,54,0,55,1001,55,1,55,2,53,55,53,4,53,1001,56,-1,56,1005,56,6,99,0,0,0,0,10", 18216},
		{"part2 answer", advent.ReadFile(t, "input.txt"), 8271623},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(t, test.in), test.want; got != want {
				t.Errorf("part2(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

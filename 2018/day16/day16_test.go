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
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kylelemons/adventofcodesolutions/advent"
)

func part1(t *testing.T, in string) (ret int) {
	for _, obs := range strings.Split(in, "\n\n") {
		lines := strings.Split(obs, "\n")

		var r0, r1, r2, r3 int
		advent.Scanner(lines[0]).Extract(t, `Before: \[(-?\d+), (-?\d+), (-?\d+), (-?\d+)\]`, &r0, &r1, &r2, &r3)
		before := []int{r0, r1, r2, r3}

		var instr, a, b, dst int
		advent.Scanner(lines[1]).Scan(t, &instr, &a, &b, &dst)

		advent.Scanner(lines[2]).Extract(t, `After:  \[(-?\d+), (-?\d+), (-?\d+), (-?\d+)\]`, &r0, &r1, &r2, &r3)
		after := []int{r0, r1, r2, r3}

		var possible []string
		try := func(instr string, f func([]int)) {
			regs := make([]int, 4)
			copy(regs, before)
			f(regs)
			if cmp.Equal(regs, after) {
				possible = append(possible, instr)
			}
		}

		try("addr", func(regs []int) {
			regs[dst] = regs[a] + regs[b]
		})
		try("addi", func(regs []int) {
			regs[dst] = regs[a] + b
		})

		try("mulr", func(regs []int) {
			regs[dst] = regs[a] * regs[b]
		})
		try("muli", func(regs []int) {
			regs[dst] = regs[a] * b
		})

		try("banr", func(regs []int) {
			regs[dst] = regs[a] & regs[b]
		})
		try("bani", func(regs []int) {
			regs[dst] = regs[a] & b
		})

		try("borr", func(regs []int) {
			regs[dst] = regs[a] | regs[b]
		})
		try("bori", func(regs []int) {
			regs[dst] = regs[a] | b
		})

		try("setr", func(regs []int) {
			regs[dst] = regs[a]
		})
		try("seti", func(regs []int) {
			regs[dst] = a
		})

		try("gtir", func(regs []int) {
			if a > regs[b] {
				regs[dst] = 1
			} else {
				regs[dst] = 0
			}
		})
		try("gtri", func(regs []int) {
			if regs[a] > b {
				regs[dst] = 1
			} else {
				regs[dst] = 0
			}
		})
		try("gtrr", func(regs []int) {
			if regs[a] > regs[b] {
				regs[dst] = 1
			} else {
				regs[dst] = 0
			}
		})

		try("eqir", func(regs []int) {
			if a == regs[b] {
				regs[dst] = 1
			} else {
				regs[dst] = 0
			}
		})
		try("eqri", func(regs []int) {
			if regs[a] == b {
				regs[dst] = 1
			} else {
				regs[dst] = 0
			}
		})
		try("eqrr", func(regs []int) {
			if regs[a] == regs[b] {
				regs[dst] = 1
			} else {
				regs[dst] = 0
			}
		})

		// t.Logf("Instruction %d could be %v", instr, possible)
		if len(possible) >= 3 {
			ret++
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
		{"part1 example 0", "Before: [3, 2, 1, 1]\n9 2 1 2\nAfter:  [3, 2, 2, 1]", 1},
		{"part1 answer", advent.ReadFile(t, "input1.txt"), 509},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part1(t, test.in), test.want; got != want {
				t.Errorf("part1(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

func part2(t *testing.T, obs, prog string) (ret int) {
	instructions := make(map[string]func(dst, a, b int, regs []int))
	addInstr := func(name string, f func(dst, a, b int, regs []int)) {
		instructions[name] = f
	}

	addInstr("addr", func(dst, a, b int, regs []int) {
		regs[dst] = regs[a] + regs[b]
	})
	addInstr("addi", func(dst, a, b int, regs []int) {
		regs[dst] = regs[a] + b
	})

	addInstr("mulr", func(dst, a, b int, regs []int) {
		regs[dst] = regs[a] * regs[b]
	})
	addInstr("muli", func(dst, a, b int, regs []int) {
		regs[dst] = regs[a] * b
	})

	addInstr("banr", func(dst, a, b int, regs []int) {
		regs[dst] = regs[a] & regs[b]
	})
	addInstr("bani", func(dst, a, b int, regs []int) {
		regs[dst] = regs[a] & b
	})

	addInstr("borr", func(dst, a, b int, regs []int) {
		regs[dst] = regs[a] | regs[b]
	})
	addInstr("bori", func(dst, a, b int, regs []int) {
		regs[dst] = regs[a] | b
	})

	addInstr("setr", func(dst, a, b int, regs []int) {
		regs[dst] = regs[a]
	})
	addInstr("seti", func(dst, a, b int, regs []int) {
		regs[dst] = a
	})

	addInstr("gtir", func(dst, a, b int, regs []int) {
		if a > regs[b] {
			regs[dst] = 1
		} else {
			regs[dst] = 0
		}
	})
	addInstr("gtri", func(dst, a, b int, regs []int) {
		if regs[a] > b {
			regs[dst] = 1
		} else {
			regs[dst] = 0
		}
	})
	addInstr("gtrr", func(dst, a, b int, regs []int) {
		if regs[a] > regs[b] {
			regs[dst] = 1
		} else {
			regs[dst] = 0
		}
	})

	addInstr("eqir", func(dst, a, b int, regs []int) {
		if a == regs[b] {
			regs[dst] = 1
		} else {
			regs[dst] = 0
		}
	})
	addInstr("eqri", func(dst, a, b int, regs []int) {
		if regs[a] == b {
			regs[dst] = 1
		} else {
			regs[dst] = 0
		}
	})
	addInstr("eqrr", func(dst, a, b int, regs []int) {
		if regs[a] == regs[b] {
			regs[dst] = 1
		} else {
			regs[dst] = 0
		}
	})

	op := make([]map[string]bool, len(instructions))

	// Mark all possible {opcode, instruction} tuples as possible
	for opcode := range op {
		op[opcode] = make(map[string]bool)
		for instr := range instructions {
			op[opcode][instr] = true
		}
	}

	for _, obs := range strings.Split(obs, "\n\n") {
		lines := strings.Split(obs, "\n")

		var r0, r1, r2, r3 int
		advent.Scanner(lines[0]).Extract(t, `Before: \[(-?\d+), (-?\d+), (-?\d+), (-?\d+)\]`, &r0, &r1, &r2, &r3)
		before := []int{r0, r1, r2, r3}

		var opcode, a, b, dst int
		advent.Scanner(lines[1]).Scan(t, &opcode, &a, &b, &dst)

		advent.Scanner(lines[2]).Extract(t, `After:  \[(-?\d+), (-?\d+), (-?\d+), (-?\d+)\]`, &r0, &r1, &r2, &r3)
		after := []int{r0, r1, r2, r3}

		regs := make([]int, 4)
		for instr := range op[opcode] {
			copy(regs, before)
			instructions[instr](dst, a, b, regs)
			if !cmp.Equal(regs, after) {
				delete(op[opcode], instr)
			}
		}
	}

	found := make(map[int]string)
	for {
		var changes, unresolved int
		for opcode, possible := range op {
			if len(possible) > 1 {
				unresolved++
				continue
			}
			var instr string
			for instr = range possible {
			}
			found[opcode] = instr
			for other := range op {
				if opcode != other && op[other][instr] {
					delete(op[other], instr)
					changes++
				}
			}
		}
		t.Logf("Simplify: %d changed, %d unresolved", changes, unresolved)
		if unresolved == 0 {
			break
		}
		if changes == 0 {
			t.Fatalf("Stuck: %v", op)
		}
	}
	for opcode := range op {
		t.Logf("Opcode %d is %q", opcode, found[opcode])
	}
	if len(found) != len(op) {
		t.Fatalf("Failed to resolve!")
	}

	regs := make([]int, 4)
	advent.Lines(prog).Scan(t, func(opcode, a, b, dst int) {
		instructions[found[opcode]](dst, a, b, regs)
	})
	return regs[0]
}

func TestPart2(t *testing.T) {
	tests := []struct {
		name string
		in1  string
		in2  string
		want int
	}{
		{"part2 answer", advent.ReadFile(t, "input1.txt"), advent.ReadFile(t, "input2.txt"), 496},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(t, test.in1, test.in2), test.want; got != want {
				t.Errorf("part2(%#v,%#v)\n = %#v, want %#v", test.in1, test.in2, got, want)
			}
		})
	}
}

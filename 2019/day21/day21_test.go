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

	"github.com/kylelemons/adventofcodesolutions/2019/intcode"
	"github.com/kylelemons/adventofcodesolutions/advent"
)

type Input struct {
	prog *intcode.Program
}

func parseInput(t *testing.T, in string) *Input {
	input := &Input{
		prog: intcode.Compile(t, in),
	}
	return input
}

func part1(t *testing.T, in string) (ret int) {
	input := parseInput(t, in)

	prog := []string{
		// Jump if A, B, or C is False
		"NOT A J", // J = !A
		"NOT B T", // T = !B
		"OR T J",  // J = !A | !B
		"NOT C T", // T = !C
		"OR T J",  // J = !A | !B | !C

		// ALTERNATE:
		// Don't jump if A and B and C
		// "OR A J",
		// "AND B J",
		// "AND C J",
		// "NOT J J",

		// Never jump if we'll land in empty space.
		"AND D J",

		"WALK",
	}
	feed := strings.Join(prog, "\n") + "\n"

	input.prog.Input = func() int {
		ch := feed[0]
		feed = feed[1:]
		return int(ch)
	}

	var line []byte
	input.prog.Output = func(v int) {
		switch {
		case v >= 256:
			ret = v
		case v == '\n':
			t.Logf("%s", line)
			line = nil
		default:
			line = append(line, byte(v))
		}
	}

	input.prog.Run(t)

	return
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part1 answer", advent.ReadFile(t, "input.txt"), 19349939},
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

	// Should NOT jump:
	//   @
	// #####.#.##.#.####
	//    ABCDEFGHI
	// !E & !H = bad
	// !(E | H) = bad
	// E | H = good

	// Should jump:
	//
	//     @
	// #####.#.##.#.####
	//      ABCDEFGHI

	// Should jump:
	//
	//         @
	// #####.#.##.#.####
	//          ABCDEFGHI

	//

	prog := []string{
		// Only jump if we're not stranding ourselves:
		"OR E J", // J = E
		"OR H J", // J = E | H

		// Try to jump if there's a blank in the next 3 spaces:
		"OR A T",  // T = A
		"AND B T", // T = A & B
		"AND C T", // T = A & B & C
		"NOT T T", // T = !(A & B & C)

		// Combine the above two conditions:
		"AND T J", // J = (E | H) & !(A & B & C)

		// Never jump if we'll land in a pit:
		"AND D J",
		"RUN",
	}
	feed := strings.Join(prog, "\n") + "\n"

	input.prog.Input = func() int {
		ch := feed[0]
		feed = feed[1:]
		return int(ch)
	}

	var line []byte
	input.prog.Output = func(v int) {
		switch {
		case v >= 256:
			ret = v
		case v == '\n':
			t.Logf("%s", line)
			line = nil
		default:
			line = append(line, byte(v))
		}
	}

	input.prog.Run(t)

	return
}

func TestPart2(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		// {"part2 example 0", "...", 0},
		{"part2 answer", advent.ReadFile(t, "input.txt"), 1142412777},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(t, test.in), test.want; got != want {
				t.Errorf("part2(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

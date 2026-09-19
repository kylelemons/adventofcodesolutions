// Copyright 2021 Kyle Lemons
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
	"bufio"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/kylelemons/adventofcodesolutions/advent"
)

type Turn struct {
	Dir   int // -1 or +1
	Count int // typically >0
}
type Input struct {
	Lines []Turn
}

var lineRE = regexp.MustCompile(`^([LR])(\d+)$`)

func parseInput(t *testing.T, in string) *Input {
	input := &Input{
		// ...
	}
	lines := bufio.NewScanner(strings.NewReader(in))
	for lines.Scan() {
		line := strings.TrimSpace(lines.Text())
		if line == "" {
			continue
		}
		match := lineRE.FindStringSubmatch(line)
		if len(match) == 0 {
			panic(fmt.Sprintf("Line %q does not match /%s/", line, lineRE))
		}
		dir, count := match[1], match[2]

		var turn Turn
		switch dir {
		case "L":
			turn.Dir = -1
		case "R":
			turn.Dir = +1
		}
		if n, err := strconv.Atoi(count); err != nil {
			panic(fmt.Sprintf("Line %q has invalid count %q", line, count))
		} else {
			turn.Count = n
		}
		input.Lines = append(input.Lines, turn)
	}

	return input
}

func part1(t *testing.T, in string) (ret int) {
	input := parseInput(t, in)

	dial := 50
	for _, turn := range input.Lines {
		dial += turn.Dir * turn.Count
		for dial >= 100 {
			dial -= 100
		}
		for dial < 0 {
			dial += 100
		}
		if dial == 0 {
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
		{"part1 example 0", `L68
L30
R48
L5
R60
L55
L1
L99
R14
L82`, 3},
		{"part1 answer", advent.ReadFile(t, "input.txt"), -1},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(t, test.in), test.want; got != want {
				t.Errorf("part1(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

func part2(t *testing.T, in string) (ret int) {
	input := parseInput(t, in)

	dial := 50
	for _, turn := range input.Lines {
		for range turn.Count {
			dial += turn.Dir
			if dial >= 100 {
				dial -= 100
			}
			if dial < 0 {
				dial += 100
			}
			if dial == 0 {
				ret++
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
		{"example", "L68\nL30\nR48\nL5\nR60\nL55\nL1\nL99\nR14\nL82", 6},
		{"single step to zero (right)", "R50", 1},
		{"single step to zero (left)", "L50", 1},
		{"full rotation passing zero (right)", "R100", 1},
		{"full rotation passing zero (left)", "L100", 1},
		{"multiple rotations (R1000)", "R1000", 10},
		{"multiple rotations (L1000)", "L1000", 10},
		{"multiple rotations landing on zero (R1050)", "R1050", 11},
		{"multiple rotations landing on zero (L1050)", "L1050", 11},
		{"move away from zero and return", "R50\nL5\nR5", 2},
		{"rotations starting from zero", "R50\nR100\nL100", 3},
		{"small steps never reaching zero", "R10\nL20\nR5", 0},
		{"overshooting zero back and forth", "R60\nL20", 2},
		{"part2 answer", advent.ReadFile(t, "input.txt"), 5847},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(t, test.in), test.want; got != want {
				t.Errorf("part1(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

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
	"fmt"
	"strings"
	"testing"

	"github.com/kylelemons/adventofcodesolutions/advent"
)

type Input struct {
	Grid [][]string
}

func parseInput(t *testing.T, in string) *Input {
	input := &Input{
		// ...
	}
	lines := strings.Split(strings.TrimSpace(in), "\n")
	for _, line := range lines {
		input.Grid = append(input.Grid, strings.Split(line, ""))
	}
	return input
}

func (i Input) Print() {
	for _, row := range i.Grid {
		fmt.Println(strings.Join(row, ""))
	}
}

func part1(t *testing.T, in string) (ret int) {
	input := parseInput(t, in)

	set := func(r, c int, s string) {
		if r < 0 || r >= len(input.Grid) || c < 0 || c >= len(input.Grid[r]) {
			return
		}
		input.Grid[r][c] = s
	}

	for col := range input.Grid[0] {
		current := input.Grid[0][col]
		switch current {
		case "S":
			// For convenience, turn the source into a beam
			set(0, col, "|")
		}
	}

	for row := 1; row < len(input.Grid); row++ {
		for col := range input.Grid[row] {
			above := input.Grid[row-1][col]
			current := input.Grid[row][col]
			switch current {
			case "^":
				// A splitter splits the beam if it is under a beam
				if above != "|" {
					continue
				}
				ret++ // count the split
				set(row, col-1, "|")
				set(row, col+1, "|")
			case ".":
				// Empty space carries beams forward
				if above == "|" {
					set(row, col, above)
				}
			}
		}
	}
	input.Print()

	return
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part1 example 0", ".......S.......\n...............\n.......^.......\n...............\n......^.^......\n...............\n.....^.^.^.....\n...............\n....^.^...^....\n...............\n...^.^...^.^...\n...............\n..^...^.....^..\n...............\n.^.^.^.^.^...^.\n...............", 21},
		{"part1 answer", advent.ReadFile(t, "input.txt"), 1651},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part1(t, test.in), test.want; got != want {
				t.Errorf("part1(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

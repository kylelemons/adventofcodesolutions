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
	"slices"
	"strings"
	"testing"

	"github.com/kylelemons/adventofcodesolutions/advent"
)

type Input struct {
	Grid      [][]string
	Timelines [][]int
}

func parseInput(t *testing.T, in string) *Input {
	input := &Input{
		// ...
	}
	lines := strings.Split(strings.TrimSpace(in), "\n")
	for _, line := range lines {
		input.Grid = append(input.Grid, strings.Split(line, ""))
		input.Timelines = append(input.Timelines, make([]int, len(line)))
	}
	return input
}

func sum(s []int) (ret int) {
	for _, v := range s {
		ret += v
	}
	return ret
}

func (i Input) Print() {
	for r, row := range i.Grid {
		fmt.Println(strings.Join(row, ""), i.Timelines[r], sum(i.Timelines[r]))
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

func part2(t *testing.T, in string) (ret int) {
	input := parseInput(t, in)

	setG := func(r, c int, s string) {
		if r < 0 || r >= len(input.Grid) || c < 0 || c >= len(input.Grid[r]) {
			return
		}
		input.Grid[r][c] = s
	}
	addT := func(r, c int, t int) {
		if r < 0 || r >= len(input.Grid) || c < 0 || c >= len(input.Grid[r]) {
			return
		}
		input.Timelines[r][c] += t
	}

	for col := range input.Grid[0] {
		current := input.Grid[0][col]
		switch current {
		case "S":
			// For convenience, turn the source into a beam
			setG(0, col, "|")
			addT(0, col, 1)
		}
	}

	for row := 1; row < len(input.Grid); row++ {
		originalRow := slices.Clone(input.Grid[row])
		for col := range input.Grid[row] {
			above := input.Grid[row-1][col]
			aboveT := input.Timelines[row-1][col]
			switch originalRow[col] { // always switch on the original character, ignore mutations
			case "^":
				// A splitter splits the beam if it is under a beam
				if above != "|" {
					continue
				}
				setG(row, col-1, "|")
				addT(row, col-1, aboveT)
				setG(row, col+1, "|")
				addT(row, col+1, aboveT)
			case ".":
				// Empty space carries beams forward
				if above == "|" {
					setG(row, col, above)
					addT(row, col, aboveT)
				}
			}
		}
	}

	input.Print()

	for _, t := range input.Timelines[len(input.Timelines)-1] {
		ret += t
	}
	return
}

func TestPart2(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part2 example 0", ".......S.......\n...............\n.......^.......\n...............\n......^.^......\n...............\n.....^.^.^.....\n...............\n....^.^...^....\n...............\n...^.^...^.^...\n...............\n..^...^.....^..\n...............\n.^.^.^.^.^...^.\n...............", 40},
		{"part2 answer", advent.ReadFile(t, "input.txt"), 108924003331749},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(t, test.in), test.want; got != want {
				t.Errorf("part2(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

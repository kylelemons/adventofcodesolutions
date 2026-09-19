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
	"strings"
	"testing"

	"github.com/kylelemons/adventofcodesolutions/advent"
)

type Input struct {
	Lines [][]byte
}

func parseInput(t *testing.T, in string) *Input {
	input := &Input{}
	for _, line := range strings.Split(strings.TrimSpace(in), "\n") {
		if line == "" {
			continue
		}
		input.Lines = append(input.Lines, []byte(line))
	}
	return input
}

func part1(t *testing.T, in string) (ret int) {
	input := parseInput(t, in)

	isRoll := func(i, j int) bool {
		if i < 0 || j < 0 {
			return false
		}
		if i >= len(input.Lines) || j >= len(input.Lines[i]) {
			return false
		}
		return input.Lines[i][j] == '@'
	}

	for i := range input.Lines {
		for j := range input.Lines[i] {
			if !isRoll(i, j) {
				continue
			}

			rolls := 0
			for _, di := range []int{-1, 0, 1} {
				for _, dj := range []int{-1, 0, 1} {
					if di == 0 && dj == 0 {
						continue
					}
					if isRoll(i+di, j+dj) {
						rolls++
					}
				}
			}
			if rolls < 4 {
				ret++
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
		{"part1 example 0", "..@@.@@@@.\n@@@.@.@.@@\n@@@@@.@.@@\n@.@@@@..@.\n@@.@@@@.@@\n.@@@@@@@.@\n.@.@.@.@@@\n@.@@@.@@@@\n.@@@@@@@@.\n@.@.@@@.@.\n", 13},
		{"part1 answer", advent.ReadFile(t, "input.txt"), 1602},
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

	isRoll := func(i, j int) bool {
		if i < 0 || j < 0 {
			return false
		}
		if i >= len(input.Lines) || j >= len(input.Lines[i]) {
			return false
		}
		return input.Lines[i][j] == '@'
	}

	remove := func() (removed int) {
		for i := range input.Lines {
			for j := range input.Lines[i] {
				if !isRoll(i, j) {
					continue
				}

				rolls := 0
				for _, di := range []int{-1, 0, 1} {
					for _, dj := range []int{-1, 0, 1} {
						if di == 0 && dj == 0 {
							continue
						}
						if isRoll(i+di, j+dj) {
							rolls++
						}
					}
				}
				if rolls < 4 {
					removed++
					input.Lines[i][j] = 'x'
				}
			}
		}
		return
	}

	for {
		removed := remove()
		if removed == 0 {
			return
		}
		ret += removed
	}
}

func TestPart2(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part2 example 0", "..@@.@@@@.\n@@@.@.@.@@\n@@@@@.@.@@\n@.@@@@..@.\n@@.@@@@.@@\n.@@@@@@@.@\n.@.@.@.@@@\n@.@@@.@@@@\n.@@@@@@@@.\n@.@.@@@.@.\n", 43},
		{"part2 answer", advent.ReadFile(t, "input.txt"), 9518},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(t, test.in), test.want; got != want {
				t.Errorf("part2(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

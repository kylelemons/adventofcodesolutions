// Copyright 2020 Kyle Lemons
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
	"testing"

	"github.com/kylelemons/adventofcodesolutions/advent"
)

type Input struct {
	Trees [][]byte
}

func parseInput(t *testing.T, in string) *Input {
	input := &Input{
		Trees: advent.Split2D(in),
	}
	return input
}

func part1(t *testing.T, in string) (ret int) {
	input := parseInput(t, in)

	dr, dc := 1, 3
	for r, c := 0, 0; r < len(input.Trees); r, c = r+dr, c+dc {
		c = c % len(input.Trees[0])
		if input.Trees[r][c] == '#' {
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
		{"part1 example 0", `..##.......
#...#...#..
.#....#..#.
..#.#...#.#
.#...##..#.
..#.##.....
.#.#.#....#
.#........#
#.##...#...
#...##....#
.#..#...#.#`, 7},
		{"part1 answer", advent.ReadFile(t, "input.txt"), 211},
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

	ret = 1
	type slope struct{ dc, dr int }
	for _, s := range []slope{
		{1, 1},
		{3, 1},
		{5, 1},
		{7, 1},
		{1, 2},
	} {
		hit := 0
		for r, c := 0, 0; r < len(input.Trees); r, c = r+s.dr, c+s.dc {
			c = c % len(input.Trees[0])
			if input.Trees[r][c] == '#' {
				hit++
			}
		}
		ret *= hit
	}
	return
}

func TestPart2(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part2 example 0", `..##.......
#...#...#..
.#....#..#.
..#.#...#.#
.#...##..#.
..#.##.....
.#.#.#....#
.#........#
#.##...#...
#...##....#
.#..#...#.#`, 336},
		{"part2 answer", advent.ReadFile(t, "input.txt"), 3584591857},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(t, test.in), test.want; got != want {
				t.Errorf("part2(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

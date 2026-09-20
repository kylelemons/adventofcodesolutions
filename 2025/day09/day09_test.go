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
	"testing"

	"github.com/kylelemons/adventofcodesolutions/advent"
)

type Tile struct{ R, C int }
type Input struct {
	Tiles []Tile
}

func parseInput(t *testing.T, in string) *Input {
	input := &Input{
		// ...
	}
	advent.Lines(in).Extract(t, `(\d+),(\d+)`, func(r, c int) {
		input.Tiles = append(input.Tiles, Tile{R: r, C: c})
	})
	return input
}

func part1(t *testing.T, in string) (ret int) {
	input := parseInput(t, in)

	for i, t0 := range input.Tiles {
		for _, t1 := range input.Tiles[i+1:] {
			minR := min(t0.R, t1.R)
			minC := min(t0.C, t1.C)
			maxR := max(t0.R, t1.R)
			maxC := max(t0.C, t1.C)
			area := (maxR - minR + 1) * (maxC - minC + 1)
			if area > ret {
				ret = area
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
		{"part1 example 0", "7,1\n11,1\n11,7\n9,7\n9,5\n2,5\n2,3\n7,3", 50},
		{"part1 answer", advent.ReadFile(t, "input.txt"), 4781235324},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part1(t, test.in), test.want; got != want {
				t.Errorf("part1(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

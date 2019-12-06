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
	"testing"

	"github.com/kylelemons/adventofcodesolutions/2019/advent"
)

type coord struct{ x, y int }

func (c coord) String() string {
	return fmt.Sprintf("%d,%d", c.x, c.y)
}

func powerLevel(c coord, serial int) int {
	rack := c.x + 10
	power := rack * c.y
	power += serial
	power *= rack
	power = power / 100 % 10
	power -= 5
	return power
}

func TestPowerLevel(t *testing.T) {
	tests := []struct {
		name string
		grid int
		x, y int
		want int
	}{
		{"powerLevel example 0", 8, 3, 5, 4},
		{"powerLevel example 1", 57, 122, 79, -5},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := powerLevel(coord{test.x, test.y}, test.grid), test.want; got != want {
				t.Errorf("powerLevel({%v,%v}, %v)\n = %#v, want %#v", test.x, test.y, test.grid, got, want)
			}
		})
	}
}

func part1(t *testing.T, serial int) (ret coord) {
	power := advent.Make2DInts(300, 300)
	var best int
	for y := range power {
		for x := range power[y] {
			power[y][x] = powerLevel(coord{x, y}, serial)
		}
	}
	for y := range power {
		if y > len(power)-3 {
			continue
		}
		for x := range power[y] {
			if x > len(power[y])-3 {
				continue
			}
			var sum int
			for yy := 0; yy < 3; yy++ {
				for xx := 0; xx < 3; xx++ {
					sum += power[y+yy][x+xx]
				}
			}
			if sum > best {
				best, ret = sum, coord{x, y}
				t.Logf("New record: %d @ %+v", best, ret)
			}
		}
	}
	return
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name string
		in   int
		want coord
	}{
		{"part1 example 0", 18, coord{33, 45}},
		{"part1 example 0", 42, coord{21, 61}},
		{"part1 answer", 2694, coord{243, 38}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part1(t, test.in), test.want; got != want {
				t.Errorf("part1(%#v)\n = %v, want %v", test.in, got, want)
			}
		})
	}
}

func part2(t *testing.T, serial int) (ret coord, bestSize int) {
	power := advent.Make2DInts(300, 300)
	var best int
	for y := range power {
		for x := range power[y] {
			power[y][x] = powerLevel(coord{x, y}, serial)
		}
	}

	const N = 20 // The max size should really be 300, but that's too slow
	for size := 1; size <= N; size++ {
		for y := range power {
			if y > len(power)-size {
				continue
			}
			for x := range power[y] {
				if x > len(power[y])-size {
					continue
				}
				var sum int
				for _, row := range power[y:][:size] {
					for _, v := range row[x:][:size] {
						sum += v
					}
				}
				if sum > best {
					best, bestSize, ret = sum, size, coord{x, y}
					t.Logf("New record: %dx%d %d @ %+v\n", size, size, best, ret)
				}
			}
		}
	}
	return
}

func TestPart2(t *testing.T) {
	tests := []struct {
		name string
		in   int
		want coord
		size int
	}{
		{"part2 example 0", 18, coord{90, 269}, 16},
		{"part2 example 0", 42, coord{232, 251}, 12},
		{"part2 answer", 2694, coord{235, 146}, 13},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c, s := part2(t, test.in)
			if got, want := c, test.want; got != want {
				t.Errorf("part2(%#v).coord\n = %v, want %v", test.in, got, want)
			}
			if got, want := s, test.size; got != want {
				t.Errorf("part2(%#v).size\n = %v, want %v", test.in, got, want)
			}
		})
	}
}

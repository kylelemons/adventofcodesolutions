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
	"math"
	"sort"
	"testing"

	"github.com/kylelemons/adventofcodesolutions/2019/advent"
)

func make2d(y, x int) [][]byte {
	back := make([]byte, y*x)
	out := make([][]byte, y)
	for r := range out {
		out[r] = back[r*x:][:x:x]
	}
	return out
}

func part1(t *testing.T, in string) (ret int) {
	type Coord struct{ x, y int }
	type State struct {
		origin   int
		coord    Coord
		distance int
		inf      bool
	}
	var level []State
	min, max := Coord{math.MaxInt32, math.MaxInt32}, Coord{math.MinInt32, math.MinInt32}
	advent.Lines(in).Extract(t, `(\d+), (\d+)`, func(x, y int) {
		level = append(level, State{
			origin: len(level),
			coord:  Coord{x, y},
		})
		if x < min.x {
			min.x = x
		}
		if x > max.x {
			max.x = x
		}
		if y < min.y {
			min.y = y
		}
		if y > max.y {
			max.y = y
		}
	})
	origins := len(level)

	visited := make(map[Coord]State)
	for len(level) > 0 {
		var next []State
		for _, s := range level {
			if s.coord.x < min.x || s.coord.x > max.x || s.coord.y < min.y || s.coord.y > max.y {
				s.inf = true
			}
			if prev, ok := visited[s.coord]; ok {
				if s.distance < prev.distance {
					// new closest neighbor found
				} else if s.origin == prev.origin {
					// we hit ourselves, no need to keep going
					continue
				} else if prev.distance == s.distance {
					// new no-man's zone found
					visited[s.coord] = State{
						origin:   -1,
						coord:    s.coord,
						distance: s.distance,
					}
					continue
				} else {
					// we're farther away, stop looking
					continue
				}
			}
			visited[s.coord] = s
			if s.inf {
				continue
			}
			for _, delta := range []Coord{{0, 1}, {0, -1}, {1, 0}, {-1, 0}} {
				next = append(next, State{
					origin:   s.origin,
					coord:    Coord{s.coord.x + delta.x, s.coord.y + delta.y},
					distance: s.distance + 1,
				})
			}
		}
		level = next
	}
	total := make([]int, origins)
	for _, s := range visited {
		if s.origin == -1 {
			continue
		}
		if s.inf || total[s.origin] < 0 {
			total[s.origin] = -1 // discard, basically
			continue
		}
		total[s.origin]++
	}
	sort.Ints(total)
	return total[len(total)-1]
}

func part2(t *testing.T, in string, limit int) (ret int) {
	type Coord struct{ x, y int }
	var starts []Coord
	min, max := Coord{math.MaxInt32, math.MaxInt32}, Coord{math.MinInt32, math.MinInt32}
	advent.Lines(in).Extract(t, `(\d+), (\d+)`, func(x, y int) {
		starts = append(starts, Coord{x, y})
		if x < min.x {
			min.x = x
		}
		if x > max.x {
			max.x = x
		}
		if y < min.y {
			min.y = y
		}
		if y > max.y {
			max.y = y
		}
	})
	// Technically this is necessary for correctness, but not for this input
	// min.x -= limit
	// min.y -= limit
	// max.x += limit
	// max.y += limit

	for y := min.y; y <= max.y; y++ {
		for x := min.x; x <= max.x; x++ {
			var safety int
			for _, start := range starts {
				safety += abs(start.x-x) + abs(start.y-y)
			}
			if safety < limit {
				ret++
			}
		}
	}
	return
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part1 example 0", `1, 1
1, 6
8, 3
3, 4
5, 5
8, 9`, 17},
		{"part1 answer", advent.ReadFile(t, "input.txt"), 3871},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part1(t, test.in), test.want; got != want {
				t.Errorf("part1(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

func TestPart2(t *testing.T) {
	tests := []struct {
		name string
		in   string
		lim  int
		want int
	}{
		{"part2 example 0", `1, 1
1, 6
8, 3
3, 4
5, 5
8, 9`, 32, 16},
		{"part2 answer", advent.ReadFile(t, "input.txt"), 10000, 44667},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(t, test.in, test.lim), test.want; got != want {
				t.Errorf("part2(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

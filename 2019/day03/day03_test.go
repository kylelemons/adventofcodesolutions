// Copyright 2019 Kyle Lemons
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
	"strings"
	"testing"

	"github.com/kylelemons/adventofcodesolutions/advent"
)

type Move struct {
	dir  string // U/L/D/R
	dist int
}

type Pos struct{ x, y int }

var dirs = map[string]Pos{
	"U": {0, 1},
	"D": {0, -1},
	"L": {-1, 0},
	"R": {1, 0},
}

func part1(t *testing.T, in string) (dist int) {
	var wires [][]Move
	advent.Lines(in).Scan(t, func(line string) {
		var steps []Move
		for _, step := range strings.Split(line, ",") {
			var move Move
			advent.Scanner(step).Extract(t, `([ULDR])(\d+)`, &move.dir, &move.dist)
			steps = append(steps, move)
		}
		wires = append(wires, steps)
	})

	visited := make(map[Pos]bool)
	var min int = math.MaxInt32
	for w, wire := range wires[:2] {
		var cur Pos
		for _, move := range wire {
			delta := dirs[move.dir]
			for i := 0; i < move.dist; i++ {
				cur = Pos{cur.x + delta.x, cur.y + delta.y}
				// t.Logf("%v %v", move, cur)
				if w == 0 {
					visited[cur] = true
					continue
				}
				if visited[cur] {
					manhat := abs(cur.x) + abs(cur.y)
					// t.Logf("Intersect at %v: %d", cur, manhat)
					if manhat < min {
						min = manhat
					}
				}
			}
		}
	}
	// t.Logf("Visited %d", len(visited))
	return min
}

func part2(t *testing.T, in string) (dist int) {
	var wires [][]Move
	advent.Lines(in).Scan(t, func(line string) {
		var steps []Move
		advent.Split(line, ',').Extract(t, `([ULDR])(\d+)`, func(dir string, dist int) {
			steps = append(steps, Move{dir, dist})
		})
		wires = append(wires, steps)
	})

	visited := make(map[Pos]int)
	var min int = math.MaxInt32
	for w, wire := range wires[:2] {
		var cur Pos
		var steps int
		for _, move := range wire {
			delta := dirs[move.dir]
			for i := 0; i < move.dist; i++ {
				cur = Pos{cur.x + delta.x, cur.y + delta.y}
				steps++
				if w == 0 {
					if _, ok := visited[cur]; !ok {
						visited[cur] = steps
					}
					continue
				}
				if prev, ok := visited[cur]; ok {
					total := prev + steps
					if total < min {
						min = total
					}
				}
			}
		}
	}
	return min
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part1 example 0", "R8,U5,L5,D3\nU7,R6,D4,L4", 6},
		{"part1 example 1", "R75,D30,R83,U83,L12,D49,R71,U7,L72\nU62,R66,U55,R34,D71,R55,D58,R83", 159},
		{"part1 example 2", "R98,U47,R26,D63,R33,U87,L62,D20,R33,U53,R51\nU98,R91,D20,R16,D67,R40,U7,R15,U6,R7", 135},
		{"part1 answer", advent.ReadFile(t, "input.txt"), 489},
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
		want int
	}{
		{"part2 example 0", "R8,U5,L5,D3\nU7,R6,D4,L4", 30},
		{"part2 example 1", "R75,D30,R83,U83,L12,D49,R71,U7,L72\nU62,R66,U55,R34,D71,R55,D58,R83", 610},
		{"part2 example 2", "R98,U47,R26,D63,R33,U87,L62,D20,R33,U53,R51\nU98,R91,D20,R16,D67,R40,U7,R15,U6,R7", 410},
		{"part2 answer", advent.ReadFile(t, "input.txt"), 93654},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(t, test.in), test.want; got != want {
				t.Errorf("part2(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

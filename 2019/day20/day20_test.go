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
	"log"
	"testing"
	"time"

	"github.com/kylelemons/adventofcodesolutions/advent"
	"github.com/kylelemons/adventofcodesolutions/advent/coords"
)

type Input struct {
	Maze   [][]byte
	Portal map[string][]coords.Coord
	In     map[coords.Coord]string
	Mid    coords.Coord
}

func parseInput(t *testing.T, in string) *Input {
	input := &Input{
		Maze:   advent.Split2D(in),
		Portal: make(map[string][]coords.Coord),
		In:     make(map[coords.Coord]string),
	}
	input.Mid = coords.RC(len(input.Maze), len(input.Maze[0]))

	for r := 0; r < len(input.Maze); r++ {
		for c := 0; c < len(input.Maze[0]); c++ {
			if input.Maze[r][c] != '.' {
				continue
			}
			cur := coords.RC(r, c)
			for _, dir := range coords.Cardinals {
				loc := cur.Add(dir)
				ch := input.Maze[loc.R()][loc.C()]
				if ch >= 'A' && ch <= 'Z' {
					label := input.portalName(cur, dir)
					input.Portal[label] = append(input.Portal[label], cur)
					input.In[loc] = label
				}
			}
		}
	}

	if _, ok := input.Portal["AA"]; !ok {
		t.Fatalf("No AA")
	}
	if _, ok := input.Portal["ZZ"]; !ok {
		t.Fatalf("No ZZ")
	}

	for portal, c := range input.Portal {
		if portal == "AA" || portal == "ZZ" {
			if got, want := len(c), 1; got != want {
				t.Fatalf("start/end portal %q has %q, want %v", portal, got, want)
			}
			continue
		}
		if got, want := len(c), 2; got != want {
			t.Fatalf("portal %q has %q, want %v", portal, got, want)
		}
		continue
	}

	// log.Println(input.Portal)

	return input
}

func (input *Input) portalName(c0 coords.Coord, dir coords.Vector) string {
	a, b := c0.Add(dir), c0.Add(dir.Scale(2))
	switch dir {
	case coords.North, coords.West:
		a, b = b, a
	}
	return string(input.Maze[a.R()][a.C()]) + string(input.Maze[b.R()][b.C()])
}

func (input *Input) Path(from, to string) (pathlen int) {
	visited := make(map[coords.Coord]bool)

	type state struct {
		loc   coords.Coord
		steps int
	}

	tick := time.Tick(1 * time.Second)

	q := []state{{loc: input.Portal[from][0]}}
	for len(q) > 0 {
		cur := q[0]
		q = q[1:]

		select {
		case <-tick:
			log.Println(cur, len(q), len(visited))
		default:
		}
		// log.Println(cur, len(q), len(visited))

		loc := cur.loc
		steps := cur.steps
		if visited[loc] {
			continue
		}
		visited[loc] = true

		ch := input.Maze[loc.R()][loc.C()]
		switch ch {
		case '#':
			continue
		case '.':
		default:
			log.Printf("Unknown ch %q", ch)
			return -1
		}

		for _, dir := range coords.Cardinals {
			next := loc.Add(dir)
			if portal, ok := input.In[next]; ok {
				if portal == to {
					return steps
				}
				if n, ok := input.Portal[portal]; ok && len(n) == 2 {
					next = n[0]
					if loc == next {
						next = n[1]
					}
					log.Printf("Taking portal %q from %v to %v", portal, loc, next)
				} else {
					log.Printf("Unknown portal %q", portal)
					continue
				}
			}
			q = append(q, state{
				loc:   next,
				steps: steps + 1,
			})
		}
	}
	return -1
}

func part1(t *testing.T, in string) (ret int) {
	input := parseInput(t, in)

	return input.Path("AA", "ZZ")
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part1 example 0", `         A           
         A           
  #######.#########  
  #######.........#  
  #######.#######.#  
  #######.#######.#  
  #######.#######.#  
  #####  B    ###.#  
BC...##  C    ###.#  
  ##.##       ###.#  
  ##...DE  F  ###.#  
  #####    G  ###.#  
  #########.#####.#  
DE..#######...###.#  
  #.#########.###.#  
FG..#########.....#  
  ###########.#####  
             Z       
             Z       `, 23},
		{"part1 example 1", `                   A               
                   A               
  #################.#############  
  #.#...#...................#.#.#  
  #.#.#.###.###.###.#########.#.#  
  #.#.#.......#...#.....#.#.#...#  
  #.#########.###.#####.#.#.###.#  
  #.............#.#.....#.......#  
  ###.###########.###.#####.#.#.#  
  #.....#        A   C    #.#.#.#  
  #######        S   P    #####.#  
  #.#...#                 #......VT
  #.#.#.#                 #.#####  
  #...#.#               YN....#.#  
  #.###.#                 #####.#  
DI....#.#                 #.....#  
  #####.#                 #.###.#  
ZZ......#               QG....#..AS
  ###.###                 #######  
JO..#.#.#                 #.....#  
  #.#.#.#                 ###.#.#  
  #...#..DI             BU....#..LF
  #####.#                 #.#####  
YN......#               VT..#....QG
  #.###.#                 #.###.#  
  #.#...#                 #.....#  
  ###.###    J L     J    #.#.###  
  #.....#    O F     P    #.#...#  
  #.###.#####.#.#####.#####.###.#  
  #...#.#.#...#.....#.....#.#...#  
  #.#####.###.###.#.#.#########.#  
  #...#.#.....#...#.#.#.#.....#.#  
  #.###.#####.###.###.#.#.#######  
  #.#.........#...#.............#  
  #########.###.###.#############  
           B   J   C               
           U   P   P               `, 58},
		{"part1 answer", advent.ReadFile(t, "input.txt"), -1},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part1(t, test.in), test.want; got != want {
				t.Errorf("part1(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

func (input *Input) outer(c coords.Coord) (ret bool) {
	height := len(input.Maze)
	width := len(input.Maze[0])
	if c.R() <= 5 || c.C() <= 5 {
		return true
	}
	if c.R() >= height-5 || c.C() >= width-5 {
		return true
	}
	return false
}

func (input *Input) RecursivePath(from, to string) (pathlen int) {
	type viskey struct {
		loc   coords.Coord
		level int
	}
	visited := make(map[viskey]bool)

	type state struct {
		loc   coords.Coord
		steps int
		level int
	}

	tick := time.Tick(1 * time.Second)

	q := []state{{loc: input.Portal[from][0]}}
	for len(q) > 0 {
		cur := q[0]
		q = q[1:]

		select {
		case <-tick:
			log.Printf("%#v, queue %d, visited %d", cur, len(q), len(visited))
		default:
		}
		// log.Println(cur, len(q), len(visited))

		loc := cur.loc
		steps := cur.steps
		if visited[viskey{loc, cur.level}] {
			continue
		}
		visited[viskey{loc, cur.level}] = true

		ch := input.Maze[loc.R()][loc.C()]
		switch ch {
		case '#':
			continue
		case '.':
		default:
			log.Printf("Unknown ch %q", ch)
			return -1
		}

		for _, dir := range coords.Cardinals {
			next := loc.Add(dir)
			level := cur.level
			if portal, ok := input.In[next]; ok {
				if portal == to {
					if cur.level == 0 {
						return steps
					}
					continue
				}
				if cur.level == 0 && input.outer(next) {
					continue
				}
				if input.outer(next) {
					level--
				} else {
					level++
				}
				if n, ok := input.Portal[portal]; ok && len(n) == 2 {
					next = n[0]
					if loc == next {
						next = n[1]
					}
					// log.Printf("Taking portal %q from %v:%v to %v:%v", portal, loc, cur.level, next, level)
				} else {
					// log.Printf("Unknown portal %q", portal)
					continue
				}
			}
			q = append(q, state{
				loc:   next,
				steps: steps + 1,
				level: level,
			})
		}
	}
	return -1
}

func part2(t *testing.T, in string) (ret int) {
	input := parseInput(t, in)

	return input.RecursivePath("AA", "ZZ")
}

func TestPart2(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{

		{"part2 example 1", `             Z L X W       C                 
             Z P Q B       K                 
  ###########.#.#.#.#######.###############  
  #...#.......#.#.......#.#.......#.#.#...#  
  ###.#.#.#.#.#.#.#.###.#.#.#######.#.#.###  
  #.#...#.#.#...#.#.#...#...#...#.#.......#  
  #.###.#######.###.###.#.###.###.#.#######  
  #...#.......#.#...#...#.............#...#  
  #.#########.#######.#.#######.#######.###  
  #...#.#    F       R I       Z    #.#.#.#  
  #.###.#    D       E C       H    #.#.#.#  
  #.#...#                           #...#.#  
  #.###.#                           #.###.#  
  #.#....OA                       WB..#.#..ZH
  #.###.#                           #.#.#.#  
CJ......#                           #.....#  
  #######                           #######  
  #.#....CK                         #......IC
  #.###.#                           #.###.#  
  #.....#                           #...#.#  
  ###.###                           #.#.#.#  
XF....#.#                         RF..#.#.#  
  #####.#                           #######  
  #......CJ                       NM..#...#  
  ###.#.#                           #.###.#  
RE....#.#                           #......RF
  ###.###        X   X       L      #.#.#.#  
  #.....#        F   Q       P      #.#.#.#  
  ###.###########.###.#######.#########.###  
  #.....#...#.....#.......#...#.....#.#...#  
  #####.#.###.#######.#######.###.###.#.#.#  
  #.......#.......#.#.#.#.#...#...#...#.#.#  
  #####.###.#####.#.#.#.#.###.###.#.###.###  
  #.......#.....#.#...#...............#...#  
  #############.#.#.###.###################  
               A O F   N                     
               A A D   M                     `, 396},
		{"part2 answer", advent.ReadFile(t, "input.txt"), 6592},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(t, test.in), test.want; got != want {
				t.Errorf("part2(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

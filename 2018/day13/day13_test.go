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
	"sort"
	"strings"
	"testing"

	"github.com/kylelemons/adventofcodesolutions/2019/advent"
)

type Dir int

const (
	North Dir = iota
	East
	South
	West
)

type Turn int

const (
	Left Turn = iota
	Straight
	Right
	NumTurns
)

type Coord struct{ X, Y int }

func (c Coord) Add(c2 Coord) Coord {
	return Coord{c.X + c2.X, c.Y + c2.Y}
}

func (c Coord) String() string {
	return fmt.Sprintf("%d,%d", c.X, c.Y)
}

type TurnKey struct {
	Dir
	Turn
}

var Turns = map[TurnKey]Dir{
	{North, Left}:     West,
	{East, Left}:      North,
	{South, Left}:     East,
	{West, Left}:      South,
	{North, Straight}: North,
	{East, Straight}:  East,
	{South, Straight}: South,
	{West, Straight}:  West,
	{North, Right}:    East,
	{East, Right}:     South,
	{South, Right}:    West,
	{West, Right}:     North,
}

var Moves = [...]Coord{
	North: {0, -1},
	East:  {1, 0},
	South: {0, 1},
	West:  {-1, 0},
}

var Carts = map[byte]Dir{
	'^': North,
	'>': East,
	'v': South,
	'<': West,
}

func part1(t *testing.T, in string) (ret string) {
	loc := strings.Split(in, "\n")

	type cart struct {
		id    int
		loc   Coord
		dir   Dir
		turns int
	}
	var carts []*cart
	current := map[Coord]*cart{}

	move := func(c *cart, loc Coord) {
		delete(current, c.loc)
		c.loc = loc
		current[loc] = c
	}

	// Find carts
	for y := range loc {
		for x := range loc[y] {
			if dir, ok := Carts[loc[y][x]]; ok {
				// Create cart
				c := &cart{id: len(carts), dir: dir}
				carts = append(carts, c)
				move(c, Coord{x, y})
			}
		}

		// Replace carts with tracks
		loc[y] = strings.NewReplacer("^", "|", "v", "|", ">", "-", "<", "-").Replace(loc[y])
	}

	for tick := 0; tick < 1e6; tick++ {
		moved := make(map[*cart]bool)
		for y := range loc {
			for x := range loc[y] {
				c, ok := current[Coord{x, y}]
				if !ok || moved[c] {
					continue
				}
				moved[c] = true

				next := c.loc.Add(Moves[c.dir])
				if _, ok := current[next]; ok {
					return next.String()
				}
				move(c, next)

				switch loc[next.Y][next.X] {
				case '+':
					// Intersection
					c.dir = Turns[TurnKey{c.dir, Turn(c.turns) % NumTurns}]
					c.turns++
				case '|':
				case '-':
				case '/':
					switch c.dir {
					case North, South:
						c.dir = Turns[TurnKey{c.dir, Right}]
					case East, West:
						c.dir = Turns[TurnKey{c.dir, Left}]
					}
				case '\\':
					switch c.dir {
					case North, South:
						c.dir = Turns[TurnKey{c.dir, Left}]
					case East, West:
						c.dir = Turns[TurnKey{c.dir, Right}]
					}
				}
			}
		}
	}
	panic("tick count exceeded")
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"part1 example 0", strings.Trim(`
|
v
|
|
|
^
|`, "\n"), "0,3"},
		{"part1 example 1", strings.Trim(`
/->-\        
|   |  /----\
| /-+--+-\  |
| | |  | v  |
\-+-/  \-+--/
\------/   `, "\n"), "7,3"},
		{"part1 answer", advent.ReadFile(t, "input.txt"), "45,34"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part1(t, test.in), test.want; got != want {
				t.Errorf("part1(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

func part2(t *testing.T, in string) (ret string) {
	loc := strings.Split(in, "\n")

	type cart struct {
		id    int
		loc   Coord
		dir   Dir
		turns int
	}
	current := map[Coord]*cart{}

	move := func(c *cart, loc Coord) {
		delete(current, c.loc)
		c.loc = loc
		current[loc] = c
	}

	// Find carts
	for y := range loc {
		for x := range loc[y] {
			if dir, ok := Carts[loc[y][x]]; ok {
				// Create cart
				c := &cart{id: len(current), dir: dir}
				move(c, Coord{x, y})
			}
		}

		// Replace carts with tracks
		loc[y] = strings.NewReplacer("^", "|", "v", "|", ">", "-", "<", "-").Replace(loc[y])
	}

	debug := func(tick int) {
		fmt.Printf("Tick %d / %d carts\n", tick, len(current))
		for y := range loc {
			fmt.Printf("%03d: ", y)
			for x := range loc[y] {
				c, ok := current[Coord{x, y}]
				if !ok {
					fmt.Printf("%c", loc[y][x])
					continue
				}
				switch c.dir {
				case North:
					fmt.Print("^")
				case East:
					fmt.Print(">")
				case South:
					fmt.Print("v")
				case West:
					fmt.Print(">")
				}
			}
			fmt.Println()
		}
		fmt.Println()
	}

	for tick := 0; tick < 1e6; tick++ {
		moved := make(map[*cart]bool)
		for y := range loc {
			for x := range loc[y] {
				c, ok := current[Coord{x, y}]
				if !ok || moved[c] {
					continue
				}
				moved[c] = true

				next := c.loc.Add(Moves[c.dir])
				if _, ok := current[next]; ok {
					delete(current, Coord{x, y})
					delete(current, next)
					continue
				}

				if next.Y < 0 || next.X < 0 || next.Y > len(loc) || next.X > len(loc[y]) {
					debug(tick)
					t.Fatalf("Tick %d: Cart %+v drove off edge to %s", tick, c, next)
				}
				move(c, next)

				switch loc[next.Y][next.X] {
				case '+':
					// Intersection
					c.dir = Turns[TurnKey{c.dir, Turn(c.turns) % NumTurns}]
					c.turns++
				case '|':
				case '-':
				case '/':
					switch c.dir {
					case North, South:
						c.dir = Turns[TurnKey{c.dir, Right}]
					case East, West:
						c.dir = Turns[TurnKey{c.dir, Left}]
					}
				case '\\':
					switch c.dir {
					case North, South:
						c.dir = Turns[TurnKey{c.dir, Left}]
					case East, West:
						c.dir = Turns[TurnKey{c.dir, Right}]
					}
				}
			}
		}
		if len(current) == 1 {
			debug(tick)
			for loc := range current {
				return loc.String()
			}
		}
	}
	panic("tick count exceeded")
}

func part2_optimized(t *testing.T, in string) (ret string) {
	loc := strings.Split(in, "\n")

	type cart struct {
		id      int
		loc     Coord
		dir     Dir
		turns   int
		crashed bool
	}
	current := map[Coord]*cart{}

	move := func(c *cart, loc Coord) {
		delete(current, c.loc)
		c.loc = loc
		current[loc] = c
	}

	// Find carts
	for y := range loc {
		for x := range loc[y] {
			if dir, ok := Carts[loc[y][x]]; ok {
				// Create cart
				c := &cart{id: len(current), dir: dir}
				move(c, Coord{x, y})
			}
		}

		// Replace carts with tracks
		loc[y] = strings.NewReplacer("^", "|", "v", "|", ">", "-", "<", "-").Replace(loc[y])
	}

	debug := func(tick int) {
		fmt.Printf("Tick %d / %d carts\n", tick, len(current))
		for y := range loc {
			fmt.Printf("%03d: ", y)
			for x := range loc[y] {
				c, ok := current[Coord{x, y}]
				if !ok {
					fmt.Printf("%c", loc[y][x])
					continue
				}
				switch c.dir {
				case North:
					fmt.Print("^")
				case East:
					fmt.Print(">")
				case South:
					fmt.Print("v")
				case West:
					fmt.Print(">")
				}
			}
			fmt.Println()
		}
		fmt.Println()
	}

	for tick := 0; tick < 1e6; tick++ {
		carts := make([]*cart, 0, len(current))
		for _, current := range current {
			carts = append(carts, current)
		}
		sort.Slice(carts, func(i, j int) bool {
			if a, b := carts[i].loc.Y, carts[j].loc.Y; a != b {
				return a < b
			}
			if a, b := carts[i].loc.X, carts[j].loc.X; a != b {
				return a < b
			}
			return false
		})
		if len(carts) < 1 {
			debug(tick)
			t.Fatalf("No carts left")
		}

		for _, c := range carts {
			if c.crashed {
				continue
			}

			next := c.loc.Add(Moves[c.dir])
			if _, ok := current[next]; ok {
				current[c.loc].crashed = true
				current[next].crashed = true
				delete(current, c.loc)
				delete(current, next)
				continue
			}

			if next.Y < 0 || next.X < 0 || next.Y > len(loc) || next.X > len(loc[0]) {
				debug(tick)
				t.Fatalf("Tick %d: Cart %+v drove off edge to %s", tick, c, next)
			}
			move(c, next)

			switch loc[next.Y][next.X] {
			case '+':
				// Intersection
				c.dir = Turns[TurnKey{c.dir, Turn(c.turns) % NumTurns}]
				c.turns++
			case '|':
			case '-':
			case '/':
				switch c.dir {
				case North, South:
					c.dir = Turns[TurnKey{c.dir, Right}]
				case East, West:
					c.dir = Turns[TurnKey{c.dir, Left}]
				}
			case '\\':
				switch c.dir {
				case North, South:
					c.dir = Turns[TurnKey{c.dir, Left}]
				case East, West:
					c.dir = Turns[TurnKey{c.dir, Right}]
				}
			}
		}
		if len(current) == 1 {
			debug(tick)
			for loc := range current {
				return loc.String()
			}
		}
	}
	debug(-1)
	panic("tick count exceeded")
}

func TestPart2(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"part2 example 0", strings.Trim(`
|
v
|
|
|
^
^`, "\n"), "0,4"},
		{"part2 example 1", strings.Trim(`
/>-<\  
|   |  
| /<+-\
| | | v
\>+</ |
  |   ^
  \<->/`, "\n"), "6,4"},
		{"part2 answer", advent.ReadFile(t, "input.txt"), "91,25"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2_optimized(t, test.in), test.want; got != want {
				t.Errorf("part2(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

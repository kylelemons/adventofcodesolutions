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
	"sort"
	"strings"
	"testing"

	"github.com/kylelemons/adventofcodesolutions/2019/advent"
)

type Coord struct {
	R, C int
}

func (c Coord) Add(c2 Coord) Coord {
	return Coord{c.R + c2.R, c.C + c2.C}
}

func (c Coord) Sub(c0 Coord) int {
	dr := c0.R - c.R
	if dr < 0 {
		dr = -dr
	}
	dc := c0.C - c.C
	if dc < 0 {
		dc = -dc
	}
	return dr + dc
}

// CoordDeltas represent adjacent coordinates in reading order.
var CoordDeltas = []Coord{{-1, 0}, {0, -1}, {0, 1}, {1, 0}}

type byReadingOrder []Coord

func (a byReadingOrder) Len() int      { return len(a) }
func (a byReadingOrder) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byReadingOrder) Less(i, j int) bool {
	if a, b := a[i].R, a[j].R; a != b {
		return a < b
	}
	if a, b := a[i].C, a[j].C; a != b {
		return a < b
	}
	return false
}

type Unit struct {
	Kind byte // G or E
	HP   int  // hit points
	AP   int  // attack power
	Loc  Coord
}

type byLoc []*Unit

func (a byLoc) Len() int      { return len(a) }
func (a byLoc) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byLoc) Less(i, j int) bool {
	if a, b := a[i].Loc.R, a[j].Loc.R; a != b {
		return a < b
	}
	if a, b := a[i].Loc.C, a[j].Loc.C; a != b {
		return a < b
	}
	return false
}

// Map represents the current state of the Day 15 simulation.
type Map struct {
	At     [][]byte
	Units  map[Coord]*Unit
	Deaths []*Unit
}

// NewMap creates a Map for running the Day 15 simulation.
func NewMap(in string, goblinAP, elfAP int) *Map {
	const DefaultHP = 200

	m := &Map{
		At:    advent.Split2D(in),
		Units: make(map[Coord]*Unit),
	}

	// Pull the units out of the map and into the Units field.
	for r := range m.At {
		for c := range m.At[r] {
			switch m.At[r][c] {
			case '#': // wall
			case '.': // floor
			case 'G': // unit
				m.Units[Coord{r, c}] = &Unit{
					Kind: m.At[r][c],
					HP:   DefaultHP,
					AP:   goblinAP,
					Loc:  Coord{r, c},
				}
				m.At[r][c] = '.'
			case 'E': // unit
				m.Units[Coord{r, c}] = &Unit{
					Kind: m.At[r][c],
					HP:   DefaultHP,
					AP:   elfAP,
					Loc:  Coord{r, c},
				}
				m.At[r][c] = '.'
			}
		}
	}

	return m
}

// Move moves the given unit to the given coordinate.
func (m *Map) Move(u *Unit, to Coord) {
	delete(m.Units, u.Loc)
	u.Loc = to
	m.Units[u.Loc] = u
}

func (m *Map) Attack(u *Unit, attackable []*Unit) {
	if len(attackable) == 0 {
		return
	}

	// Sort by attack preference.
	sort.Slice(attackable, func(i, j int) bool {
		// First, sort by HP
		if a, b := attackable[i].HP, attackable[j].HP; a != b {
			return a < b
		}
		// Then sort by direction
		if a, b := attackable[i].Loc.R, attackable[j].Loc.R; a != b {
			return a < b
		}
		if a, b := attackable[i].Loc.C, attackable[j].Loc.C; a != b {
			return a < b
		}
		return false
	})

	// Attack the preferred unit.
	target := attackable[0]
	target.HP -= u.AP
	if target.HP <= 0 {
		delete(m.Units, target.Loc)
		m.Deaths = append(m.Deaths, target)
	}
}

// NextMove performs a breadth-first-search outward from u until it hits one or
// more destination in the map, chooses the preferred destination, and returns
// the preferred first step to get there.
func (m *Map) NextMove(u *Unit, dest map[Coord]bool) Coord {
	// Our BFS needs the current location and the path to get there.
	type state struct {
		loc  Coord
		path []Coord
	}

	// Initialize the BFS with the current location.
	q := []state{{loc: u.Loc}}

	// Track visited locations so we don't do extra work or get into cycles.
	visited := make(map[Coord]bool)

	// If we get to a distance where we start finding destinations, collect them
	// so that we can determine which one is preferred.
	var best []state

	// Run the BFS in chunks, so that we consider all destionations a certain
	// number of steps away together.
	//
	// The BFS stops when we run out of places to try (in q) or we find paths to
	// destinations (best).
	for steps := 0; len(best) == 0 && len(q) > 0; steps++ {
		next := make([]state, 0, 4*len(q))
		for _, cur := range q {
			if visited[cur.loc] {
				continue
			}
			visited[cur.loc] = true

			// Move from the current location to each of the adjacent spaces.
			for _, d := range CoordDeltas {
				c := cur.loc.Add(d)
				if m.At[c.R][c.C] != '.' { // wall
					continue
				}
				if _, ok := m.Units[c]; ok { // already occupied
					continue
				}

				// Store the next state.
				//
				// Note that the path limts the capacity so that states never
				// overwrite one another's backing storage.
				s := state{
					loc:  c,
					path: append(cur.path[:len(cur.path):len(cur.path)], c),
				}
				next = append(next, s)

				// If we've found a destination, our BFS is over, but we may not
				// be at the first "reading order" destination, so we collect
				// all of them and then sort later.
				if dest[c] {
					best = append(best, s)
				}
			}
		}
		q = next
	}
	if len(best) == 0 {
		// Nobody is reachable, stay put.
		return u.Loc
	}

	// Pick the best path based on reading order of destination and first step.
	sort.Slice(best, func(i, j int) bool {
		// First, sort by destination
		if a, b := best[i].loc.R, best[j].loc.R; a != b {
			return a < b
		}
		if a, b := best[i].loc.C, best[j].loc.C; a != b {
			return a < b
		}
		// Then sort by first step
		if a, b := best[i].path[0].R, best[j].path[0].R; a != b {
			return a < b
		}
		if a, b := best[i].path[0].C, best[j].path[0].C; a != b {
			return a < b
		}
		return false
	})

	// The best step is the first one in the first path.
	return best[0].path[0]
}

func (m *Map) Combat() (rounds, goblinHP, elfHP int) {
	for round := 0; ; round++ {
		// Figure out the order in which units are going to act.
		ready := make([]*Unit, 0, len(m.Units))
		for _, u := range m.Units {
			ready = append(ready, u)
		}
		sort.Sort(byLoc(ready)) // sorted in "reading" order

		// Make each unit act in order (unless they are killed first).
		for _, u := range ready {
			if u.HP <= 0 {
				continue
			}

			// Store attackable units and potential destinations.
			attackable := make([]*Unit, 0, len(m.Units))
			dest := make(map[Coord]bool, 4*len(m.Units))

			// Count the enemies we see, in case combat ends.
			var enemies int

			// Find enemy units
			for _, e := range m.Units {
				if e.Kind == u.Kind {
					continue // not an enemy
				}
				enemies++

				// Check if the enemy is close enough to attack directly.
				if e.Loc.Sub(u.Loc) == 1 {
					attackable = append(attackable, e)
					continue
				}

				// Collect squares we can stand in to attack this enemy.
				for _, d := range CoordDeltas {
					c := e.Loc.Add(d)
					if m.At[c.R][c.C] != '.' { // wall
						continue
					}
					if _, ok := m.Units[c]; ok { // occupied already
						continue
					}
					dest[c] = true
				}
			}
			if enemies == 0 {
				// No enemies, we're done
				for _, u := range m.Units {
					switch u.Kind {
					case 'G':
						goblinHP += u.HP
					case 'E':
						elfHP += u.HP
					}
				}
				return round, goblinHP, elfHP
			} else if len(attackable) > 0 {
				m.Attack(u, attackable)
			} else if len(dest) > 0 {
				m.Move(u, m.NextMove(u, dest))
				for _, d := range CoordDeltas {
					c := u.Loc.Add(d)
					e, ok := m.Units[c]
					if !ok {
						continue
					}
					if u.Kind == e.Kind {
						continue
					}
					attackable = append(attackable, e)
				}
				m.Attack(u, attackable)
			}
		}
	}
}

func part1(t *testing.T, in string) (ret int) {
	m := NewMap(in, 3, 3)
	rounds, goblinHP, elfHP := m.Combat()
	return (goblinHP + elfHP) * rounds
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part1 example 0", strings.Trim(`
#######   
#.G...#   
#...EG#   
#.#.#G#   
#..G#E#   
#.....#   
#######   
`, "\n"), 27730},
		{"part1 example 1", strings.Trim(`
#######   
#G..#E#   
#E#E.E#   
#G.##.#   
#...#E#   
#...E.#   
#######   `, "\n"), 36334},
		{"part1 example 2", strings.Trim(`
#######
#E..EG#
#.#G.E#
#E.##E#
#G..#.#
#..E#.#
#######
`, "\n"), 39514},
		{"part1 example 3", strings.Trim(`
#######
#.E...#
#.#..G#
#.###.#
#E#G#G#
#...#G#
#######
`, "\n"), 28944},
		{"part1 example 4", strings.Trim(`
#########
#G......#
#.E.#...#
#..##..G#
#...##..#
#...#...#
#.G...G.#
#.....G.#
#########
`, "\n"), 18740},
		{"part1 answer", advent.ReadFile(t, "input.txt"), 239010},
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
moreAP:
	for elfAP := 3; ; elfAP++ {
		m := NewMap(in, 3, elfAP)
		rounds, _, elfHP := m.Combat()

		// If any elves died, try again with more AP.
		for _, death := range m.Deaths {
			if death.Kind == 'E' {
				continue moreAP
			}
		}

		// No elves died, we're good.
		return elfHP * rounds
	}
}

func TestPart2(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part2 example 0", strings.Trim(`
#######   
#.G...#   
#...EG#   
#.#.#G#   
#..G#E#   
#.....#   
#######   
`, "\n"), 4988},
		{"part2 example 2", strings.Trim(`
#######
#E..EG#
#.#G.E#
#E.##E#
#G..#.#
#..E#.#
#######
`, "\n"), 31284},
		{"part2 example 4", strings.Trim(`
#########
#G......#
#.E.#...#
#..##..G#
#...##..#
#...#...#
#.G...G.#
#.....G.#
#########
`, "\n"), 1140},
		{"part2 answer", advent.ReadFile(t, "input.txt"), 62468},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(t, test.in), test.want; got != want {
				t.Errorf("part2(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

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

func part1(t *testing.T, in string) (ret int) {
	const (
		DefaultHP = 200
		DefaultAP = 3
	)

	field := advent.Split2D(in)

	units := make(map[Coord]*Unit)
	move := func(u *Unit, to Coord) {
		delete(units, u.Loc)
		u.Loc = to
		units[u.Loc] = u
	}

	for r := range field {
		for c := range field[r] {
			switch field[r][c] {
			case '#': // wall
			case '.': // floor
			case 'G', 'E': // unit
				units[Coord{r, c}] = &Unit{
					Kind: field[r][c],
					HP:   DefaultHP,
					AP:   DefaultAP,
					Loc:  Coord{r, c},
				}
				field[r][c] = '.'
			}
		}
	}

	for round := 0; round < 1e6; round++ {
		ready := make([]*Unit, 0, len(units))
		for _, u := range units {
			ready = append(ready, u)
		}
		sort.Sort(byLoc(ready))
		for _, u := range ready {
			if u.HP <= 0 {
				continue
			}

			attackable := make([]*Unit, 0, len(units))
			dest := make(map[Coord]bool, 4*len(units))
			var enemies int
			for _, e := range units {
				if e.Kind == u.Kind {
					continue // not an enemy
				}
				enemies++
				if e.Loc.Sub(u.Loc) == 1 {
					attackable = append(attackable, e)
					continue
				}
				for _, d := range []Coord{{-1, 0}, {0, -1}, {0, 1}, {1, 0}} {
					c := e.Loc.Add(d)
					if field[c.R][c.C] != '.' {
						continue
					}
					if _, ok := units[c]; ok {
						continue
					}
					dest[c] = true
				}
			}
			attack := func() {
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
				target := attackable[0]
				target.HP -= u.AP
				if target.HP <= 0 {
					delete(units, target.Loc)
				}
			}
			if enemies == 0 {
				// No enemies, we're done
				for _, u := range units {
					ret += u.HP
				}
				return ret * round
			} else if len(attackable) > 0 {
				attack()
			} else if len(dest) > 0 {
				type state struct {
					loc  Coord
					path []Coord
				}
				q := []state{{loc: u.Loc}}
				visited := make(map[Coord]bool)
				var best []state
				for steps := 0; len(best) == 0 && len(q) > 0; steps++ {
					// fmt.Printf("Round %d: %+v: Step %d: %d starts\n", round, u, steps, len(q))
					next := make([]state, 0, 4*len(q))
					for _, cur := range q {
						if visited[cur.loc] {
							continue
						}
						visited[cur.loc] = true
					nextDir:
						for _, d := range []Coord{{-1, 0}, {0, -1}, {0, 1}, {1, 0}} {
							c := cur.loc.Add(d)
							if field[c.R][c.C] != '.' {
								continue
							}
							if _, ok := units[c]; ok {
								continue
							}
							for _, prev := range cur.path {
								if prev == c {
									continue nextDir
								}
							}
							s := state{
								loc:  c,
								path: append(cur.path[:len(cur.path):len(cur.path)], c),
							}
							next = append(next, s)
							if dest[c] {
								best = append(best, s)
							}
						}
					}
					q = next
				}
				if len(best) == 0 {
					continue // no enemies reachable
				}
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
				move(u, best[0].path[0])
				for _, d := range []Coord{{-1, 0}, {0, -1}, {0, 1}, {1, 0}} {
					c := u.Loc.Add(d)
					e, ok := units[c]
					if !ok {
						continue
					}
					if u.Kind == e.Kind {
						continue
					}
					attackable = append(attackable, e)
				}
				if len(attackable) > 0 {
					attack()
				}
			}
		}

		// if round < 50 {
		// } else if round < 100 {
		// 	fmt.Printf("Finished round %d\n", round)
		// } else {
		// 	fmt.Printf("After %d round(s):\n", round+1)
		// 	for r := range field {
		// 		hp := ""
		// 		for c := range field[r] {
		// 			if u, ok := units[Coord{r, c}]; ok {
		// 				fmt.Printf("%c", u.Kind)
		// 				hp += fmt.Sprintf(" %c:%d", u.Kind, u.HP)
		// 				continue
		// 			}
		// 			fmt.Printf("%c", field[r][c])
		// 		}
		// 		fmt.Printf("%s\n", hp)
		// 	}
		// 	fmt.Println()
		// }
	}
	panic("round limit exceeded")
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
	const (
		DefaultHP = 200
		DefaultAP = 3
	)

	var ElfAP = 2
tryAgain:

	field := advent.Split2D(in)

	units := make(map[Coord]*Unit)
	move := func(u *Unit, to Coord) {
		delete(units, u.Loc)
		u.Loc = to
		units[u.Loc] = u
	}

	var elves int
	for r := range field {
		for c := range field[r] {
			switch field[r][c] {
			case '#': // wall
			case '.': // floor
			case 'G', 'E': // unit
				ap := DefaultAP
				if field[r][c] == 'E' {
					ap = ElfAP
					elves++
				}
				units[Coord{r, c}] = &Unit{
					Kind: field[r][c],
					HP:   DefaultHP,
					AP:   ap,
					Loc:  Coord{r, c},
				}
				field[r][c] = '.'
			}
		}
	}
	debug := func() string {
		out := new(strings.Builder)
		for r := range field {
			hp := ""
			for c := range field[r] {
				if u, ok := units[Coord{r, c}]; ok {
					fmt.Fprintf(out, "%c", u.Kind)
					hp += fmt.Sprintf(" %c:%d", u.Kind, u.HP)
					continue
				}
				fmt.Fprintf(out, "%c", field[r][c])
			}
			fmt.Fprintf(out, "%s\n", hp)
		}
		return out.String()
	}

	for round := 0; round < 1e6; round++ {
		ready := make([]*Unit, 0, len(units))
		for _, u := range units {
			ready = append(ready, u)
		}
		sort.Sort(byLoc(ready))
		for _, u := range ready {
			if u.HP <= 0 {
				continue
			}

			attackable := make([]*Unit, 0, len(units))
			dest := make(map[Coord]bool, 4*len(units))
			var enemies int
			for _, e := range units {
				if e.Kind == u.Kind {
					continue // not an enemy
				}
				enemies++
				if e.Loc.Sub(u.Loc) == 1 {
					attackable = append(attackable, e)
					continue
				}
				for _, d := range []Coord{{-1, 0}, {0, -1}, {0, 1}, {1, 0}} {
					c := e.Loc.Add(d)
					if field[c.R][c.C] != '.' {
						continue
					}
					if _, ok := units[c]; ok {
						continue
					}
					dest[c] = true
				}
			}
			attack := func() {
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
				target := attackable[0]
				target.HP -= u.AP
				if target.HP <= 0 {
					delete(units, target.Loc)
				}
			}
			if enemies == 0 {
				// No enemies, we're done
				hp := 0
				for _, u := range units {
					// If the last enemy is a goblin, try again
					if u.Kind == 'G' {
						ElfAP++
						t.Logf("Goblins win; Increasing Elf AP to %d", ElfAP)
						goto tryAgain
					}
					hp += u.HP
				}
				if len(units) < elves {
					ElfAP++
					t.Logf("Lost elves; Increasing Elf AP to %d", ElfAP)
					goto tryAgain
				}
				t.Logf("Elves win with %d total HP left after %d full rounds:\n%s", hp, round, debug())
				return hp * round
			} else if len(attackable) > 0 {
				attack()
			} else if len(dest) > 0 {
				type state struct {
					loc  Coord
					path []Coord
				}
				q := []state{{loc: u.Loc}}
				visited := make(map[Coord]bool)
				var best []state
				for steps := 0; len(best) == 0 && len(q) > 0; steps++ {
					// fmt.Printf("Round %d: %+v: Step %d: %d starts\n", round, u, steps, len(q))
					next := make([]state, 0, 4*len(q))
					for _, cur := range q {
						if visited[cur.loc] {
							continue
						}
						visited[cur.loc] = true
					nextDir:
						for _, d := range []Coord{{-1, 0}, {0, -1}, {0, 1}, {1, 0}} {
							c := cur.loc.Add(d)
							if field[c.R][c.C] != '.' {
								continue
							}
							if _, ok := units[c]; ok {
								continue
							}
							for _, prev := range cur.path {
								if prev == c {
									continue nextDir
								}
							}
							s := state{
								loc:  c,
								path: append(cur.path[:len(cur.path):len(cur.path)], c),
							}
							next = append(next, s)
							if dest[c] {
								best = append(best, s)
							}
						}
					}
					q = next
				}
				if len(best) == 0 {
					continue // no enemies reachable
				}
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
				move(u, best[0].path[0])
				for _, d := range []Coord{{-1, 0}, {0, -1}, {0, 1}, {1, 0}} {
					c := u.Loc.Add(d)
					e, ok := units[c]
					if !ok {
						continue
					}
					if u.Kind == e.Kind {
						continue
					}
					attackable = append(attackable, e)
				}
				if len(attackable) > 0 {
					attack()
				}
			}
		}
	}
	panic("round limit exceeded")
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

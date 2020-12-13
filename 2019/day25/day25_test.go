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
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/kylelemons/adventofcodesolutions/2019/intcode"
	"github.com/kylelemons/adventofcodesolutions/advent"
	"github.com/kylelemons/adventofcodesolutions/advent/coords"
)

type Input struct {
	Base *intcode.Program
}

func parseInput(t *testing.T, in string) *Input {
	input := &Input{
		Base: intcode.Compile(t, in),
	}
	return input
}

func part1(t *testing.T, in string) (passcode int) {
	input := parseInput(t, in)

	type state struct {
		loc   coords.Coord
		prog  *intcode.Program
		items []string
		prev  []coords.Coord
	}

	q := []state{{prog: input.Base}}

	itemLocations := make(map[string]coords.Coord)
	display := make(map[coords.Coord]byte)

	type vk struct {
		loc   coords.Coord
		items string // item1|item2 (sorted)
	}
	vkFor := func(loc coords.Coord, items []string) vk {
		sort.Strings(items)
		return vk{
			loc:   loc,
			items: strings.Join(items, "|"),
		}
	}
	visited := make(map[vk]bool)
	seenRooms := make(map[string]bool)

	doorDirs := map[string]coords.Coord{
		"north": coords.North,
		"east":  coords.East,
		"west":  coords.West,
		"south": coords.South,
	}

	type scanMode int
	const (
		startMode scanMode = iota
		doorMode
		itemMode
		seenMode
	)

	for i := 0; i < 1000 && len(q) > 0; i++ {
		s := q[0]
		q = q[1:]

		key := vkFor(s.loc, s.items)
		if visited[key] {
			continue
		}
		visited[key] = true

		prefix := fmt.Sprintf("%s|%d|%s", s.loc, len(s.prev), strings.Join(s.items, "|"))

		// Run the program to figure out what's in this room
		var (
			mode     scanMode
			doors    []string
			items    []string
			rejected bool // rejected by a security checkpoint
		)
		s.prog.ASCII(
			func() string { t.Fatalf("Not ready for input"); return "" },
			func(v string) {
				line := advent.Scanner(v)
				switch {
				case v == "":
				case v == "Command?":
					s.prog.Halt()
				case v == "Doors here lead:":
					mode = doorMode
				case v == "Items here:":
					mode = itemMode
				case v == "Analyzing...":
					// t.Logf("%s : Trying security checkpoint", prefix)
				case line.CanExtract(t, `.* (\d+) .*airlock.*`, &passcode):
					passcode = passcode
				case strings.Contains(v, "heavier than the detected value"), strings.Contains(v, "lighter than the detected value"):
					prev := s.prev[len(s.prev)-1]
					// t.Logf("%s : Rejected, returning to %v", prefix, prev)
					display[s.loc] = 'N'
					s.loc = prev
					rejected = true
				case mode == doorMode && strings.HasPrefix(v, "- "):
					doors = append(doors, v[2:])
				case mode == itemMode && strings.HasPrefix(v, "- "):
					items = append(items, v[2:])
				case mode == seenMode:
				case strings.HasPrefix(v, "=="):
					if seenRooms[v] {
						mode = seenMode
					} else {
						seenRooms[v] = true
						t.Logf("%s : %s", prefix, v)
					}
				default:
					t.Logf("%s : %s", prefix, v)
				}
			},
			func(v int) { t.Fatalf("Unexpected Output(%v)", v) },
		)
		s.prog.Run(t)
		if passcode > 0 {
			t.Logf("Santa gave us the passcode: %v", passcode)
			return
		}

		if rejected {
			t.Logf("%s : Rejected carrying %q", prefix, s.items)
			continue
		}

		if _, ok := display[s.loc]; !ok {
			display[s.loc] = '.'
			if len(items) > 0 {
				display[s.loc] = 'o'
			}
		}

		for _, item := range items {
			switch item {
			case "infinite loop", "photons", "molten lava", "escape pod", "giant electromagnet":
				continue // well played
			}

			if _, ok := itemLocations[item]; !ok {
				itemLocations[item] = s.loc
			}

			// Provide the program with the input to pick up the item
			s.prog.ASCII(
				func() string {
					return "take " + item
				},
				func(v string) {
					switch v {
					case "Command?":
						s.prog.Halt()
					case "":
					default:
						if v == "You take the "+item+"." {
							break
						}
						t.Logf("%s : (item) %q", prefix, v)
					}
				},
				func(v int) { t.Fatalf("Unexpected Output(%v)", v) },
			)
			s.prog.Run(t)
			s.items = append(s.items[:len(s.items):len(s.items)], item)
		}

		for _, door := range doors {
			dir := doorDirs[door]

			// Provide the program with the correct input
			p := s.prog.Snapshot()
			p.ASCII(
				func() string {
					// t.Logf("%s : Sending %q", prefix, door)
					return door
				},
				func(v string) { p.Halt() },
				func(v int) { t.Fatalf("Unexpected Output(%v)", v) },
			)
			p.Run(t)
			q = append(q, state{
				loc:   s.loc.Add(dir),
				prog:  p,
				items: s.items,
				prev:  append(s.prev[:len(s.prev):len(s.prev)], s.loc),
			})
		}
	}

	t.Logf("%d remaining in queue", len(q))
	t.Logf("Map:\n%s", advent.String2DMap(display))
	for item, loc := range itemLocations {
		t.Logf("Item %q @ %v", item, loc)
	}

	return
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		// {"part1 example 0", "...", 0},
		{"part1 answer", advent.ReadFile(t, "input.txt"), 25166400},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part1(t, test.in), test.want; got != want {
				t.Errorf("part1(...)\n = %#v, want %#v", got, want)
			}
		})
	}
}

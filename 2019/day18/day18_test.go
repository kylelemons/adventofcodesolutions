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
	"time"

	"github.com/kylelemons/adventofcodesolutions/advent"
	"github.com/kylelemons/adventofcodesolutions/advent/coords"
)

type Input struct {
	Maze [][]byte
}

func parseInput(t *testing.T, in string) *Input {
	input := &Input{
		Maze: advent.Split2D(strings.TrimSpace(in)),
	}
	return input
}

func part1(t *testing.T, in string) (ret int) {
	input := parseInput(t, in)

	keyLocations := make(map[string]coords.Coord)
	var keys []string

	for r, row := range input.Maze {
		for c, ch := range row {
			if ch == '@' || ch >= 'a' && ch <= 'z' {
				keyLocations[string(ch)] = coords.RC(r, c)
				keys = append(keys, string(ch))
			}
		}
	}
	sort.Strings(keys)

	// t.Logf("Finding paths between keys...")

	type pathKey struct {
		start, end string
	}
	type path struct {
		steps int
		doors []string
	}
	paths := make(map[pathKey][]path)

	for _, startKey := range keys {
		startLoc := keyLocations[startKey]

		visited := make(map[coords.Coord]bool)

		type state struct {
			cur   coords.Coord
			steps int
			doors []string
		}
		q := []state{{cur: startLoc}}
		for len(q) > 0 {
			s := q[0]
			q = q[1:]

			cur, steps, doors := s.cur, s.steps, s.doors
			if visited[cur] {
				continue
			}
			visited[cur] = true

			r, c := s.cur.R(), s.cur.C()
			ch := input.Maze[r][c]
			switch {
			case ch == '#':
				// wall
				continue
			case ch >= 'A' && ch <= 'Z':
				// door
				doors = plus(doors, string(ch-'A'+'a'))
				sort.Strings(doors)
			case ch >= 'a' && ch <= 'z':
				// key
				key := pathKey{startKey, string(ch)}
				paths[key] = append(paths[key], path{
					steps: steps,
					doors: doors,
				})
				// t.Logf("Path from %q to %q: %d steps through %q", startKey, ch, steps, doors)
			}

			for _, dir := range coords.Cardinals {
				q = append(q, state{
					cur:   cur.Add(dir),
					steps: steps + 1,
					doors: doors,
				})
			}
		}
	}

	// t.Logf("Finding best possible route...")

	type state struct {
		curKey string
		steps  int
		keys   []string
	}
	var q advent.PriorityQueue
	q.Push(
		state{curKey: "@", keys: []string{"@"}},
		0,  // steps
		-1, // -keys
	)

	status := time.NewTicker(1 * time.Second)

	type visitedKey struct {
		loc  coords.Coord
		keys string
	}
	visited := make(map[visitedKey]int)
	for q.Len() > 0 {
		var s state
		q.Pop(&s)
		if len(s.keys) == len(keys) {
			return s.steps
		}

		vk := visitedKey{keyLocations[s.curKey], strings.Join(s.keys, "")}
		if prev, ok := visited[vk]; ok && prev <= s.steps {
			continue
		}
		visited[vk] = s.steps

		select {
		case <-status.C:
			t.Logf("%#v (%d queued)", s, q.Len())
		default:
		}

		for _, nextKey := range keys {
			if contains(s.keys, nextKey) {
				continue
			}

			var best path
		nextPath:
			for _, path := range paths[pathKey{s.curKey, nextKey}] {
				for _, required := range path.doors {
					if !contains(s.keys, required) {
						continue nextPath
					}
				}
				if best.steps == 0 || path.steps < best.steps {
					best = path
				}
			}
			if best.steps == 0 {
				// no available paths
				continue
			}

			keys := plus(s.keys, nextKey)
			sort.Strings(keys)
			steps := s.steps + best.steps

			nk := visitedKey{keyLocations[nextKey], strings.Join(keys, "")}
			if prev, ok := visited[nk]; ok && prev <= steps {
				continue
			}

			q.Push(
				state{
					curKey: nextKey,
					steps:  steps,
					keys:   keys,
				},
				steps,
				-len(s.keys),
			)
		}
	}

	return -1
}

func plus(slice []string, s string) []string {
	return append(slice[:len(slice):len(slice)], s)
}

func contains(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part1 example 0", `#########
#b.A.@.a#
#########`, 8},
		{"part1 example 1", `########################
#f.D.E.e.C.b.A.@.a.B.c.#
######################.#
#d.....................#
########################`, 86},
		{"part1 example 2", `########################
#...............b.C.D.f#
#.######################
#.....@.a.B.c.d.A.e.F.g#
########################`, 132},
		{"part1 example 3", `#################
#i.G..c...e..H.p#
########.########
#j.A..b...f..D.o#
########@########
#k.E..a...g..B.n#
########.########
#l.F..d...h..C.m#
#################`, 136},
		{"part1 example 4", `########################
#@..............ac.GI.b#
###d#e#f################
###A#B#C################
###g#h#i################
########################`, 81},
		{"part1 answer", advent.ReadFile(t, "input.txt"), 4676},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part1(t, test.in), test.want; got != want {
				t.Errorf("part1(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

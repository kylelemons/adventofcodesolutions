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
	"log"
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

func part2(t *testing.T, in string) (ret int) {
	input := parseInput(t, in)

	// Modify the input to have the four starting positions.
modify:
	for r, row := range input.Maze {
		for c, ch := range row {
			if ch == '@' {
				center := coords.RC(r, c)
				walls := []coords.Coord{
					coords.North, coords.East, coords.South, coords.West,
					coords.RC(0, 0),
				}
				starts := []coords.Coord{
					coords.NorthEast, coords.NorthWest,
					coords.SouthWest, coords.SouthEast,
				}
				for _, delta := range walls {
					c := center.Add(delta)
					input.Maze[c.R()][c.C()] = '#'
				}
				for i, delta := range starts {
					c := center.Add(delta)
					input.Maze[c.R()][c.C()] = '1' + byte(i)
				}
				break modify
			}
		}
	}

	// Find all of the keys and note their locations.
	keyLocations := make(map[string]coords.Coord)
	var keys []string
	for r, row := range input.Maze {
		for c, ch := range row {
			if ch >= '1' && ch <= '4' || ch >= 'a' && ch <= 'z' {
				keyLocations[string(ch)] = coords.RC(r, c)
				keys = append(keys, string(ch))
			}
		}
	}
	sort.Strings(keys)

	// Prepare a memoized BFS for finding distances to all keys.
	rows := len(input.Maze)
	cols := len(input.Maze[0])
	distVisited := make([]bool, rows*cols)
	type memoKey struct {
		from string
		has  string
	}
	distMemo := make(map[memoKey]map[string]int, len(keys)*len(keys)*len(keys))
	distances := func(from string, has []string) (ret map[string]int) {
		mk := memoKey{from, strings.Join(has, "")}
		if memo, ok := distMemo[mk]; ok {
			return memo
		}
		defer func() { distMemo[mk] = ret }()

		for i := range distVisited {
			distVisited[i] = false
		}

		key2dist := make(map[string]int)
		type state struct {
			cur   coords.Coord
			steps int
		}
		q := make([]state, 0, rows*cols)
		q = append(q, state{keyLocations[from], 0})
		for len(q) > 0 {
			c := q[0]
			q = q[1:]

			vk := c.cur.R()*cols + c.cur.C()
			if distVisited[vk] {
				continue
			}
			distVisited[vk] = true

			switch ch := c.cur.In2D(input.Maze); {
			case ch == '#':
				continue
			case ch >= 'A' && ch <= 'Z':
				key := string(ch - 'A' + 'a')
				if !contains(has, key) {
					continue
				}
			case ch >= 'a' && ch <= 'z':
				key := string(ch)
				if !contains(has, key) {
					key2dist[key] = c.steps
				}
			}

			for _, delta := range coords.Cardinals {
				next := c.cur.Add(delta)
				q = append(q, state{
					cur:   next,
					steps: c.steps + 1,
				}) // save space TPDP
			}
		}
		return key2dist
	}

	// Print out the map for humans to appreciate.
	for _, row := range input.Maze {
		t.Logf("  %s", row)
	}

	// Don't visit states that have the same starting locations and keys.
	type visKey struct {
		bots [4]string
		keys string
	}
	visited := make(map[visKey]bool)

	// A* Setup.
	type state struct {
		bots  [4]string
		steps int
		keys  []string
	}
	var q advent.PriorityQueue

	start := state{
		bots: [4]string{"1", "2", "3", "4"},
		keys: []string{"1", "2", "3", "4"},
	}
	q.Push(start, start.steps)

	// Provide some friendly output so humans can tell it's working.
	var best state
	status := time.NewTicker(1 * time.Second)
	defer status.Stop()

	// Run the A*.
	for q.Len() > 0 {
		var s state
		q.Pop(&s)

		// Visited check.
		vk := visKey{s.bots, strings.Join(s.keys, "")}
		if visited[vk] {
			continue
		}
		visited[vk] = true

		// Human state tracing.
		select {
		case <-status.C:
			log.Printf("Queue: %d | %d remaining | Best: %+v",
				q.Len(), len(keys)-len(best.keys), best)
		default:
		}
		if len(s.keys) > len(best.keys) {
			best = s
		}

		// Try to move each bot (since they can only move one at a time).
		for i, cur := range s.bots {
			for nextKey, dist := range distances(cur, s.keys) {
				bots := s.bots
				bots[i] = nextKey
				next := state{
					bots:  bots,
					steps: s.steps + dist,
					keys:  plus(s.keys, nextKey),
				}
				sort.Strings(next.keys)

				if !visited[visKey{next.bots, strings.Join(s.keys, "")}] {
					q.Push(next, next.steps)
				}

				if len(next.keys) == len(keys) {
					t.Logf("Done: %+v", next)
					return next.steps
				}
			}
		}
	}

	return -1
}

func TestPart2(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part2 example 0", `#######
#a.#Cd#
##...##
##.@.##
##...##
#cB#Ab#
#######`, 8},
		{"part2 example 1", `###############
#d.ABC.#.....a#
######...######
######.@.######
######...######
#b.....#.....c#
###############`, 24},
		{"part2 example 2", `#############
#DcBa.#.GhKl#
#.###...#I###
#e#d##@##j#k#
###C#...###J#
#fEbA.#.FgHi#
#############`, 32},
		{"part2 example 3", `#############
#g#f.D#..h#l#
#F###e#E###.#
#dCba...BcIJ#
######@######
#nK.L...G...#
#M###N#H###.#
#o#m..#i#jk.#
#############`, 72},
		{"part2 answer", advent.ReadFile(t, "input.txt"), 2066},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(t, test.in), test.want; got != want {
				t.Errorf("part2(%#v)\n = %#v, want %#v", "test.in", got, want)
			}
		})
	}
}

func plus(slice []string, s string) []string {
	return append(slice[:len(slice):len(slice)], s)
}

func plusCoord(slice []coords.Coord, s coords.Coord) []coords.Coord {
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

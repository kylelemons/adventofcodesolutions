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

func part1(t *testing.T, in string) (ret int) {
	input := parseInput(t, in)

	type retCode int
	const (
		hitWall  retCode = 0
		success  retCode = 1
		atTarget retCode = 2
	)
	tryDir := func(p *intcode.Program, dir coords.Vector) (p2 *intcode.Program, ret retCode) {
		prog := p.Snapshot()
		prog.Input = func() int {
			switch dir {
			case coords.North:
				return 1
			case coords.South:
				return 2
			case coords.West:
				return 3
			case coords.East:
				return 4
			default:
				panic(dir)
			}
		}
		prog.Output = func(v int) {
			ret = retCode(v)
			prog.Halt()
		}
		prog.Run(t)
		return prog, ret
	}
	type state struct {
		loc   coords.Coord
		prog  *intcode.Program
		steps int
	}

	q := []state{{loc: coords.RC(0, 0), prog: input.Base}}
	visited := make(map[coords.Coord]byte)
	var rows, cols advent.RangeTracker

	sprint := func() string {
		if !rows.Valid {
			return "[]"
		}
		out := new(strings.Builder)
		for r := rows.Min; r <= rows.Max; r++ {
			for c := cols.Min; c <= cols.Max; c++ {
				ch := byte(' ')
				if v, ok := visited[coords.RC(r, c)]; ok {
					ch = v
				}
				out.WriteByte(ch)
			}
			out.WriteByte('\n')
		}
		return out.String()
	}

	for len(q) > 0 {
		s := q[0]
		q = q[1:]

		// log.Printf("Current map:\n%s", sprint())

		for _, next := range coords.Cardinals {
			nextLoc := s.loc.Add(next)
			if _, ok := visited[nextLoc]; ok {
				continue
			}
			rows.Track(nextLoc.R())
			cols.Track(nextLoc.C())

			p2, ret := tryDir(s.prog, next)
			switch ret {
			case hitWall:
				visited[nextLoc] = '.'
				continue
			case success:
				visited[nextLoc] = ' '
			case atTarget:
				visited[nextLoc] = 'O'
				t.Logf("Final map:\n%s", sprint())
				return s.steps + 1
			}

			q = append(q, state{
				loc:   nextLoc,
				prog:  p2,
				steps: s.steps + 1,
			})
		}
	}
	t.Fatalf("Ran out of areas to visit:\n%s", sprint())
	return
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		// {"part1 example 0", "...", 0},
		{"part1 answer", advent.ReadFile(t, "input.txt"), 222},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part1(t, test.in), test.want; got != want {
				t.Errorf("part1(...)\n = %#v, want %#v", got, want)
			}
		})
	}
}

func part2(t *testing.T, in string) (ret int) {
	input := parseInput(t, in)

	type retCode int
	const (
		hitWall  retCode = 0
		success  retCode = 1
		atTarget retCode = 2
	)
	tryDir := func(p *intcode.Program, dir coords.Vector) (p2 *intcode.Program, ret retCode) {
		prog := p.Snapshot()
		prog.Input = func() int {
			switch dir {
			case coords.North:
				return 1
			case coords.South:
				return 2
			case coords.West:
				return 3
			case coords.East:
				return 4
			default:
				panic(dir)
			}
		}
		prog.Output = func(v int) {
			ret = retCode(v)
			prog.Halt()
		}
		prog.Run(t)
		return prog, ret
	}
	type discoverState struct {
		loc   coords.Coord
		prog  *intcode.Program
		steps int
	}

	dq := []discoverState{{loc: coords.RC(0, 0), prog: input.Base}}
	visited := make(map[coords.Coord]byte)
	var rows, cols advent.RangeTracker

	sprint := func() string {
		if !rows.Valid {
			return "[]"
		}
		out := new(strings.Builder)
		for r := rows.Min; r <= rows.Max; r++ {
			for c := cols.Min; c <= cols.Max; c++ {
				ch := byte(' ')
				if v, ok := visited[coords.RC(r, c)]; ok {
					ch = v
				}
				out.WriteByte(ch)
			}
			out.WriteByte('\n')
		}
		return out.String()
	}

	var dest coords.Coord
	for len(dq) > 0 {
		s := dq[0]
		dq = dq[1:]

		// log.Printf("Current map:\n%s", sprint())

		for _, next := range coords.Cardinals {
			nextLoc := s.loc.Add(next)
			if _, ok := visited[nextLoc]; ok {
				continue
			}
			rows.Track(nextLoc.R())
			cols.Track(nextLoc.C())

			p2, ret := tryDir(s.prog, next)
			switch ret {
			case hitWall:
				visited[nextLoc] = '.'
				continue
			case success:
				visited[nextLoc] = ' '
			case atTarget:
				visited[nextLoc] = 'O'
				dest = nextLoc
			}

			dq = append(dq, discoverState{
				loc:   nextLoc,
				prog:  p2,
				steps: s.steps + 1,
			})
		}
	}
	t.Logf("Entire map: (dest = %v)\n%s", dest, sprint())

	type fillState struct {
		t   int
		loc coords.Coord
	}

	fq := []fillState{{t: 0, loc: dest}}
	visited[dest] = ' '

	for len(fq) > 0 {
		s := fq[0]
		fq = fq[1:]

		if visited[s.loc] != ' ' {
			continue
		}
		// Fill the current location
		visited[s.loc] = 'O'

		// Track the last timestep that filled a space
		ret = s.t // will be 0 when filling the original spot

		for _, next := range coords.Cardinals {
			fq = append(fq, fillState{
				t:   s.t + 1,
				loc: s.loc.Add(next),
			})
		}
	}
	return
}

func TestPart2(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		// {"part2 example 0", "...", 0},
		{"part2 answer", advent.ReadFile(t, "input.txt"), 394},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(t, test.in), test.want; got != want {
				t.Errorf("part2(...)\n = %#v, want %#v", got, want)
			}
		})
	}
}

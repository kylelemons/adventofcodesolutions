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
	"strings"
	"testing"

	"github.com/kylelemons/adventofcodesolutions/2019/intcode"
	"github.com/kylelemons/adventofcodesolutions/advent"
	"github.com/kylelemons/adventofcodesolutions/advent/coords"
)

func part1(t *testing.T, in string) (ret int) {
	prog := intcode.Compile(t, in)

	out := new(strings.Builder)
	prog.Output = func(v int) {
		out.WriteByte(byte(v))
	}
	prog.Run(t)
	t.Logf("Output:\n%s", out)

	field := advent.Split2D(strings.TrimSpace(out.String()))

	var total int
	for r := range field {
		if r <= 0 || r >= len(field)-1 {
			continue
		}
		for c := range field[r] {
			if c <= 0 || c >= len(field[r])-1 {
				continue
			}
			pos := coords.RC(r, c)
			if field[pos.R()][pos.C()] != '#' {
				continue
			}

			var count int
			for _, delta := range coords.Cardinals {
				next := pos.Add(delta)
				if field[next.R()][next.C()] == '#' {
					count++
				}
			}
			if count == 4 {
				total += pos.X() * pos.Y()
			}
		}
	}

	return total
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		// {"part1 example 0", "...", 0},
		{"part1 answer", advent.ReadFile(t, "input.txt"), 3936},
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
	var field [][]byte
	startProg := intcode.Compile(t, in)

	initialState := new(strings.Builder)
	startProg.Output = func(v int) {
		initialState.WriteByte(byte(v))
	}
	startProg.Run(t)

	field = advent.Split2D(strings.TrimSpace(initialState.String()))

	var cur coords.Coord
	var dir coords.Vector
findStart:
	for r := range field {
		for c := range field[r] {
			switch field[r][c] {
			case '^':
				dir = coords.North
			case '<':
				dir = coords.East
			case '>':
				dir = coords.West
			case 'v':
				dir = coords.South
			default:
				continue
			}
			cur = coords.RC(r, c)
			field[r][c] = '#'
			break findStart
		}
	}

	var spaces int
	for r := range field {
		for c := range field[r] {
			if field[r][c] == '#' {
				spaces++
			}
		}
	}

	get := func(next coords.Coord) byte {
		if next.R() < 0 || next.R() >= len(field) || next.C() < 0 || next.C() >= len(field[0]) {
			return '.'
		}
		return field[next.R()][next.C()]
	}

	var path []string
	for {
		if left := cur.Add(dir.Left()); get(left) == '#' {
			path = append(path, "L")
			dir = dir.Left()
		} else if right := cur.Add(dir.Right()); get(right) == '#' {
			path = append(path, "R")
			dir = dir.Right()
		} else {
			break
		}

		var steps int
		for {
			if forward := cur.Add(dir); get(forward) == '#' {
				steps++
				cur = forward
			} else {
				path = append(path, fmt.Sprint(steps))
				break
			}
		}
	}

	long := strings.Join(path, ",")
	t.Logf("Start: %s", long)

	start := 0
	var funcs []string
	for len(funcs) < 3 {
		start += strings.IndexAny(long[start:], "RL0123456789")

		max, at := 0, 0
		for length := 1; start+length < len(long) && length < 20; length++ {
			sub := long[start:][:length]
			if sub[len(sub)-1] == ',' {
				continue
			}
			count := strings.Count(long, sub)
			reduction := count * len(sub)
			if reduction > max {
				max, at = reduction, length
			}
		}
		fname, best := string(rune('A'+len(funcs))), long[start:][:at]
		long = strings.ReplaceAll(long, best, fname)
		t.Logf("%s: %s (reduces by %d)", fname, best, max)
		funcs = append(funcs, best)
		start += 2
	}
	t.Logf("Main: %s", long)

	pending := strings.Join([]string{
		long,
		funcs[0],
		funcs[1],
		funcs[2],
		"n",
	}, "\n") + "\n"

	prog := intcode.Compile(t, in)
	prog.Memory[0] = 2

	prog.Input = func() (v int) {
		c := pending[0]
		pending = pending[1:]
		return int(c)
	}

	line := ""
	prog.Output = func(v int) {
		ret = v
		if v == 10 {
			t.Log(line)
			line = ""
			return
		}
		line += string(byte(v))
	}
	prog.Run(t)
	return
}

func TestPart2(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part2 answer", advent.ReadFile(t, "input.txt"), 785733},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(t, test.in), test.want; got != want {
				t.Errorf("part2(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

func part2manual(t *testing.T, in string) (ret int) {
	prog := intcode.Compile(t, in)
	prog.Memory[0] = 2

	pending := strings.Join([]string{
		"A,B,A,C,A,B,A,C,B,C",
		"R,4,L,12,L,8,R,4",
		"L,8,R,10,R,10,R,6",
		"R,4,R,10,L,12",
		"n",
	}, "\n") + "\n"

	prog.Input = func() (v int) {
		c := pending[0]
		pending = pending[1:]
		return int(c)
	}

	line := ""
	prog.Output = func(v int) {
		ret = v
		if v == 10 {
			t.Log(line)
			line = ""
			return
		}
		line += string(byte(v))
	}
	prog.Run(t)
	return
}

func TestPart2manual(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part2manual answer", advent.ReadFile(t, "input.txt"), 785733},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2manual(t, test.in), test.want; got != want {
				t.Errorf("part2manual(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

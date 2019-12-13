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
package main

import (
	"testing"

	"github.com/kylelemons/adventofcodesolutions/2019/intcode"
	"github.com/kylelemons/adventofcodesolutions/advent"
	"github.com/kylelemons/adventofcodesolutions/advent/coords"
)

func part1(t *testing.T, in string) (ret int) {
	prog := intcode.Compile(t, in)

	var next []int
	prog.Output = func(v int) {
		next = append(next, v)
		if len(next) == 3 {
			got := next
			next = nil

			rc := coords.XY(got[0], got[1])
			if rc == coords.XY(-1, 0) {
				return
			}
			switch got[2] {
			case 2:
				ret++
			}
		}
	}
	prog.Run(nil)

	return
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part1 answer", advent.ReadFile(t, "input.txt"), 258},
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
	prog := intcode.Compile(t, in)
	prog.Memory[0] = 2

	m := make(map[coords.Coord]byte)

	var ball, paddle coords.Coord
	var r, c advent.RangeTracker

	var round int
	prog.Input = func() (v int) {
		// if round%1000 == 0 {
		// 	for rr := r.Min; rr <= r.Max; rr++ {
		// 		for cc := c.Min; cc <= c.Max; cc++ {
		// 			fmt.Printf("%c", m[coords.RC(rr, cc)])
		// 		}
		// 		fmt.Println()
		// 	}
		// 	fmt.Println()
		// }
		round++
		switch {
		case paddle.C() < ball.C():
			return 1
		case paddle.C() > ball.C():
			return -1
		default:
			return 0
		}
	}

	var next []int
	prog.Output = func(v int) {
		next = append(next, v)
		if len(next) == 3 {
			got := next
			next = nil

			rc := coords.XY(got[0], got[1])
			if rc == coords.XY(-1, 0) {
				ret = got[2]
				return
			}
			r.Track(rc.R())
			c.Track(rc.C())
			switch got[2] {
			case 0:
				m[rc] = ' '
			case 1:
				m[rc] = '|'
			case 2:
				m[rc] = '#'
			case 3:
				m[rc] = '='
				paddle = rc
			case 4:
				m[rc] = 'o'
				ball = rc
			}
		}
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
		{"part2 answer", advent.ReadFile(t, "input.txt"), 12765},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(t, test.in), test.want; got != want {
				t.Errorf("part2(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

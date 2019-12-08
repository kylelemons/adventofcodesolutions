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
	"strings"
	"testing"

	"github.com/kylelemons/adventofcodesolutions/advent"
)

type point struct {
	x, y   int
	dx, dy int
}

func parse(t *testing.T, in string) (ret []point) {
	advent.Lines(in).Extract(t, `position=<\s*(-?\d+),\s*(-?\d+)> velocity=<\s*(-?\d+),\s*(-?\d+)>`, func(x, y, dx, dy int) {
		ret = append(ret, point{x, y, dx, dy})
	})
	return
}

func advance(points []point) (x, y advent.RangeTracker) {
	for i, p := range points {
		points[i].x = x.Track(p.x + p.dx)
		points[i].y = y.Track(p.y + p.dy)
		// fmt.Println(points[i], x, y)
	}
	return
}

func backtrack(points []point) (x, y advent.RangeTracker) {
	for i, p := range points {
		points[i].x = x.Track(p.x - p.dx)
		points[i].y = y.Track(p.y - p.dy)
		// fmt.Println(points[i], x, y)
	}
	return
}

func part1(t *testing.T, in string) (ret int) {
	points := parse(t, in)

	lastX, lastY := advance(points)
	ret++
	for {
		curX, curY := advance(points)
		ret++
		if curX.Delta() > lastX.Delta() && curY.Delta() > lastY.Delta() {
			break
		}
		lastX, lastY = curX, curY
	}

	x, y := backtrack(points)
	ret--
	disp := advent.Make2D(y.Delta()+1, x.Delta()+1)
	for _, p := range points {
		disp[p.y-y.Min][p.x-x.Min] = '#'
	}
	for _, row := range disp {
		t.Logf("%s", strings.ReplaceAll(string(row), "\x00", " "))
	}

	return
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part1 example 0", strings.TrimSpace(`
position=< 9,  1> velocity=< 0,  2>
position=< 7,  0> velocity=<-1,  0>
position=< 3, -2> velocity=<-1,  1>
position=< 6, 10> velocity=<-2, -1>
position=< 2, -4> velocity=< 2,  2>
position=<-6, 10> velocity=< 2, -2>
position=< 1,  8> velocity=< 1, -1>
position=< 1,  7> velocity=< 1,  0>
position=<-3, 11> velocity=< 1, -2>
position=< 7,  6> velocity=<-1, -1>
position=<-2,  3> velocity=< 1,  0>
position=<-4,  3> velocity=< 2,  0>
position=<10, -3> velocity=<-1,  1>
position=< 5, 11> velocity=< 1, -2>
position=< 4,  7> velocity=< 0, -1>
position=< 8, -2> velocity=< 0,  1>
position=<15,  0> velocity=<-2,  0>
position=< 1,  6> velocity=< 1,  0>
position=< 8,  9> velocity=< 0, -1>
position=< 3,  3> velocity=<-1,  1>
position=< 0,  5> velocity=< 0, -1>
position=<-2,  2> velocity=< 2,  0>
position=< 5, -2> velocity=< 1,  2>
position=< 1,  4> velocity=< 2,  1>
position=<-2,  7> velocity=< 2, -2>
position=< 3,  6> velocity=<-1, -1>
position=< 5,  0> velocity=< 1,  0>
position=<-6,  0> velocity=< 2,  0>
position=< 5,  9> velocity=< 1, -2>
position=<14,  7> velocity=<-2,  0>
position=<-3,  6> velocity=< 2, -1>
`), 3},
		{"part1 answer", advent.ReadFile(t, "input.txt"), 10905},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part1(t, test.in), test.want; got != want {
				t.Errorf("part1(%#v)\n part2 = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

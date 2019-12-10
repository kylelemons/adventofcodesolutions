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
	"math"
	"math/big"
	"sort"
	"strings"
	"testing"

	"github.com/kylelemons/adventofcodesolutions/advent"
)

type Point struct{ X, Y int }

func (p Point) String() string { return fmt.Sprintf("%v,%v", p.X, p.Y) }

type Input struct {
	Grid      [][]byte
	Asteroids []Point
}

func parseInput(t *testing.T, in string) *Input {
	input := &Input{
		Grid: advent.Split2D(in),
	}
	for r := range input.Grid {
		for c := range input.Grid[r] {
			if input.Grid[r][c] == '#' {
				input.Asteroids = append(input.Asteroids, Point{c, r})
			}
		}
	}
	return input
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func part1(t *testing.T, in string) (ret int, base Point) {
	input := parseInput(t, in)

	var max int
	var at Point
	for _, from := range input.Asteroids {
		type Key struct {
			slope      string
			negY, negX bool
		}
		canSee := map[Key]bool{}
		for _, to := range input.Asteroids {
			if to == from {
				continue
			}

			dx := to.X - from.X
			dy := to.Y - from.Y

			var k Key
			switch {
			case to.Y == from.Y:
				k = Key{negX: dx < 0, slope: "horiz"}
			case to.X == from.X:
				k = Key{negY: dy < 0, slope: "vert"}
			default:
				slope := big.NewRat(int64(abs(dy)), int64(abs(dx))).RatString()
				k = Key{negY: dy < 0, negX: dx < 0, slope: slope}
			}
			if canSee[k] {
				continue
			}
			canSee[k] = true

		}
		if m := len(canSee); m > max {
			max, at = m, from
		}
	}

	return max, at
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name  string
		in    string
		want1 int
		want2 Point
	}{
		{"part1 example 0", `
......#.#.
#..#.#....
..#######.
.#.#.###..
.#..#.....
..#....#.#
#..#....#.
.##.#..###
##...#..#.
.#....####
`, 33, Point{5, 8}},
		{"part1 example 1", `
#.#...#.#.
.###....#.
.#....#...
##.#.#.#.#
....#.#.#.
.##..###.#
..#...##..
..##....##
......#...
.####.###.
`, 35, Point{1, 2}},
		{"part1 example 2", `
.#..#..###
####.###.#
....###.#.
..###.##.#
##.##.#.#.
....###..#
..#.#..#.#
#..#.#.###
.##...##.#
.....#.#..
`, 41, Point{6, 3}},
		{"part1 answer", advent.ReadFile(t, "input.txt"), 299, Point{X: 26, Y: 29}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got1, got2 := part1(t, strings.Trim(test.in, "\n"))
			if got, want := got1, test.want1; got != want {
				t.Errorf("part1(%#v).1\n = %#v, want %#v", test.in, got, want)
			}
			if got, want := got2, test.want2; got != want {
				t.Errorf("part1(%#v).1\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

func part2(t *testing.T, in string, from Point) (bet Point) {
	input := parseInput(t, in)
	t.Logf("Starting with %v", from)

	type Asteroid struct {
		Location Point
		Angle    float64 // measured from the Y axis clockwise
		Distance float64
	}

	type Slope struct {
		Slope      string
		NegY, NegX bool
	}
	sub := func(dy, dx int) Slope {
		switch {
		case dy == 0:
			return Slope{Slope: "horiz", NegX: dx < 0}
		case dx == 0:
			return Slope{Slope: "vert", NegY: dy < 0}
		default:
			slope := big.NewRat(int64(abs(dy)), int64(abs(dx))).RatString()
			return Slope{Slope: slope, NegY: dy < 0, NegX: dx < 0}
		}
	}
	ycw := func(dy, dx int) float64 {
		dy = -dy

		// Determine the ccw angle from +x
		angle := math.Atan2(float64(dy), float64(dx))
		// orig := angle
		if angle < 0 {
			angle += 2 * math.Pi
		}

		// Determine the cw angle from +x
		angle = 2*math.Pi - angle

		// Determine the cw angle from +y
		angle += math.Pi / 2
		if angle >= 2*math.Pi {
			angle -= 2 * math.Pi
		}
		// t.Logf("[x:%d,y:%d] / angle = %v / ycw = %v", dx, dy, orig*180/math.Pi, angle*180/math.Pi)
		return angle
	}
	dist := func(dy, dx int) float64 {
		return math.Sqrt(float64(dy*dy) + float64(dx*dx))
	}

	first := make(map[Slope]Asteroid)
	groups := make(map[float64][]Asteroid)
	for _, target := range input.Asteroids {
		if target == from {
			continue
		}

		dx := target.X - from.X
		dy := target.Y - from.Y
		slope := sub(dy, dx)
		angle := ycw(dy, dx)
		dist := dist(dy, dx)

		// t.Logf("%v: Angle is %v / %v", target, angle*180/math.Pi, angle)

		if prev, ok := first[slope]; ok {
			angle = prev.Angle
		} else {
			first[slope] = Asteroid{Location: target, Angle: angle, Distance: dist}
		}
		groups[angle] = append(groups[angle], Asteroid{Location: target, Angle: angle, Distance: dist})
	}

	type OrderedGroup struct {
		Angle     float64
		Asteroids []Asteroid
	}
	var ordered []*OrderedGroup
	for angle, asteroids := range groups {
		ordered = append(ordered, &OrderedGroup{angle, asteroids})
	}
	sort.Slice(ordered, func(i, j int) bool {
		return ordered[i].Angle < ordered[j].Angle
	})

	// t.Logf("Groups: %d, Asteroids: %d", len(groups), len(input.Asteroids))
	for _, og := range ordered {
		sort.Slice(og.Asteroids, func(i, j int) bool {
			return og.Asteroids[i].Distance < og.Asteroids[j].Distance
		})
		// t.Logf(" - %+v", og)
	}

	index := 1
	for {
		var vaporized bool
		for _, og := range ordered {
			if len(og.Asteroids) == 0 {
				continue
			}

			destroy := og.Asteroids[0]
			og.Asteroids = og.Asteroids[1:]
			t.Logf("%3d: Destroying %+v  \t(%6.2fdeg)", index, destroy.Location, destroy.Angle*180/math.Pi)
			if index == 200 {
				return destroy.Location
			}
			index++
			vaporized = true
		}
		if !vaporized {
			return Point{-1, -1}
		}
	}
}

func TestPart2(t *testing.T) {
	tests := []struct {
		name  string
		in    string
		start Point
		want  Point
	}{
		{"part2 example 1", `
.#....#####...#..
##...##.#####..##
##...#...#.#####.
..#.....X...###..
..#.#.....#....##
`, Point{8, 3}, Point{-1, -1}},
		{"part2 example 2", `
.#..##.###...#######
##.############..##.
.#.######.########.#
.###.#######.####.#.
#####.##.#.##.###.##
..#####..#.#########
####################
#.####....###.#.#.##
##.#################
#####.##.###..####..
..######..##.#######
####.##.####...##..#
.#####..#.######.###
##...#.##########...
#.##########.#######
.####.#.###.###.#.##
....##.##.###..#####
.#.#.###########.###
#.#.#.#####.####.###
###.##.####.##.#..##
`, Point{11, 13}, Point{8, 2}},
		{"part2 answer", advent.ReadFile(t, "input.txt"), Point{X: 26, Y: 29}, Point{14, 19}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := part2(t, strings.Trim(test.in, "\n"), test.start)
			if got, want := got, test.want; got != want {
				t.Errorf("part2(%#v)\n = %v, want %v", test.in, got, want)
			}
			t.Log(got.X*100 + got.Y)
		})
	}
}

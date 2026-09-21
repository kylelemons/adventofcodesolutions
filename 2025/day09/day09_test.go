// Copyright 2021 Kyle Lemons
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
package aocday

import (
	"iter"
	"testing"

	"github.com/kylelemons/adventofcodesolutions/advent"
)

type Tile struct{ R, C int }
type Input struct {
	Tiles []Tile
}

func parseInput(t *testing.T, in string) *Input {
	input := &Input{
		// ...
	}
	advent.Lines(in).Extract(t, `(\d+),(\d+)`, func(r, c int) {
		input.Tiles = append(input.Tiles, Tile{R: r, C: c})
	})
	return input
}

func part1(t *testing.T, in string) (ret int) {
	input := parseInput(t, in)

	for i, t0 := range input.Tiles {
		for _, t1 := range input.Tiles[i+1:] {
			minR := min(t0.R, t1.R)
			minC := min(t0.C, t1.C)
			maxR := max(t0.R, t1.R)
			maxC := max(t0.C, t1.C)
			area := (maxR - minR + 1) * (maxC - minC + 1)
			if area > ret {
				ret = area
			}
		}
	}

	return
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part1 example 0", "7,1\n11,1\n11,7\n9,7\n9,5\n2,5\n2,3\n7,3", 50},
		{"part1 answer", advent.ReadFile(t, "input.txt"), 4781235324},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part1(t, test.in), test.want; got != want {
				t.Errorf("part1(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

// loopPairs yields consecutive (e0, e1) pairs of vertices representing each edge,
// wrapping around to close the loop between the last and first vertex.
func loopPairs(tiles []Tile) iter.Seq2[Tile, Tile] {
	return func(yield func(Tile, Tile) bool) {
		for i := 0; i+1 < len(tiles); i++ {
			if !yield(tiles[i], tiles[i+1]) {
				return
			}
		}
		yield(tiles[len(tiles)-1], tiles[0]) // close the loop
	}
}

// edgePiercesRect returns true if the edge (e0, e1) cuts into the strict interior
// of the bounding box [minR, maxR] x [minC, maxC].
//
// Edges that run along the border of the rectangle (or only touch corners) are
// valid boundary lines and do not pierce the strict interior.
func edgePiercesRect(minR, maxR, minC, maxC int, e0, e1 Tile) bool {
	// If the edge is vertical (e0.C == e1.C), transpose (R, C) with (C, R) so the
	// same horizontal collision check can be reused for both orientations.
	if e0.C == e1.C {
		minR, maxR, minC, maxC = minC, maxC, minR, maxR
		e0.R, e0.C, e1.R, e1.C = e0.C, e0.R, e1.C, e1.R
	}
	r := e0.R
	cLo, cHi := min(e0.C, e1.C), max(e0.C, e1.C)
	// Strict inequalities allow edges lying on rectangle borders (e.g. r == minR)
	// while catching any edge passing through the open interval (minR, maxR) x (minC, maxC).
	return minR < r && r < maxR && cHi > minC && cLo < maxC
}

// pointInPolygon determines whether continuous coordinates (r, c) lie inside
// the orthogonal polygon using ray casting in the +C (rightward) direction.
//
// Since the ray travels horizontally along row r, only perpendicular (vertical)
// edges can intersect it. We count crossings with vertical edges located to the
// right of the test point.
func pointInPolygon(r, c float64, tiles []Tile) bool {
	inside := false
	for e0, e1 := range loopPairs(tiles) {
		if e0.C == e1.C {
			edgeC := float64(e0.C)
			rLo := float64(min(e0.R, e1.R))
			rHi := float64(max(e0.R, e1.R))
			if edgeC > c && rLo < r && r < rHi {
				inside = !inside
			}
		}
	}
	return inside
}

func part2(t *testing.T, in string) (ret int) {
	input := parseInput(t, in)

	for i, t0 := range input.Tiles {
	nextTile:
		for _, t1 := range input.Tiles[i+1:] {
			minR, maxR := min(t0.R, t1.R), max(t0.R, t1.R)
			minC, maxC := min(t0.C, t1.C), max(t0.C, t1.C)

			// 1. Ensure no polygon boundary edge enters the interior of the rectangle.
			// Because rectangles are convex, any exterior void or cavity inside the bounding
			// box must be bounded by edges that cut into this region.
			for e0, e1 := range loopPairs(input.Tiles) {
				if edgePiercesRect(minR, maxR, minC, maxC, e0, e1) {
					continue nextTile
				}
			}

			// 2. Ensure the rectangle is inside the polygon (and not in an exterior void).
			// We sample the continuous midpoint. A tiny epsilon offset ensures non-integer
			// ray coordinates, eliminating collinear ray / vertex intersection singularities.
			const eps = 1e-4
			midR := float64(minR+maxR)/2.0 + eps
			midC := float64(minC+maxC)/2.0 + eps
			if pointInPolygon(midR, midC, input.Tiles) {
				area := (maxR - minR + 1) * (maxC - minC + 1)
				if area > ret {
					ret = area
				}
			}
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
		{"part2 example 0", "7,1\n11,1\n11,7\n9,7\n9,5\n2,5\n2,3\n7,3", 24},
		{"part2 answer", advent.ReadFile(t, "input.txt"), 1566935900},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(t, test.in), test.want; got != want {
				t.Errorf("part2(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

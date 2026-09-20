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
	"bufio"
	"cmp"
	// "fmt"
	"math"
	"slices"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/kylelemons/adventofcodesolutions/advent"
)

type Loc [3]int

func (a Loc) Dist(b Loc) float64 {
	return math.Sqrt(float64(
		(a[0]-b[0])*(a[0]-b[0]) +
			(a[1]-b[1])*(a[1]-b[1]) +
			(a[2]-b[2])*(a[2]-b[2]),
	))
}

type Input struct {
	Boxes []Loc
}

func parseInput(t *testing.T, in string) *Input {
	input := &Input{
		// ...
	}
	lines := bufio.NewScanner(strings.NewReader(strings.TrimSpace(in)))
	for lines.Scan() {
		line := lines.Text()
		xyz := strings.Split(line, ",")
		if len(xyz) != 3 {
			t.Fatalf("invalid input line %q", line)
		}
		x, _ := strconv.Atoi(xyz[0])
		y, _ := strconv.Atoi(xyz[1])
		z, _ := strconv.Atoi(xyz[2])
		input.Boxes = append(input.Boxes, Loc{x, y, z})
	}
	return input
}

func part1(t *testing.T, in string, n int) (ret int) {
	input := parseInput(t, in)

	type Dist struct {
		A, B Loc
		Dist float64
	}
	var dists []Dist
	for i, a := range input.Boxes {
		for _, b := range input.Boxes[i+1:] {
			if a == b {
				continue
			}
			dists = append(dists, Dist{a, b, a.Dist(b)})
		}
	}

	slices.SortFunc(dists, func(a, b Dist) int {
		return cmp.Or(
			cmp.Compare(a.Dist, b.Dist),
			cmp.Compare(a.A[0], b.A[0]),
			cmp.Compare(a.A[1], b.A[1]),
			cmp.Compare(a.A[2], b.A[2]),
			cmp.Compare(a.B[0], b.B[0]),
			cmp.Compare(a.B[1], b.B[1]),
			cmp.Compare(a.B[2], b.B[2]),
		)
	})

	// disjoint set forest
	rootOf := make(map[Loc]Loc)
	// everything starts in its own set
	for _, box := range input.Boxes {
		rootOf[box] = box
	}
	// joining two forests
	join := func(a, b Loc) {
		rootOf[b] = a
		for box, root := range rootOf {
			if root == b {
				rootOf[box] = a
			}
		}
	}

	for i := range dists[:n] {
		if n <= 0 {
			break
		}
		// fmt.Printf("Closest: %v - %v (dist: %v)\n", dists[i].A, dists[i].B, dists[i].Dist)
		a, b := rootOf[dists[i].A], rootOf[dists[i].B]
		if a == b {
			// fmt.Printf(" ... Same circuit!\n")
			continue
		}
		join(a, b)
	}

	sizes := make(map[Loc]int)
	for _, root := range rootOf {
		sizes[root]++
	}

	type netsize struct {
		Root Loc
		Size int
	}
	var nets []netsize
	for root, size := range sizes {
		nets = append(nets, netsize{root, size})
	}
	sort.Slice(nets, func(i, j int) bool {
		return nets[i].Size > nets[j].Size
	})
	// for i := range nets {
	// 	fmt.Printf("Net %v: %v\n", i, nets[i])
	// }

	return nets[0].Size * nets[1].Size * nets[2].Size
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name string
		in   string
		n    int
		want int
	}{
		{"part1 example 0", "162,817,812\n57,618,57\n906,360,560\n592,479,940\n352,342,300\n466,668,158\n542,29,236\n431,825,988\n739,650,466\n52,470,668\n216,146,977\n819,987,18\n117,168,530\n805,96,715\n346,949,466\n970,615,88\n941,993,340\n862,61,35\n984,92,344\n425,690,689", 10, 40},
		{"part1 answer", advent.ReadFile(t, "input.txt"), 1000, 46398},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part1(t, test.in, test.n), test.want; got != want {
				t.Errorf("part1(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

func part2(t *testing.T, in string, n int) (ret int) {
	input := parseInput(t, in)

	type Dist struct {
		A, B Loc
		Dist float64
	}
	var dists []Dist
	for i, a := range input.Boxes {
		for _, b := range input.Boxes[i+1:] {
			if a == b {
				continue
			}
			dists = append(dists, Dist{a, b, a.Dist(b)})
		}
	}

	slices.SortFunc(dists, func(a, b Dist) int {
		return cmp.Or(
			cmp.Compare(a.Dist, b.Dist),
			cmp.Compare(a.A[0], b.A[0]),
			cmp.Compare(a.A[1], b.A[1]),
			cmp.Compare(a.A[2], b.A[2]),
			cmp.Compare(a.B[0], b.B[0]),
			cmp.Compare(a.B[1], b.B[1]),
			cmp.Compare(a.B[2], b.B[2]),
		)
	})

	// disjoint set forest
	roots := make(map[Loc]struct{})
	rootOf := make(map[Loc]Loc)
	// everything starts in its own set
	for _, box := range input.Boxes {
		roots[box] = struct{}{}
		rootOf[box] = box
	}
	// joining two forests
	join := func(a, b Loc) (joined bool) {
		a = rootOf[a]
		b = rootOf[b]
		if a == b {
			return false
		}

		rootOf[b] = a
		delete(roots, b)

		for box, root := range rootOf {
			if root == b {
				rootOf[box] = a
			}
		}
		return true
	}

	for i := range dists {
		if n <= 0 {
			break
		}
		a, b := dists[i].A, dists[i].B
		join(a, b)
		if len(roots) == 1 {
			return a[0] * b[0]
		}
	}
	panic("unreachable")
}

func TestPart2(t *testing.T) {
	tests := []struct {
		name string
		in   string
		n    int
		want int
	}{
		{"part2 example 0", "162,817,812\n57,618,57\n906,360,560\n592,479,940\n352,342,300\n466,668,158\n542,29,236\n431,825,988\n739,650,466\n52,470,668\n216,146,977\n819,987,18\n117,168,530\n805,96,715\n346,949,466\n970,615,88\n941,993,340\n862,61,35\n984,92,344\n425,690,689", 10, 25272},
		{"part2 answer", advent.ReadFile(t, "input.txt"), 1000, 8141888143},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(t, test.in, test.n), test.want; got != want {
				t.Errorf("part2(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

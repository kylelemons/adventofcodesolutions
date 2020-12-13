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
	"strings"
	"testing"
	"text/tabwriter"

	"github.com/kylelemons/adventofcodesolutions/advent"
)

const Axes = 3

type Planet struct {
	Pos [Axes]int
	Vel [Axes]int
}

type Input struct {
	Io, Europa, Ganymede, Callisto Planet

	Planets []*Planet
}

func parseInput(t testing.TB, in string) *Input {
	input := &Input{}
	input.Planets = []*Planet{&input.Io, &input.Europa, &input.Ganymede, &input.Callisto}

	var i int
	advent.Lines(in).Extract(t, `<x=(-?\d+), y=(-?\d+), z=(-?\d+)>`, func(x, y, z int) {
		input.Planets[i].Pos = [3]int{x, y, z}
		i++
	})
	return input
}

func (in *Input) StepAxis(axis int) {
	// Apply gravity
	for i, p0 := range in.Planets {
		for j, p1 := range in.Planets {
			if i < j {
				continue
			}
			switch {
			case p0.Pos[axis] < p1.Pos[axis]:
				p0.Vel[axis]++
				p1.Vel[axis]--
			case p0.Pos[axis] > p1.Pos[axis]:
				p0.Vel[axis]--
				p1.Vel[axis]++
			}
		}
	}
	// Apply velocity
	for _, p := range in.Planets {
		p.Pos[axis] += p.Vel[axis]
	}
}

func (in *Input) Axis(axis int) [6]int {
	return [6]int{
		in.Planets[0].Pos[axis],
		in.Planets[1].Pos[axis],
		in.Planets[2].Pos[axis],
		in.Planets[0].Vel[axis],
		in.Planets[1].Vel[axis],
		in.Planets[2].Vel[axis],
	}
}

func (in *Input) TimeStep() {
	for i := 0; i < Axes; i++ {
		in.StepAxis(i)
	}
}

func (in *Input) Energy() (energy int) {
	abs := func(a int) int {
		if a < 0 {
			return -a
		}
		return a
	}
	for _, p := range in.Planets {
		var kinetic, potential int
		for _, v := range p.Vel {
			kinetic += abs(v)
		}
		for _, d := range p.Pos {
			potential += abs(d)
		}
		energy += kinetic * potential
	}
	return
}

func (in *Input) String() string {
	buf := new(strings.Builder)
	w := tabwriter.NewWriter(buf, 0, 0, 0, ' ', tabwriter.AlignRight)
	for i, p := range in.Planets {
		fmt.Fprintf(w, "%d: \tpos=<x=\t%d\t,y=\t%d\t,z=\t%d\t>, vel=<x=\t%d\t,y=\t%d\t,z=\t%d\t>\n",
			i, p.Pos[0], p.Pos[1], p.Pos[2], p.Vel[0], p.Vel[1], p.Vel[2])
	}
	w.Flush()
	return buf.String()
}

func part1(t *testing.T, in string, N int) (ret int) {
	input := parseInput(t, in)

	for i := 0; i < N; i++ {
		if i%100 == 0 { //}= 100 {
			t.Logf("After %d steps: (e=%d)\n%s", i, input.Energy(), input)
		}
		input.TimeStep()
	}
	e := input.Energy()
	t.Logf("After %d steps: (e=%d)\n%s", N, e, input)

	return e
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name  string
		in    string
		steps int
		want  int
	}{
		{"part1 example 0", `
<x=-1, y=0, z=2>
<x=2, y=-10, z=-7>
<x=4, y=-8, z=8>
<x=3, y=5, z=-1>`, 10, 179},
		{"part1 example 1", `
<x=-8, y=-10, z=0>
<x=5, y=5, z=10>
<x=2, y=-7, z=3>
<x=9, y=-8, z=-3>`, 100, 1940},
		{"part1 answer", advent.ReadFile(t, "input.txt"), 1000, 9441},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part1(t, strings.TrimLeft(test.in, "\n"), test.steps), test.want; got != want {
				t.Errorf("part1(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

func part2(t *testing.T, in string) (ret int) {
	input := parseInput(t, in)

	var periods []int
nextAxis:
	for axis := 0; axis < Axes; axis++ {
		start := input.Axis(axis)

		for steps := 1; ; steps++ {
			input.StepAxis(axis)
			if input.Axis(axis) == start {
				periods = append(periods, steps)
				continue nextAxis
			}
		}
	}
	return advent.LCM(periods...)
}

func TestPart2(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part2 example 0", `
<x=-1, y=0, z=2>
<x=2, y=-10, z=-7>
<x=4, y=-8, z=8>
<x=3, y=5, z=-1>`, 2772},
		{"part2 example 1", `
<x=-8, y=-10, z=0>
<x=5, y=5, z=10>
<x=2, y=-7, z=3>
<x=9, y=-8, z=-3>`, 4686774924},
		{"part2 answer", advent.ReadFile(t, "input.txt"), 503560201099704},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(t, strings.TrimLeft(test.in, "\n")), test.want; got != want {
				t.Errorf("part2(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

func BenchmarkTimeStep(b *testing.B) {
	input := parseInput(b, advent.ReadFile(b, "input.txt"))

	for i := 0; i < b.N; i++ {
		input.TimeStep()
	}
}

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

	"github.com/google/go-cmp/cmp"
	"github.com/kylelemons/adventofcodesolutions/advent"
)

type part1in struct {
	width, height int
	input         string
}

type part1ret struct {
	FewestZeros int
	FZLayer     int
	Answer      int
}

func part1(t *testing.T, in part1in) (ret part1ret) {
	var layers [][]string
	rem := in.input
	for len(rem) > 0 {
		image := make([]string, in.height)
		zeros, ones, twos := 0, 0, 0
		for i := range image {
			image[i], rem = rem[:in.width], rem[in.width:]
			zeros += strings.Count(image[i], "0")
			ones += strings.Count(image[i], "1")
			twos += strings.Count(image[i], "2")
		}
		layers = append(layers, image)
		if len(layers) == 1 || zeros < ret.FewestZeros {
			ret.FewestZeros = zeros
			ret.FZLayer = len(layers) - 1
			ret.Answer = ones * twos
			t.Logf("New best:")
			for i := range image {
				t.Logf("%q", image[i])
			}
		}
	}
	t.Logf("%d layers", len(layers))
	return
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name string
		in   part1in
		want part1ret
	}{
		{"part1 answer", part1in{25, 6, advent.ReadFile(t, "input.txt")}, part1ret{
			FewestZeros: 8,
			FZLayer:     5,
			Answer:      2016,
		}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part1(t, test.in), test.want; got != want {
				t.Errorf("part1(%#v)\n = %#v, want %#v", test.in, got, want)
				t.Errorf("Diff: (-got +want)\n%s", cmp.Diff(got, want))
			}
		})
	}
}

type part2in struct {
	width, height int
	input         string
}

type part2ret struct {
	Image []string
}

func part2(t *testing.T, in part2in) (ret part2ret) {

	var layers [][]string
	rem := in.input
	for len(rem) > 0 {
		image := make([]string, in.height)
		for i := range image {
			image[i], rem = rem[:in.width], rem[in.width:]
		}
		layers = append(layers, image)
	}
	t.Logf("%d layers", len(layers))

	out := advent.Make2D(in.height, in.width)
	for l := len(layers) - 1; l >= 0; l-- {
		layer := layers[l]
		for r := range layer {
			for c, ch := range []byte(layer[r]) {
				switch ch {
				case '0': // black
					out[r][c] = ' '

				case '1': // white
					out[r][c] = '#'
				}
			}
		}
	}

	for i := range out {
		ret.Image = append(ret.Image, string(out[i]))
	}

	return
}

func TestPart2(t *testing.T) {
	tests := []struct {
		name string
		in   part2in
		want part2ret
	}{
		{"part2 answer", part2in{25, 6, advent.ReadFile(t, "input.txt")}, part2ret{
			Image: []string{
				"#  # ####  ##  #### #  # ",
				"#  #    # #  #    # #  # ",
				"####   #  #      #  #  # ",
				"#  #  #   #     #   #  # ",
				"#  # #    #  # #    #  # ",
				"#  # ####  ##  ####  ##  ",
			},
		}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(t, test.in), test.want; !cmp.Equal(got, want) {
				t.Errorf("part2(%#v)\n = %#v, want %#v", test.in, got, want)
				t.Errorf("Diff: (-got +want)\n%s", cmp.Diff(got, want))
			}
		})
	}
}

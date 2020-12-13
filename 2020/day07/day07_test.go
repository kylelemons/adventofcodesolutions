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
package aocday

import (
	"regexp"
	"testing"

	"github.com/kylelemons/adventofcodesolutions/advent"
)

type edge struct {
	count int
	color string
}

type Input struct {
	edges map[string][]edge // edges[color] = list_of_edges
}

func parseInput(t *testing.T, in string) *Input {
	input := &Input{
		edges: make(map[string][]edge),
	}
	advent.Lines(in).Each(func(i int, line advent.Scanner) {
		var color, rules string
		line.Extract(t, `(.*) bags contain (.*)`, &color, &rules)
		input.edges[color] = nil
		if m, _ := regexp.MatchString("no other bags", rules); m {
			return
		}
		advent.Split(rules, ',').Extract(t, `(\d+) (.*) bags?`, func(count int, innerColor string) {
			input.edges[color] = append(input.edges[color], edge{
				count: count,
				color: innerColor,
			})
		})
	})
	return input
}

func part1(t *testing.T, in string) (ret int) {
	input := parseInput(t, in)

	var walk func(color string) (hasGold bool)
	walk = func(color string) (hasGold bool) {
		if color == "shiny gold" {
			return true
		}
		for _, edge := range input.edges[color] {
			if walk(edge.color) {
				return true
			}
		}
		return false
	}
	for color := range input.edges {
		if color == "shiny gold" {
			continue
		}
		if walk(color) {
			ret++
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
		{"part1 example 0", `light red bags contain 1 bright white bag, 2 muted yellow bags.
dark orange bags contain 3 bright white bags, 4 muted yellow bags.
bright white bags contain 1 shiny gold bag.
muted yellow bags contain 2 shiny gold bags, 9 faded blue bags.
shiny gold bags contain 1 dark olive bag, 2 vibrant plum bags.
dark olive bags contain 3 faded blue bags, 4 dotted black bags.
vibrant plum bags contain 5 faded blue bags, 6 dotted black bags.
faded blue bags contain no other bags.
dotted black bags contain no other bags.`, 4},
		{"part1 answer", advent.ReadFile(t, "input.txt"), 259},
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

	var walk func(color string) (contained int)
	walk = func(color string) (contained int) {
		contained++
		for _, edge := range input.edges[color] {
			contained += edge.count * walk(edge.color)
		}
		return
	}
	return walk("shiny gold") - 1 // don't count the gold bag
}

func TestPart2(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part2 example 0", `light red bags contain 1 bright white bag, 2 muted yellow bags.
dark orange bags contain 3 bright white bags, 4 muted yellow bags.
bright white bags contain 1 shiny gold bag.
muted yellow bags contain 2 shiny gold bags, 9 faded blue bags.
shiny gold bags contain 1 dark olive bag, 2 vibrant plum bags.
dark olive bags contain 3 faded blue bags, 4 dotted black bags.
vibrant plum bags contain 5 faded blue bags, 6 dotted black bags.
faded blue bags contain no other bags.
dotted black bags contain no other bags.`, 32},
		{"part2 example 1", `shiny gold bags contain 2 dark red bags.
dark red bags contain 2 dark orange bags.
dark orange bags contain 2 dark yellow bags.
dark yellow bags contain 2 dark green bags.
dark green bags contain 2 dark blue bags.
dark blue bags contain 2 dark violet bags.
dark violet bags contain no other bags.`, 126},
		{"part2 answer", advent.ReadFile(t, "input.txt"), 45018},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(t, test.in), test.want; got != want {
				t.Errorf("part2(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

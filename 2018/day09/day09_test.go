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
	"container/ring"
	"fmt"
	"testing"

	"github.com/kylelemons/adventofcodesolutions/advent"
)

func part1(t *testing.T, in string) (ret int) {
	var players, last int
	advent.Scanner(in).Extract(t, `(\d+) players; last marble is worth (\d+) points`, &players, &last)

	circle := make([]int, 1, last)
	circle[0] = 0

	player := 1
	current := 0
	next := 1

	insert := func(v, at int) (i int) {
		target := at % len(circle)
		if target == 0 {
			circle = append(circle, v)
			return len(circle) - 1
		}

		circle = append(circle, 0)
		copy(circle[target+1:], circle[target:])
		circle[target] = v
		return target
	}
	remove := func(at int) (v, idx int) {
		target := at % len(circle)
		v = circle[target]
		circle = append(circle[:target], circle[target+1:]...)
		if target >= len(circle) {
			return v, 0
		}
		return v, target
	}

	scores := make(map[int]int) // scores[player] = player_score

	for next <= last {
		if next%100000 == 0 {
			fmt.Println(next)
		}
		if next%23 == 0 {
			v, at := remove(current + len(circle) - 7)
			scores[player] += next + v
			current = at
		} else {
			current = insert(next, current+2)
		}

		// t.Log(player, circle[:current], circle[current], circle[current+1:])
		next++
		player++
		if player > players {
			player = 1
		}
	}
	for _, score := range scores {
		if score > ret {
			ret = score
		}
	}

	return
}

type Marble struct {
	Value int
	Next  *Marble
}

func part2(t *testing.T, in string) (ret int) {
	var players, last int
	advent.Scanner(in).Extract(t, `(\d+) players; last marble is worth (\d+) points`, &players, &last)

	circle := &ring.Ring{Value: 0}

	player := 1
	next := 1

	scores := make(map[int]int) // scores[player] = player_score

	for next <= last {
		if next%23 == 0 {
			circle = circle.Move(-8)
			scores[player] += next + circle.Unlink(1).Value.(int)
			circle = circle.Next()
		} else {
			circle = circle.Move(1).Link(&ring.Ring{Value: next}).Prev()
		}

		// if last < 100 {
		// 	start := circle
		// 	current := circle.Next()
		// 	best := circle
		// 	for current != start {
		// 		if current.Value.(int) < best.Value.(int) {
		// 			best = current
		// 		}
		// 		current = current.Next()
		// 	}
		// 	var out []int
		// 	best.Do(func(v interface{}) {
		// 		out = append(out, v.(int))
		// 	})
		// 	t.Log(player, circle.Value, out)
		// }

		next++
		player++
		if player > players {
			player = 1
		}
	}
	for _, score := range scores {
		if score > ret {
			ret = score
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
		{"part1 example 0", "9 players; last marble is worth 25 points", 32},
		{"part1 answer", "459 players; last marble is worth 71790 points", 386151},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part1(t, test.in), test.want; got != want {
				t.Errorf("part1(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

func TestPart2(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"part1 example 0", "9 players; last marble is worth 25 points", 32},
		{"part1 answer", "459 players; last marble is worth 71790 points", 386151},
		{"part2 answer", "459 players; last marble is worth 7179000 points", 3211264152},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(t, test.in), test.want; got != want {
				t.Errorf("part2(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

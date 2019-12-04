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
	"container/heap"
	"sort"
	"strings"
	"testing"

	"github.com/kylelemons/adventofcodesolutions/2019/advent"
)

func parse(t *testing.T, in string) map[string][]string {
	prereq := make(map[string][]string)
	advent.Lines(in).Extract(t, `Step (\S+) must be finished before step (\S+) can begin.`, func(prev, step string) {
		prereq[step] = append(prereq[step], prev)
		if len(prereq[prev]) == 0 {
			prereq[prev] = nil
		}
	})
	return prereq
}

func part1(t *testing.T, in string) (ret string) {
	prereq := parse(t, in)

	done := make(map[string]bool)

	var order []string
	for len(prereq) > 0 {
		var ready []string
	nextStep:
		for step, prev := range prereq {
			for _, p := range prev {
				if !done[p] {
					continue nextStep
				}
			}
			ready = append(ready, step)
		}
		if len(ready) == 0 {
			t.Fatalf("No steps ready; done: %v", done)
		}
		sort.Strings(ready)
		do := ready[0]
		done[do] = true
		delete(prereq, do)
		order = append(order, do)
	}

	return strings.Join(order, "")
}

func part2(t *testing.T, in string, workers, extraTime int) (ret int) {
	prereq := parse(t, in)

	done := make(map[string]bool)
	started := make(map[string]bool)

	var order []string
	var queue EventHeap
	var now int
	for len(prereq) > 0 {
		var ready []string
	nextStep:
		for step, prev := range prereq {
			if started[step] {
				continue
			}
			for _, p := range prev {
				if !done[p] {
					continue nextStep
				}
			}
			ready = append(ready, step)
		}
		if len(ready) == 0 && len(queue) == 0 {
			t.Fatalf("No steps ready; done: %v", done)
		}
		sort.Strings(ready)

		for workers > 0 && len(ready) > 0 {
			do := ready[0]
			started[do] = true
			end := now + int(do[0]-'A'+1) + extraTime
			t.Logf("Starting %v at %v", do, now)
			heap.Push(&queue, TimedEvent{
				T: end,
				K: do,
				F: func() {
					done[do] = true
					delete(prereq, do)
					order = append(order, do)
					workers++
				},
			})
			ready = ready[1:]
			workers--
		}

		next := heap.Pop(&queue).(TimedEvent)
		now = next.T
		next.F()
		t.Logf("Finished %v at %v", next.K, now)
	}

	return now
}

type TimedEvent struct {
	T int
	K string
	F func()
}

type EventHeap []TimedEvent

func (h EventHeap) Len() int              { return len(h) }
func (h EventHeap) Swap(i, j int)         { h[i], h[j] = h[j], h[i] }
func (h *EventHeap) Push(x interface{})   { *h = append(*h, x.(TimedEvent)) }
func (h *EventHeap) Pop() (x interface{}) { x, *h = (*h)[len(*h)-1], (*h)[:len(*h)-1]; return }

func (h EventHeap) Less(i, j int) bool {
	if a, b := h[i].T, h[j].T; a != b {
		return a < b
	}
	if a, b := h[i].K, h[j].K; a != b {
		return a < b
	}
	return false
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"part1 example 0", strings.TrimSpace(`
Step C must be finished before step A can begin.
Step C must be finished before step F can begin.
Step A must be finished before step B can begin.
Step A must be finished before step D can begin.
Step B must be finished before step E can begin.
Step D must be finished before step E can begin.
Step F must be finished before step E can begin.
		`), "CABDFE"},
		{"part1 answer", advent.ReadFile(t, "input.txt"), "BETUFNVADWGPLRJOHMXKZQCISY"},
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
		name    string
		in      string
		workers int
		time    int
		want    int
	}{
		{"part2 example 0", strings.TrimSpace(`
Step C must be finished before step A can begin.
Step C must be finished before step F can begin.
Step A must be finished before step B can begin.
Step A must be finished before step D can begin.
Step B must be finished before step E can begin.
Step D must be finished before step E can begin.
Step F must be finished before step E can begin.
		`), 2, 0, 15},
		{"part2 answer", advent.ReadFile(t, "input.txt"), 5, 60, 848},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part2(t, test.in, test.workers, test.time), test.want; got != want {
				t.Errorf("part2(%#v)\n = %#v, want %#v", test.in, got, want)
			}
		})
	}
}

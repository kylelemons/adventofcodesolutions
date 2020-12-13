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
	"sync"
	"testing"

	"github.com/kylelemons/adventofcodesolutions/2019/intcode"
	"github.com/kylelemons/adventofcodesolutions/advent"
)

type Helper struct {
	Base *intcode.Program
}

func parseInput(t *testing.T, in string) *Helper {
	helper := &Helper{
		Base: intcode.Compile(t, in),
	}
	return helper
}

func (h *Helper) RunNetwork(t *testing.T, nodeCount int) (part1, part2 int) {
	type packet struct {
		Dest int
		X, Y int

		temp int
	}
	type queue struct {
		sync.Mutex
		*sync.Cond
		partial *packet
		pending []packet
	}

	queues := make([]*queue, nodeCount)
	nodes := make([]*intcode.Program, nodeCount)
	for i := range queues {
		q := new(queue)
		q.Cond = sync.NewCond(&q.Mutex)

		queues[i] = q
		nodes[i] = h.Base.Snapshot()
	}

	// NAT
	natStore := make(chan packet)
	natTickle := make(chan bool, 1)
	done := make(chan bool)
	go func() {
		var pending *packet

		seen := make(map[int]bool)
		scan := func() {
			var active []int
			for id, q := range queues {
				q.Lock()
				defer q.Unlock()
				if q.partial == nil && len(q.pending) == 0 {
					continue
				}
				active = append(active, id)
			}
			if act := len(active); act > 0 {
				// t.Logf("NAT: Scan: %v are active (no action)", active)
				return
			}
			if pending == nil {
				// t.Logf("NAT: No nodes are active, but no packet is pending")
				for _, q := range queues {
					q.Broadcast()
				}
				return
			}
			// t.Logf("NAT: No nodes are active; queueing pending packet: %+v", pending)
			if seen[pending.Y] && part2 == 0 {
				part2 = pending.Y
				t.Logf("Part 2: %d", part2)
				close(done)
			}
			seen[pending.Y] = true

			// queue[0].Mutex is already held
			queues[0].pending = append(queues[0].pending, *pending)
			pending = nil
		}
		for {
			select {
			case p := <-natStore:
				// t.Logf("NAT: Storing %+v", p)
				pending = &p
				if part1 == 0 {
					t.Logf("Part 1: %v", p.Y)
					part1 = p.Y
				}
				continue // We received a packet, don't scan
			case <-natTickle:
				scan()
			}
		}
	}()

	for id, n := range nodes {
		id, n := id, n

		setup := make(chan int, 1)
		setup <- id

		n.Input = func() (v int) {
			// defer func() { t.Logf("%d.Input() = %v", id, v) }()

			select {
			case v = <-setup:
				return
			default:
			}

			q := queues[id]
			q.Lock()
			defer q.Unlock()

			switch {
			case q.partial == nil && len(q.pending) == 0:
				if id == 0 {
					select {
					case natTickle <- true:
					default:
					}
				} else {
					q.Wait()
				}
				return -1
			case q.partial == nil:
				q.partial = &q.pending[0]
				v = q.partial.X
			default:
				v = q.partial.Y
				q.partial = nil
				q.pending = q.pending[1:]
			}
			return
		}

		var partial packet
		n.Output = func(v int) {
			// t.Logf("%d.Output(%v)", id, v)

			partial.temp++
			switch partial.temp {
			case 1:
				partial.Dest = v
				return
			case 2:
				partial.X = v
				return
			case 3:
				partial.Y = v
			}

			// ensure partial gets reset after use
			defer func() { partial.temp = 0 }()

			if partial.Dest > len(queues) {
				natStore <- partial
				return
			}

			q := queues[partial.Dest]
			q.Lock()
			defer q.Unlock()
			q.pending = append(q.pending, partial)
			q.Broadcast()
		}

		// if id == 0 {
		// 	n.Debugf = t.Logf
		// }

		go n.Run(t)
		defer n.Halt()
	}

	<-done
	return
}

func TestPart1AndPart2FullSimulation(t *testing.T) {
	t.Skip("Flaky (not sure why)")

	tests := []struct {
		name  string
		in    string
		want1 int
		want2 int
	}{
		{"answers", advent.ReadFile(t, "input.txt"), 22134, 16084},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			helper := parseInput(t, test.in)
			part1, part2 := helper.RunNetwork(t, 50)
			if got, want := part1, test.want1; got != want {
				t.Errorf("part1(...)\n = %#v, want %#v", got, want)
			}
			if got, want := part2, test.want2; got != want {
				t.Errorf("part2(...)\n = %#v, want %#v", got, want)
			}
		})
	}
}

func (h *Helper) RunFast(t *testing.T, numNodes int) (part1, part2 int) {
	type packet struct {
		Dest int
		X, Y int
	}
	var pending []packet

	// I/O helpers:
	var (
		// rawInput provides the slice values as input and then halts.
		rawInput = func(id int, node *intcode.Program, inputs ...int) func() int {
			return func() (v int) {
				if len(inputs) == 0 {
					node.Halt()
					return -1
				}
				v, inputs = inputs[0], inputs[1:]
				// t.Logf("%2d : Input() = %v", id, v)
				return
			}
		}

		// packetInput provides the packet as input and then halts.
		packetInput = func(id int, node *intcode.Program, current packet) func() int {
			inputs := []int{current.X, current.Y}
			return rawInput(id, node, inputs...)
		}

		// packetOutput queues up any messages the node sends while it runs.
		packetOutput = func(id int, node *intcode.Program) func(int) {
			var outputs []int
			return func(v int) {
				// t.Logf("%2d : Output(%v)", id, v)
				outputs = append(outputs, v)
				if len(outputs) < 3 {
					return
				}
				defer func() { outputs = nil }()

				packet := packet{
					Dest: outputs[0],
					X:    outputs[1],
					Y:    outputs[2],
				}
				// t.Logf("Queueing packet %+v for delivery", packet)
				pending = append(pending, packet)
			}
		}

		traceExec = func(id int, node *intcode.Program) func(string, ...interface{}) {
			return func(format string, args ...interface{}) {
				t.Logf("%2d : TRACE %p : %s", id, node.Output, fmt.Sprintf(format, args...))
			}
		}
		_ = traceExec
	)

	// Setup nodes
	nodes := make([]*intcode.Program, numNodes)
	for id := range nodes {
		node := h.Base.Snapshot()

		node.Output = packetOutput(id, node)
		node.Input = rawInput(id, node, id)
		node.Run(t)

		nodes[id] = node
	}

	var nat *packet
	seen := make(map[int]bool)
	for part1 == 0 || part2 == 0 {
		for len(pending) > 0 {
			// Deliver the next message.
			current := pending[0]
			pending = pending[1:]

			id := current.Dest
			if id == 255 {
				// t.Logf("NAT: Storing %+v", current)
				nat = &current
				if part1 == 0 {
					part1 = current.Y
					t.Logf("Part 1: %d", part1)
				}
				continue
			}

			node := nodes[id]
			node.Input = packetInput(id, node, current)
			node.Output = packetOutput(id, node)

			// Run until it asks for a second packet.
			node.Run(t)
		}

		// Notify all nodes that they have nothing waiting.
		for id, node := range nodes {
			node.Input = rawInput(id, node, -1)
			node.Output = packetOutput(id, node)
			node.Run(t)
		}

		// Queue up the NAT packet if we're still idle.
		if len(pending) == 0 {
			if nat == nil {
				t.Fatalf("Still idle with no NAT packet")
				continue
			}
			// t.Logf("NAT: Idle detected, queueing %+v", nat)

			if seen[nat.Y] {
				part2 = nat.Y
				t.Logf("Part 2: %d", part2)
			}
			seen[nat.Y] = true

			pending = append(pending, packet{
				Dest: 0,
				X:    nat.X,
				Y:    nat.Y,
			})
			nat = nil
		}
	}

	return
}

func TestPart1AndPart2Fast(t *testing.T) {
	tests := []struct {
		name  string
		in    string
		want1 int
		want2 int
	}{
		{"answers", advent.ReadFile(t, "input.txt"), 22134, 16084},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			helper := parseInput(t, test.in)
			part1, part2 := helper.RunFast(t, 50)
			if got, want := part1, test.want1; got != want {
				t.Errorf("part1(...)\n = %#v, want %#v", got, want)
			}
			if got, want := part2, test.want2; got != want {
				t.Errorf("part2(...)\n = %#v, want %#v", got, want)
			}
		})
	}
}

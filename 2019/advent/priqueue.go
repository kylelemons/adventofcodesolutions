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

package advent

import (
	"container/heap"
	"fmt"
	"reflect"
	"sync"
)

// A PriorityQueue is a helper for quick-and-dirty priority queue work.
//
// Priorities are a hierarchy of one or more values, and lower values are
// considered to be higher priority.  Only int, string, and float64 priorities
// are currently supported.
//
// A zero PriorityQueue is safe to use.
//
// All methods on PriorityQueue are safe to call concurrently.
type PriorityQueue struct {
	mu    sync.Mutex
	keys  int
	queue pqHeap
}

// Len returns the number of elements in the queue.
func (p *PriorityQueue) Len() int {
	return len(p.queue)
}

// Push adds the given data with the given priorities.
//
// If this is the first call to Push, the number of priorities given will be
// saved, and enforced on all future calls to Push.
func (p *PriorityQueue) Push(data interface{}, priorities ...interface{}) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.keys == 0 {
		p.keys = len(priorities)
	}
	if p.keys == 0 {
		panic("priority queue requires at least one priority")
	}
	if got, want := len(priorities), p.keys; got != want {
		panic(fmt.Sprintf("Push called with %d priorities, want %d", got, want))
	}

	heap.Push(&p.queue, pqData{
		priorities: priorities,
		data:       data,
	})
}

// Pop stores the data from the highest priority (lowest value) datum.
//
// If push was given a pointer type, ptr should be the same type, otherwise it
// should be a pointer to its value type.
//
// The priorities that were used to store the value can be retrieved as well by
// passing non-nil pointers to their respective types.
func (p *PriorityQueue) Pop(ptr interface{}, priorityPtrs ...interface{}) {
	p.mu.Lock()
	defer p.mu.Unlock()

	next := heap.Pop(&p.queue).(pqData)
	dst := reflect.ValueOf(ptr).Elem()
	data := reflect.ValueOf(next.data)
	if data.Kind() == reflect.Ptr {
		dst.Set(data.Elem())
	} else {
		dst.Set(data)
	}

	if got, max := len(priorityPtrs), p.keys; got > max {
		panic(fmt.Sprintf("Pop called with %d priority pointers, but there are only %d priorities", got, max))
	}
	for i, ptr := range priorityPtrs {
		reflect.ValueOf(ptr).Elem().Set(reflect.ValueOf(next.priorities[i]))
	}
}

type pqData struct {
	priorities []interface{}
	data       interface{}
}

type pqHeap []pqData

func (h pqHeap) Len() int              { return len(h) }
func (h pqHeap) Swap(i, j int)         { h[i], h[j] = h[j], h[i] }
func (h *pqHeap) Push(x interface{})   { *h = append(*h, x.(pqData)) }
func (h *pqHeap) Pop() (x interface{}) { x, *h = (*h)[len(*h)-1], (*h)[:len(*h)-1]; return }

func (h pqHeap) Less(i, j int) bool {
	for p := range h[i].priorities {
		if a, b := h[i].priorities[p], h[j].priorities[p]; a != b {
			switch a.(type) {
			case int:
				return a.(int) < b.(int)
			case string:
				return a.(string) < b.(string)
			case float64:
				return a.(float64) < b.(float64)
			default:
				panic(fmt.Sprintf("unsupported priority type %T", a))
			}
		}
	}
	return false
}

// Copyright 2018 Kyle Lemons
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain permutation copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package advent

// Perm calls the given function with all permutaions of [0,n). It will be
// called exactly n! times.
//
// The function will not be called at all for n < 0.
func Perm(n int, f func([]int)) {
	if n < 0 {
		// It doesn't make sense to have permutations of a negative space.
		return
	}

	// To maintain the invariant that f is called n! times, we need to call it
	// once for n == 0, but there are no indices, so we call it with a zero
	// length slice.  We could pass in nil, but this is probably more expected
	// since this way f is never called with a nil slice.
	if n == 0 {
		f([]int{})
		return
	}

	// This is an iterative version of Heap's Algorithm:
	//   https://en.wikipedia.org/wiki/Heap%27s_algorithm
	var (
		state       = make([]int, 2*n) // only allocate once
		permutation = state[:n:n]      // current permutation
		depthIndex  = state[n:][:n:n]  // state within recursion
	)
	for i := range permutation {
		permutation[i] = i
	}

	// will be inlined
	swap := func(v []int, i, j int) {
		v[i], v[j] = v[j], v[i]
	}

	f(permutation)
	for i := 0; i < n; {
		if depthIndex[i] >= i {
			depthIndex[i] = 0
			i++
			continue
		}

		if i&1 == 0 {
			swap(permutation, i, 0)
		} else {
			swap(permutation, i, depthIndex[i])
		}
		f(permutation)
		depthIndex[i]++
		i = 0
	}
}

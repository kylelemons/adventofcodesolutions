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
	"fmt"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPerm(t *testing.T) {
	tests := []struct {
		n int
		v []string
	}{
		{
			n: -1,
			v: nil,
		},
		{
			n: 0,
			v: []string{
				"[]",
			},
		},
		{
			n: 1,
			v: []string{
				"[0]",
			},
		},
		{
			n: 3,
			v: []string{
				"[0 1 2]",
				"[0 2 1]",
				"[1 0 2]",
				"[1 2 0]",
				"[2 0 1]",
				"[2 1 0]",
			},
		},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("n=%d", test.n), func(t *testing.T) {
			var out []string
			Perm(test.n, func(v []int) { out = append(out, fmt.Sprint(v)) })
			sort.Strings(out)
			if diff := cmp.Diff(out, test.v); diff != "" {
				t.Errorf("Perm(%d) produced incorrect results: (-got +want)\n%s", test.n, diff)
			}
		})
	}
}

func TestPrimes(t *testing.T) {
	want := []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29}
	primes := Primes()

	var got []int
	for len(got) < len(want) {
		got = append(got, primes())
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Primes() returned incorrect sequence: (-got +want)\n%s", diff)
	}
}

func BenchmarkPerm(b *testing.B) {
	for _, n := range []int{0, 1, 3, 5, 7, 10} {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			var total int
			for i := 0; i < b.N; i++ {
				Perm(n, func(v []int) { total += v[0] })
			}
			_ = total
		})
	}
}

func BenchmarkPrimes(b *testing.B) {
	for _, n := range []int{100, 1000, 10000} {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			b.SetBytes(int64(n))
			var count, sum int
			for i := 0; i < b.N; i++ {
				primes := Primes()
				for j := 0; j < n; j++ {
					sum += primes()
					count++
				}
			}
			b.Logf("%d primes (total %d)", count, sum)
		})
	}
}

func ExamplePerm() {
	colors := []string{
		"red",
		"blue",
		"green",
	}
	Perm(len(colors), func(idx []int) {
		var comma string
		for _, i := range idx {
			fmt.Printf("%s%s", comma, colors[i])
			comma = ", "
		}
		fmt.Println()
	})

	// Output:
	// red, blue, green
	// blue, red, green
	// green, red, blue
	// red, green, blue
	// blue, green, red
	// green, blue, red
}

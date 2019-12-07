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
			n: 0,
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

func BenchmarkPerm(b *testing.B) {
	for _, n := range ([]int{0,1,3,5,7,10}) {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			var total int
			for i := 0; i < b.N; i++ {
				Perm(n, func(v []int){total += v[0]})
			}
			_ = total
		})
	}
}

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

package advent

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFactorize(t *testing.T) {
	tests := []struct {
		n    int
		want map[int]int
	}{
		{
			n: 12300,
			want: map[int]int{
				2:  2,
				3:  1,
				5:  2,
				41: 1,
			},
		},
		{
			n: 100200300400,
			want: map[int]int{
				2:        4,
				5:        2,
				23:       1,
				10891337: 1,
			},
		},
	}
	for _, test := range tests {
		t.Run(fmt.Sprint(test.n), func(t *testing.T) {
			if diff := cmp.Diff(factorize(test.n), test.want); diff != "" {
				t.Errorf("factorize(%d) returned incorrect factors: (-got +want)\n%s", test.n, diff)
			}
		})
	}
}

func TestLCM(t *testing.T) {
	tests := []struct {
		v    []int
		want int
	}{
		{
			v:    []int{2 * 2 * 3, 3 * 3 * 5, 5 * 5 * 7},
			want: 2 * 2 * 3 * 3 * 5 * 5 * 7,
		},
		{
			v:    []int{16, 99, 117, 412, 1010},
			want: 1071092880,
		},
		{
			v:    []int{924, 924, 2772, 2772},
			want: 2772,
		},
		{
			v:    []int{924, 2772, 2772, 924},
			want: 2772,
		},
	}
	for _, test := range tests {
		t.Run(fmt.Sprint(test.v), func(t *testing.T) {
			if got, want := LCM(test.v...), test.want; got != want {
				t.Errorf("LCM(%v) = %v, want %v", test.v, got, want)
			}
		})
	}
}

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

// Package advent contains utilities for Advent of Code.
package advent

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRecords(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{"", nil},
		{"\n", nil},
		{"\n\n", nil},
		{"\n\n\n", nil},
		{"\n\n\n\n", nil},
		{"a", []string{"a"}},
		{"a\n", []string{"a"}},
		{"a\n\n", []string{"a"}},
		{"a\n\n\n", []string{"a"}},
		{"a\n\n\n\n", []string{"a"}},
		{"a\n\nb", []string{"a", "b"}},
		{"a\n\nb\n", []string{"a", "b"}},
		{"a\n\nb\n\n", []string{"a", "b"}},
		{"\na\n\nb\n\n", []string{"a", "b"}},
		{"a\n\n\nb\n\n", []string{"a", "b"}},
		{"a\n\nb\n\n\n", []string{"a", "b"}},
		{"a\n\n\n\nb\n\n", []string{"a", "b"}},
		{"a\n\n\n\n\nb\n\n", []string{"a", "b"}},
	}
	for _, test := range tests {
		if got := Records(test.input).All(t); !cmp.Equal(got, test.want) {
			t.Errorf("Records(%q) = %#q, want %#q", test.input, got, test.want)
		}
	}
}

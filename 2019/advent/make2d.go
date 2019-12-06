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

// Make2D makes a 2d slice of bytes of the given dimensions.
func Make2D(rows, cols int) [][]byte {
	back := make([]byte, rows*cols)
	out := make([][]byte, rows)
	for r := range out {
		out[r] = back[r*cols:][:cols:cols]
	}
	return out
}

// Make2DInts makes a 2d slice of bytes of the given dimensions.
func Make2DInts(rows, cols int) [][]int {
	back := make([]int, rows*cols)
	out := make([][]int, rows)
	for r := range out {
		out[r] = back[r*cols:][:cols:cols]
	}
	return out
}

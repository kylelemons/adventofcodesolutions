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

// Package coords contains helper functions for coordinate systems for
// Advent of Code problems.
package coords

import "fmt"

// An Coord represents a coordinate or vector in an X-Y coordinate system wher
// +Y and +R are "South" and +X and +C are "East".
type Coord struct {
	x, y int
}

// XY returns the (x,y) coordinate as a Coord.
func XY(x, y int) Coord { return Coord{x: x, y: y} }

// RC returns the (r,c) coordinate as a Coord.
func RC(r, c int) Coord { return Coord{x: c, y: r} }

// R returns the row of the coord in the Row-Column coordinate system.
func (c Coord) R() int { return c.y }

// C returns the column of the coord in the Row-Column coordinate system.
func (c Coord) C() int { return c.x }

// X returns the x coordinate in the X-Y coordinate system.
func (c Coord) X() int { return c.x }

// Y returns the y coordinate in the X-Y coordinate system.
func (c Coord) Y() int { return c.y }

// Add adds the vector to the point.
func (c Coord) Add(v Vector) Coord { return Coord{x: c.x + v.x, y: c.y + v.y} }

// Sub returns a vector such that c2 + v = c.
func (c Coord) Sub(c2 Coord) (v Vector) { return Vector{x: c2.x - c.x, y: c2.y - c.y} }

// String returns a string value of the coordinate.
func (c Coord) String() string {
	return fmt.Sprintf("(x=%+d, y=%+d)", c.x, c.y)
}

// RCString returns a string value of the coordinate in the R-C system.
func (c Coord) RCString() string {
	return fmt.Sprintf("%d,%d", c.R(), c.C())
}

// XYString returns a string value of the coordinate in the R-C system.
func (c Coord) XYString() string {
	return fmt.Sprintf("%d,%d", c.X(), c.Y())
}

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
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

// Scanner is a helper for doing simple linewise scanning.
type Scanner string

// Scan scans the string into the given pointers using fmt.Sscan.
func (s Scanner) Scan(t *testing.T, ptrs ...interface{}) bool {
	t.Helper()
	if _, err := fmt.Sscan(string(s), ptrs...); err != nil {
		if err == io.EOF {
			return true
		}
		t.Fatalf("Sscan: %s", err)
	}
	return true
}

// Extract extracts sequential capture groups from the scanner into the given pointers.
func (s Scanner) Extract(t *testing.T, re string, ptrs ...interface{}) bool {
	t.Helper()
	r, err := regexp.Compile(re)
	if err != nil {
		t.Fatalf("bad regexp %q: %s", re, err)
	}
	if got, want := len(ptrs), r.NumSubexp(); got != want {
		t.Fatalf("bad scan: %d pointers found, want %d (number of groups)", got, want)
	}
	matches := r.FindStringSubmatch(string(s))
	if matches == nil {
		return false
	}
	for i, ptr := range ptrs {
		val := matches[i+1]
		if got, want := reflect.TypeOf(ptr).Kind(), reflect.Ptr; got != want {
			t.Fatalf("can't scan into group %d: got %v, want %v", i+1, got, want)
		}
		if _, err := fmt.Sscan(val, ptr); err != nil {
			t.Fatalf("failed to scan %q into %T: %s", val, ptr, err)
		}
	}
	return true
}

// ReadFile reads the named file and returns it as a string.
func ReadFile(t *testing.T, filename string) string {
	t.Helper()
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatalf("failed to read %q: %s", filename, err)
	}
	return string(data)
}

// Lines calls each once for each line in s, parsing each line into the pointers using fmt.Sscan.
func Lines(t *testing.T, s string, each interface{}) {
	fval := reflect.ValueOf(each)
	pointers, values := inputsFor(fval)
	lines := bufio.NewScanner(strings.NewReader(s))
	for lines.Scan() {
		Scanner(lines.Text()).Scan(t, pointers...)
		fval.Call(values)
	}
}

func inputsFor(fval reflect.Value) (scanInputs []interface{}, callInputs []reflect.Value) {
	var values, pointers []reflect.Value
	var rawPointers []interface{}
	ftyp := fval.Type()
	for i, n := 0, ftyp.NumIn(); i < n; i++ {
		it := ftyp.In(i)
		ptr := reflect.New(it)
		values = append(values, ptr.Elem())
		pointers = append(pointers, ptr)
		rawPointers = append(rawPointers, ptr.Interface())
	}
	return rawPointers, values
}

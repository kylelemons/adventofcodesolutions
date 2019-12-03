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
	"unicode/utf8"
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
	return strings.TrimSpace(string(data))
}

// Delimited is a helper for processing delimited inputs.
type Delimited struct {
	Scanner *bufio.Scanner
}

// Lines returns a Delimited helper for line-delimited inputs.
func Lines(input string) *Delimited {
	return &Delimited{
		Scanner: bufio.NewScanner(strings.NewReader(input)),
	}
}

func withSplitter(input string, splitter bufio.SplitFunc) *Delimited {
	s := bufio.NewScanner(strings.NewReader(input))
	s.Split(splitter)
	return &Delimited{
		Scanner: s,
	}
}

// Words returns a Delimited helper for word-delimited inputs.
func Words(input string) *Delimited { return withSplitter(input, bufio.ScanWords) }

// Split returns a Delimited helper for rune-delimited inputs (e.g. comma).
func Split(input string, delim rune) *Delimited {
	// stolen from bufio.ScanWords
	return withSplitter(input, func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		// Scan until rune, marking end of token.
		for width, i := 0, 0; i < len(data); i += width {
			var r rune
			r, width = utf8.DecodeRune(data[i:])
			if r == delim {
				return i + width, data[:i], nil
			}
		}
		// If we're at EOF, we have a final, non-empty, non-terminated word. Return it.
		if atEOF && len(data) > 0 {
			return len(data), data, nil
		}
		// Request more data.
		return 0, nil, nil
	})
}

// Scan calls Scanner.Scan on each item to provide inputs to the each function.
//
// Example:
//   advent.Lines(input).Scan(t, func(command string, arg int) { ... })
func (d *Delimited) Scan(t *testing.T, each interface{}) {
	fval := reflect.ValueOf(each)
	pointers, values := inputsFor(fval)
	for d.Scanner.Scan() {
		Scanner(d.Scanner.Text()).Scan(t, pointers...)
		fval.Call(values)
	}
}

// Extract calls Scanner.Extract on each item to provide inputs to the each function.
//
// Example:
//   advent.Lines(input).Extract(t, input, `([ULDR])(\d+)`, func(dir string, steps int) { ... })
func (d *Delimited) Extract(t *testing.T, re string, each interface{}) {
	fval := reflect.ValueOf(each)
	pointers, values := inputsFor(fval)
	for d.Scanner.Scan() {
		Scanner(d.Scanner.Text()).Extract(t, re, pointers...)
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

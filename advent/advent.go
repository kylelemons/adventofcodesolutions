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
	"log"
	"reflect"
	"regexp"
	"strings"
	"unicode/utf8"
)

// Scanner is a helper for doing simple linewise scanning.
type Scanner string

// An OptionalT represents the poritons of *testing.T that are used by the
// advent package.
//
// Any API that accepts an OptionalT can also be passed nil, in which case
// a default implementation will be created and used in its place.
type OptionalT interface {
	Fatalf(format string, args ...interface{})
	Helper()
}

type defaultT struct{}

func (t defaultT) Helper() {}

func (t defaultT) Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

func maybeT(t OptionalT) OptionalT {
	if t != nil {
		return t
	}
	return defaultT{}
}

// Scan scans the string into the given pointers using fmt.Sscan.
func (s Scanner) Scan(t OptionalT, ptrs ...interface{}) {
	t = maybeT(t)
	t.Helper()
	if _, err := fmt.Sscan(string(s), ptrs...); err != nil {
		if err == io.EOF {
		}
		t.Fatalf("Sscan: %s", err)
	}
}

// Extract extracts sequential capture groups from the scanner into the given pointers.
func (s Scanner) Extract(t OptionalT, re string, ptrs ...interface{}) {
	t = maybeT(t)
	t.Helper()
	if !s.CanExtract(t, re, ptrs...) {
		t.Fatalf("Input %q does not match /%s/", s, re)
	}
}

// CanExtract returns true if it can extract sequential capture groups from the
// scanner into the given pointers.
//
// CanExtact returns true if the line is fully scanned, false if the regex does
// not match, and will fatal out if the regex is invalid or if a matched value
// cannot be correctly stored.
func (s Scanner) CanExtract(t OptionalT, re string, ptrs ...interface{}) bool {
	t = maybeT(t)
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
func ReadFile(t OptionalT, filename string) string {
	t = maybeT(t)
	t.Helper()
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatalf("failed to read %q: %s", filename, err)
	}
	return strings.TrimRight(string(data), "\n")
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

// Each calls f for each successive non-empty token with the token index and
// a scanner to use in parsing the token.
func (d *Delimited) Each(each func(i int, token Scanner)) {
	for i := 0; d.Scanner.Scan(); {
		token := d.Scanner.Text()
		if token == "" {
			continue
		}
		each(i, Scanner(token))
		i++
	}
}

// Scan calls Scanner.Scan on each item to provide inputs to the each function.
//
// Example:
//   advent.Lines(input).Scan(t, func(command string, arg int) { ... })
func (d *Delimited) Scan(t OptionalT, each interface{}) {
	t = maybeT(t)
	t.Helper()

	fval := reflect.ValueOf(each)
	pointers, values := inputsFor(fval)
	d.Each(func(_ int, token Scanner) {
		token.Scan(t, pointers...)
		fval.Call(values)
	})
}

// Extract calls Scanner.Extract on each item to provide inputs to the each function.
//
// Example:
//   advent.Lines(input).Extract(t, input, `([ULDR])(\d+)`, func(dir string, steps int) { ... })
func (d *Delimited) Extract(t OptionalT, re string, each interface{}) {
	t = maybeT(t)
	t.Helper()

	fval := reflect.ValueOf(each)
	pointers, values := inputsFor(fval)
	d.Each(func(_ int, token Scanner) {
		token.Extract(t, re, pointers...)
		fval.Call(values)
	})
}

// inputsFor returns two slices based on the given reflective function value:
//  - a []interface{} for passing to fmt.Scan* corresponding 1:1 to params
//  - a []reflect.Value for calling the fval function with the internal params
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

// Fields returns pointers to the slice's fields in source order, suitable for
// passing to Scan, Extract, or CanExtract.
func Fields(structPtr interface{}) (fieldPtrs []interface{}) {
	v := reflect.ValueOf(structPtr).Elem()
	for i, n := 0, v.Type().NumField(); i < n; i++ {
		fieldPtrs = append(fieldPtrs, v.Field(i).Addr().Interface())
	}
	return
}

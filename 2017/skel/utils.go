package aocday

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"regexp"
	"testing"
)

type scanner string

func (s scanner) scan(t *testing.T, re string, ptrs ...interface{}) bool {
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

func read(t *testing.T, filename string) string {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatalf("failed to read %q: %s", filename, err)
	}
	return string(data)
}

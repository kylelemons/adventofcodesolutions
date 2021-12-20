package advent

import "testing"

func Must[T any](value T, err error) func(t *testing.T) T {
	return func(t *testing.T) T {
		t.Helper()
		if err != nil {
			t.Fatalf("Must[%T]: %s", *new(T), err)
		}
		return value
	}
}

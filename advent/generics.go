//go:build go1.18

package advent

func Must[T any](value T, err error) func(t OptionalT) T {
	return func(t OptionalT) T {
		t = maybeT(t)
		t.Helper()
		if err != nil {
			t.Fatalf("Must[%T]: %s", *new(T), err)
		}
		return value
	}
}

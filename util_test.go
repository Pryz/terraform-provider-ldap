package main

import (
	"testing"
)

func TestObjectClassesRoundTrip(t *testing.T) {

	tests := []map[string]struct{}{
		map[string]struct{}{
			"one":   struct{}{},
			"two":   struct{}{},
			"three": struct{}{},
			"four":  struct{}{},
			"five":  struct{}{},
			"six":   struct{}{},
		},

		map[string]struct{}{
			"unquoted":        struct{}{},
			"\"quoted\"":      struct{}{},
			"  with blanks  ": struct{}{},
		},
	}

	for _, test := range tests {
		set := unmarshalObjectClasses(marshalObjectClasses(test))
		if len(test) != len(set) {
			t.Error("for", test, "result was", set)
		}
		for item := range test {
			if _, ok := set[item]; !ok {
				t.Error("for", test, "item", item, "is missing")
			}
		}
	}
}

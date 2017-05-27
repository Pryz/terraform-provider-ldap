package main

import "testing"

func setup() (*set, *set) {
	s1 := &set{
		data: map[string]struct{}{},
	}
	s1.Add("a")
	s1.Add("b")
	s1.Add("c")
	s1.Add("d")

	s1.Add("e")
	s1.Add("f")

	s2 := NewSet()
	s2.Add("a")
	s2.Add("b")
	s2.Add("c")
	s2.Add("d")

	s2.Add("g")
	s2.Add("h")

	return s1, s2
}

func TestAdd(t *testing.T) {
	s1 := &set{
		data: map[string]struct{}{},
	}
	if !s1.Add("a") {
		t.Errorf("Invalid result (false) from add of a new string")
	}
	if s1.Add("a") != false {
		t.Error("Invalid result (true) from add of an existing string")
	}
	if s1.Add("") != false {
		t.Error("Invalid result (true) from add of a null string")
	}
}

func TestIntersection(t *testing.T) {
	s1, s2 := setup()
	intersection := s1.Intersection(s2)
	if intersection.Len() != 4 {
		t.Errorf("Invalid intersection, expected 4 got %d", intersection.Len())
	}
	for _, v := range []string{"a", "b", "c", "d"} {
		if !intersection.Contains(v) {
			t.Errorf("Invalid intersection, it does not contain %q", v)
		}
	}
}

func TestDifference(t *testing.T) {
	s1, s2 := setup()
	difference := s1.Difference(s2)
	if difference.Len() != 2 {
		t.Errorf("Invalid difference, expected 2 got %d", difference.Len())
	}
	for _, v := range []string{"e", "f"} {
		if !difference.Contains(v) {
			t.Errorf("Invalid difference, it does not contain %q", v)
		}
	}

	difference = s2.Difference(s1)
	if difference.Len() != 2 {
		t.Errorf("Invalid difference, expected 2 got %d", difference.Len())
	}
	for _, v := range []string{"g", "h"} {
		if !difference.Contains(v) {
			t.Errorf("Invalid difference, it does not contain %q", v)
		}
	}
}

func TestUnion(t *testing.T) {
	s1, s2 := setup()
	union := s1.Union(s2)
	if union.Len() != 8 {
		t.Errorf("Invalid union, expected 8 got %d", union.Len())
	}
	for _, v := range []string{"a", "b", "c", "d", "e", "f", "g", "h"} {
		if !union.Contains(v) {
			t.Errorf("Invalid union, it does not contain %q", v)
		}
	}
}

func TestSymmetricDifference(t *testing.T) {
	s1, s2 := setup()
	simmetric := s1.SymmetricDifference(s2)
	if simmetric.Len() != 4 {
		t.Errorf("Invalid simmetric difference, expected 4 got %d", simmetric.Len())
	}
	for _, v := range []string{"e", "f", "g", "h"} {
		if !simmetric.Contains(v) {
			t.Errorf("Invalid simmetric difference, it does not contain %q", v)
		}
	}
}

func TestAllTogether(t *testing.T) {
	s1, s2 := setup()
	if !s1.Union(s2).Difference(s1.Intersection(s2)).Equals(s1.SymmetricDifference(s2)) {
		t.Errorf("Invalid overall concatenation")
	}
}

func TestToString(t *testing.T) {
	s1, _ := setup()
	if s1.String() != "{ \"a\", \"b\", \"c\", \"d\", \"e\", \"f\", }" {
		t.Errorf("Invalid string, got %s", s1.String())
	}
}

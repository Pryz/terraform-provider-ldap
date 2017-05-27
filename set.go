package main

import (
	"bytes"
	"fmt"
	"sort"
)

// Set represents a set of strings.
type Set struct {
	data map[string]struct{}
}

// NewSet creates a new string set.
func NewSet(values ...string) *Set {
	set := &Set{
		data: map[string]struct{}{},
	}
	for _, value := range values {
		set.Add(value)
	}
	return set
}

// Add adds a string to the set, returning whether the string was actually added;
// if the same value was already in the set, it will therefore return false.
func (s *Set) Add(attribute string) bool {
	if attribute == "" {
		return false
	}
	_, ok := s.data[attribute]
	if !ok {
		s.data[attribute] = struct{}{}
		return true
	}
	return false
}

// Remove removes the given value from the set, returning if such a value was
// present and therefore actually removed.
func (s *Set) Remove(attribute string) bool {
	if attribute == "" {
		return false
	}
	_, ok := s.data[attribute]
	if ok {
		delete(s.data, attribute)
		return true
	}
	return false
}

// Clear removes all values from the set.
func (s *Set) Clear() {
	for k := range s.data {
		delete(s.data, k)
	}
}

// Equals checks if two sets have the same elements.
func (s *Set) Equals(other *Set) bool {
	if other == nil {
		return false
	}
	if s.Len() != other.Len() {
		return false
	}
	for k := range s.data {
		if !other.Contains(k) {
			return false
		}
	}
	return true
}

// Contains checks if the set contains the given value.
func (s *Set) Contains(value string) bool {
	if value == "" {
		return false
	}
	_, ok := s.data[value]
	return ok
}

// Len returns the number of elements in the set.
func (s *Set) Len() int {
	return len(s.data)
}

// Difference compares the set against another one and returns a new set
// containing the elements that are only in this set.
func (s *Set) Difference(other *Set) *Set {
	if other == nil {
		return nil
	}
	result := &Set{
		data: map[string]struct{}{},
	}
	for k := range s.data {
		if !other.Contains(k) {
			result.Add(k)
		}
	}
	return result
}

// Intersection compares the set against another one, returning a new set
// containing only the elements that are in both.
func (s *Set) Intersection(other *Set) *Set {
	if other == nil {
		return nil
	}
	result := &Set{
		data: map[string]struct{}{},
	}
	for k := range s.data {
		if other.Contains(k) {
			result.Add(k)
		}
	}
	return result
}

// Union returns a new set containing the elements in the set and in the other
// set.
func (s *Set) Union(other *Set) *Set {
	if other == nil {
		return nil
	}
	result := &Set{
		data: map[string]struct{}{},
	}
	for k := range s.data {
		result.Add(k)
	}
	for k := range other.data {
		result.Add(k)
	}
	return result
}

// SymmetricDifference compares the set agains the other one and returns a new
// set containing only the elements that are not in common.
func (s *Set) SymmetricDifference(other *Set) *Set {
	if other == nil {
		return nil
	}
	result := &Set{
		data: map[string]struct{}{},
	}
	for k := range s.data {
		if !other.Contains(k) {
			result.Add(k)
		}
	}
	for k := range other.data {
		if !s.Contains(k) {
			result.Add(k)
		}
	}
	return result
}

// List returns the elements in the set as a slice of strings.
func (s *Set) List() []string {
	sorted := make([]string, 0, s.Len())
	for k := range s.data {
		sorted = append(sorted, k)
	}
	sort.Strings(sorted)
	return sorted
}

// String returns a string representation of the set.
func (s *Set) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("{ ")
	sorted := make([]string, 0, s.Len())
	for k := range s.data {
		sorted = append(sorted, k)
	}
	sort.Strings(sorted)
	for _, k := range sorted {
		buffer.WriteString(fmt.Sprintf("%q, ", k))
	}
	buffer.WriteRune('}')
	return buffer.String()
}

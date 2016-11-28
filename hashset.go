package colog

import (
	"bytes"
	"encoding/json"
	"sort"
)

// Hash is the base58 string representation of a multihash
type Hash string

// String returns a Hash's string representation
func (h Hash) String() string {
	if h == "" {
		return "null"
	}
	return string(h)
}

// HashSet is a set of Hashes
type HashSet map[Hash]struct{}

// NewHashSet returns a new empty HashSet.
func NewHashSet() HashSet {
	return make(HashSet)
}

// Drop removes Hash h from the set
func (s HashSet) Drop(h Hash) {
	delete(s, h)
}

// Add adds Hash h to the set
func (s HashSet) Add(h Hash) {
	s[h] = struct{}{}
}

// Contains returns whether Hash h is in the set
func (s HashSet) Contains(h Hash) bool {
	_, ok := s[h]
	return ok
}

// Count returns the number of elements in s
func (s HashSet) Count() int {
	return len(s)
}

// Copy returns a copy of the hash set
func (s HashSet) Copy() HashSet {
	s_ := NewHashSet()

	for k := range s {
		s_.Add(k)
	}

	return s_
}

// Sorted returns an alphabetically ordered slice of the hashes in the set
func (s HashSet) Sorted() []Hash {
	hs := make([]Hash, 0, s.Count())
	ss := make([]string, 0, s.Count())

	for h := range s {
		ss = append(ss, string(h))
	}

	// TODO find out if this sorts good enough
	sort.Strings(ss)

	for _, s := range ss {
		hs = append(hs, Hash(s))
	}

	return hs
}

// MarshalJSON returns the JSON representation of the set
func (s HashSet) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Sorted())
}

// UnmarshalJSON adds hashes from JSON (usually called on empty sets)
func (s HashSet) UnmarshalJSON(in []byte) error {
	var hs []Hash

	err := json.Unmarshal(in, &hs)
	if err != nil {
		return err
	}

	for _, h := range hs {
		s.Add(h)
	}

	return nil
}

// String returns the string representation of the set
func (s HashSet) String() string {
	hs := s.Sorted()

	buf := bytes.NewBufferString("{ ")

	if s.Count() > 0 {
		buf.WriteString(hs[0].String())

		hs = hs[1:]

		for _, h := range hs {
			buf.WriteString(", " + h.String())
		}

	}

	buf.WriteString(" }")

	return buf.String()
}

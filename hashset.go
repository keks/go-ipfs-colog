package colog

import (
	"sort"
)

type HashSet map[Hash]struct{}

func (s HashSet) Unset(h Hash) {
	delete(s, h)
}

func (s HashSet) Set(h Hash) {
	s[h] = struct{}{}
}

func (s HashSet) IsSet(h Hash) bool {
	_, ok := s[h]
	return ok
}

func (s HashSet) Copy() HashSet {
	s_ := make(HashSet)

	for k, v := range s {
		s_[k] = v
	}

	return s_
}

func (s HashSet) Sorted() []Hash {
	hs := make([]Hash, 0, len(s))
	ss := make([]string, 0, len(s))

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

package colog

import (
	"bytes"
	"encoding/json"
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

func (s HashSet) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Sorted())
}

func (s HashSet) UnmarshalJSON(in []byte) error {
	var hs []Hash

	err := json.Unmarshal(in, &hs)
	if err != nil {
		return err
	}

	for _, h := range hs {
		s.Set(h)
	}

	return nil
}

func (s HashSet) String() string {

	hs := s.Sorted()

	buf := bytes.NewBufferString("{ ")

	if len(hs) > 0 {
		buf.WriteString(hs[0].String())

		hs = hs[1:]

		for _, h := range hs {
			buf.WriteString(", " + h.String())
		}

	}

	buf.WriteString(" }")

	return buf.String()
}

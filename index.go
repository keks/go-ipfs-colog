package colog

// Index maps a hash to a set of hashes.
type Index map[Hash]HashSet

// Add adds hash g to the set of hashes stored for hash h.
func (i Index) Add(h, g Hash) {
	s, ok := i[h]
	if !ok {
		s = NewHashSet()
		i[h] = s
	}

	s.Add(g)
}

// Get returns the set of hashes stored for h.
func (i Index) Get(h Hash) HashSet {
	return i[h]
}

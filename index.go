package colog

type Index map[Hash]HashSet

func (i Index) Add(h, g Hash) {
	s, ok := i[h]
	if !ok {
		s = NewHashSet()
		i[h] = s
	}

	s.Add(g)
}

func (i Index) Get(h Hash) HashSet {
	return i[h]
}

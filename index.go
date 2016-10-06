package colog

type Index map[Hash]HashSet

func (i Index) Add(h, g Hash) {
	var s HashSet

	s, ok := i[h]
	if !ok {
		s = make(HashSet)
		i[h] = s
	}

	s.Set(g)
}

func (i Index) Get(h Hash) HashSet {
	return i[h]
}

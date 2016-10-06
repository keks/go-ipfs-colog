package colog

type Entry struct {
	Hash
	Value []byte
	Prev  HashSet
}

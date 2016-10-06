package colog

import (
	"errors"
	"github.com/keks/go-ipfs-colog/immutabledb"
	"log"
)

type Hash string

type CoLog struct {
	Id string
	db immutabledb.ImmutableDB

	next, prev Index

	heads HashSet
}

func New(id string, db immutabledb.ImmutableDB) *CoLog {
	// TODO: iterate over db to build index
	// TODO: make index persistent

	return &CoLog{
		Id: id,
		db: db,

		next: make(Index),
		prev: make(Index),

		heads: make(HashSet),
	}
}

func (l *CoLog) Add(data []byte) *Entry {
	hash := Hash(l.db.Put(data))

	e := &Entry{
		Hash:  Hash(hash),
		Value: data,
		Prev:  l.heads.Copy(),
	}

	l.heads = make(HashSet)
	l.heads.Set(hash)

	if len(e.Prev) == 0 {
		e.Prev.Set("")
	}

	for h := range e.Prev {
		l.next.Add(h, e.Hash)
		l.prev.Add(e.Hash, h)
	}

	return e
}

func (l *CoLog) EntryFromHash(h Hash) Entry {
	data := l.db.Get(string(h))

	return Entry{
		Hash:  h,
		Value: data,
		Prev:  l.prev[h],
	}
}

var CorruptLogErr = errors.New("other log corrupt, hash doesn't match")

func (l *CoLog) Contains(h Hash) bool {
	_, ok := l.prev[h]
	return ok
}

func (l *CoLog) Join(other *CoLog) error {
	newHeads := make(HashSet)

	for h := range other.prev {
		// skip known hashes
		if _, ok := l.prev[h]; ok {
			continue
		}

		e := other.EntryFromHash(h)

		// add to db
		h_ := Hash(l.db.Put(e.Value))

		// check if hash matches
		if h != h_ {
			return CorruptLogErr
		}

		// fix up index
		for hPrev := range e.Prev {
			l.next.Add(hPrev, e.Hash)
			l.prev.Add(e.Hash, hPrev)
		}

		// fix heads
		for head := range l.heads {

			// case 1: hash is head in both logs => remains head
			if other.heads.IsSet(head) {
				newHeads.Set(head)
				continue
			}

			// case 2: hash is head in l1, but not in l2 and is not part of l2
			//  => remains head
			if !other.Contains(head) {
				newHeads.Set(head)
				continue
			}

			// case 3: hash is head in l1, but not in l2 and is part of l2
			//  => not head anymore
			// do nothing
		}

		for head := range other.heads {
			// we had those already
			if l.heads.IsSet(head) {
				continue
			}

			// case 4: hash is head in l2, but not in l1 and is not part of l1
			//  => remains head
			if !l.Contains(head) {
				newHeads.Set(head)
				continue
			}

			// case 5: hash is head in l2, but not in l1 and is part of l1
			//  => not head anymore
			// do nothing
		}
	}

	l.heads = newHeads

	return nil
}

// Items returns the Entries in canonical order
func (l *CoLog) Items() []Entry {
	out := make([]Entry, 0, len(l.prev))

	// set up hash stack to track concurrentness
	var stack []Hash

	push := func(hs []Hash) {
		stack = append(stack, hs...)
	}
	pop := func() Hash {
		h := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		return h
	}

	// start with root nodes: nodes with prev=[""]
	stack = l.next[""].Sorted()

	for len(stack) > 0 {
		h := pop()
		e := l.EntryFromHash(h)
		out = append(out, e)

		push(l.next[h].Sorted())
	}

	return out
}

func (l *CoLog) Print() {
	for _, e := range l.Items() {
		log.Println("Entry:", e.Hash)

		data := l.db.Get(string(e.Hash))

		log.Println("Data:", string(data))

		if len(e.Prev) > 0 {
			for p := range e.Prev {
				log.Println("Prev:", p)
			}
		}

		log.Println()
	}
}

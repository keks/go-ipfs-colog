package colog

import (
	"encoding/json"
	"sync"

	"github.com/keks/go-ipfs-colog/immutabledb"
)

// CoLog is a concurrent log
type CoLog struct {
	mutex sync.Mutex

	db immutabledb.ImmutableDB

	next, prev Index

	heads HashSet

	chans *entryChanSet
}

// New returns a concurrent log
func New(db immutabledb.ImmutableDB) *CoLog {
	// TODO: make index persistent

	return &CoLog{
		db: db,

		next: make(Index),
		prev: make(Index),

		heads: NewHashSet(),
		chans: newEntryChanSet(),
	}
}

// Add adds data to the colog and returns the resulting entry
func (l *CoLog) Add(data interface{}) (*Entry, error) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	// prepare entry
	e := &Entry{
		Prev: l.heads.Copy(),
	}

	// set value
	err := e.set(data)
	if err != nil {
		return nil, err
	}

	// use empty string to mark root entry
	if e.Prev.Count() == 0 {
		e.Prev.Add("")
	}

	eBytes, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}

	// store entry
	hStr, err := l.db.Put(eBytes)
	if err != nil {
		return nil, err
	}

	// set hash
	e.Hash = Hash(hStr)

	// update index
	for h := range e.Prev {
		l.next.Add(h, e.Hash)
		l.prev.Add(e.Hash, h)
	}

	// update local heads
	l.heads = NewHashSet()
	l.heads.Add(e.Hash)

	return e, err
}

// Get returns an entry from the db.
func (l *CoLog) Get(h Hash) (*Entry, error) {
	e := NewEntry()

	eBytes, err := l.db.Get(string(h))
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(eBytes, e)
	if err != nil {
		return nil, err
	}

	e.Hash = h

	return e, nil
}

// Contains returns whether an Entry with Hash h is stored
func (l *CoLog) Contains(h Hash) bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	return l.contains(h)
}

func (l *CoLog) contains(h Hash) bool {
	hs, ok := l.prev[h]

	delete(hs, "")

	return ok && len(h) > 0
}

// Join merges colog `other' into `l'
func (l *CoLog) Join(other *CoLog) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	newHeads := make(HashSet)

	for h := range other.prev {
		// skip known hashes
		if l.contains(h) {
			continue
		}

		e, err := other.Get(h)
		if err != nil {
			return err
		}

		eBytes, err := json.Marshal(e)
		if err != nil {
			return err
		}

		// add to db
		_, err = l.db.Put(eBytes)
		if err != nil {
			return err
		}

		// fix heads
		for head := range l.heads {

			// case 1: hash is head in both logs => remains head
			if other.heads.Contains(head) {
				newHeads.Add(head)
				continue
			}

			// case 2: hash is head in l1, but not in l2 and is not part of l2
			//  => remains head
			if !other.contains(head) {
				newHeads.Add(head)
				continue
			}

			// case 3: hash is head in l1, but not in l2 and is part of l2
			//  => not head anymore
			// do nothing
		}

		for head := range other.heads {
			// we had those already
			if l.heads.Contains(head) {
				continue
			}

			// case 4: hash is head in l2, but not in l1 and is not part of l1
			//  => remains head
			if !l.contains(head) {
				newHeads.Add(head)
				continue
			}

			// case 5: hash is head in l2, but not in l1 and is part of l1
			//  => not head anymore
			// do nothing
		}

		// fix up index
		for hPrev := range e.Prev {
			l.next.Add(hPrev, e.Hash)
			l.prev.Add(e.Hash, hPrev)
		}

	}

	l.heads = newHeads

	return nil
}

// Items returns the Entries in canonical order
func (l *CoLog) Items() []*Entry {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	// output Entry slice
	out := make([]*Entry, 0, len(l.prev))

	// keeps track how many of the revious pointers have been added to out already
	addedPrevs := map[Hash]int{}

	// set up hash stack to track concurrentness
	var stack = []Hash{}

	// push to stack
	push := func(hs ...Hash) {
		stack = append(stack, hs...)
	}

	//pop from stack
	pop := func() Hash {
		h := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		return h
	}

	// start with root nodes: nodes with prev=[""]
	push(l.next[""].Sorted()...)

	for len(stack) > 0 {
		// pop hash from stack
		h := pop()

		//ignore root
		if h == "" {
			continue
		}

		// get Entry
		e, err := l.Get(h)
		if err != nil {
			continue
		}

		// append Entry
		out = append(out, e)

		// mark that an Entry was added in all next hashes
		for hNext := range l.next[h] {
			addedPrevs[hNext]++
		}

		// push next hashes, but only if all past hashes have been added
		for _, hNext := range l.next[h].Sorted() {
			if addedPrevs[hNext] == len(l.prev[hNext]) {
				push(hNext)
			}
		}
	}

	return out
}

func (l *CoLog) Heads() []Hash {
	return l.heads.Sorted()
}

func (l *CoLog) FetchFromHead(head Hash) error {
	if _, ok := l.prev[head]; ok {
		return nil
	}

	// set up hash stack to track concurrentness
	var stack = []Hash{}

	// push to stack
	push := func(hs ...Hash) {
		stack = append(stack, hs...)
	}

	//pop from stack
	pop := func() Hash {
		h := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		return h
	}

	// walk colog from head
	push(head)

	for len(stack) > 0 {
		// pop hash from stack
		h := pop()

		//ignore root
		if h == "" {
			continue
		}

		// check if already known
		if _, ok := l.prev[h]; ok {
			continue
		}

		// set as head if no followups known
		if nexts := l.next[h]; len(nexts) == 0 {
			l.heads.Add(h)
		}

		// get Entry
		e, err := l.Get(h)
		if err != nil {
			continue
		}

		for hPrev := range e.Prev {
			// remove from heads in case it was there
			l.heads.Drop(hPrev)

			// add to index
			l.prev.Add(h, hPrev)
			l.next.Add(hPrev, h)

			// mark for recursion
			push(hPrev)
		}
	}

	return nil
}

func (l *CoLog) Watch() <-chan *Entry {
	ch := make(chan *Entry)

	l.chans.Add(ch)
	return ch
}

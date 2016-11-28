package colog

import (
	"encoding/json"
	"sync"
)

// Entry is an element of the log.
type Entry struct {
	Hash  Hash            `json:"-"`
	Value json.RawMessage `json:"payload"`
	Prev  HashSet         `json:"next"`
}

// NewEntry returns a new Entry.
func NewEntry() *Entry {
	return &Entry{Prev: NewHashSet()}
}

func (e *Entry) set(v interface{}) (err error) {
	if e.Prev == nil {
		e.Prev = NewHashSet()
	}

	e.Value, err = json.Marshal(v)
	return err
}

// Get writes the value that is serialized in the entry into v. v needs to be a pointer.
func (e *Entry) Get(v interface{}) (err error) {
	return json.Unmarshal(e.Value, v)
}

// GetString returns the string stored in the Entry. If the value stored is not a string, an empty string is returned.
func (e *Entry) GetString() string {
	var s string

	e.Get(&s)
	return s
}

// String returns the string representation of the Entry.
func (e *Entry) String() string {
	return "{ " + e.Hash.String() + ": " + string(e.Value) + " " + e.Prev.String() + " }"
}

// set of entry channels
type entryChanSet struct {
	sync.Mutex
	chans  map[chan<- *Entry]struct{}
	rchans map[<-chan *Entry]chan<- *Entry
}

func newEntryChanSet() *entryChanSet {
	return &entryChanSet{
		chans:  make(map[chan<- *Entry]struct{}),
		rchans: make(map[<-chan *Entry]chan<- *Entry),
	}
}

func (cs *entryChanSet) New() <-chan *Entry {
	ch := make(chan *Entry)

	cs.Lock()
	cs.chans[ch] = struct{}{}
	cs.rchans[ch] = ch
	cs.Unlock()

	return ch
}

func (cs *entryChanSet) Drop(ch <-chan *Entry) {
	cs.Lock()
	delete(cs.chans, cs.rchans[ch])
	delete(cs.rchans, ch)
	cs.Unlock()
}

func (cs *entryChanSet) Send(e *Entry) {
	cs.Lock()

	for ch := range cs.chans {
		go func(ch_ chan<- *Entry) {
			ch_ <- e
		}(ch)
	}

	cs.Unlock()
}

func (cs *entryChanSet) Count() int {
	return len(cs.chans)
}

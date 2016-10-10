package colog

import (
	"encoding/json"
	"sync"
)

type Entry struct {
	Hash  Hash            `json:"-"`
	Value json.RawMessage `json:"payload"`
	Prev  HashSet         `json:"next"`
}

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

func (e *Entry) Get(v interface{}) (err error) {
	return json.Unmarshal(e.Value, v)
}

func (e *Entry) GetString() string {
	var s string

	e.Get(&s)
	return s
}

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
		go func() {
			ch <- e
		}()
	}

	cs.Unlock()
}

func (cs *entryChanSet) Count() int {
	return len(cs.chans)
}

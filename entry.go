package colog

import "encoding/json"

type Entry struct {
	Hash  Hash            `json:"-"`
	Value json.RawMessage `json:"payload"`
	Prev  HashSet         `json:"next"`
}

func NewEntry() *Entry {
	return &Entry{Prev: make(HashSet)}
}

func (e *Entry) set(v interface{}) (err error) {
	if e.Prev == nil {
		e.Prev = make(HashSet)
	}

	e.Value, err = json.Marshal(v)
	return err
}

func (e *Entry) Get(v interface{}) (err error) {
	return json.Unmarshal(e.Value, v)
}

func (e *Entry) GetString() string {
	var s string

	json.Unmarshal(e.Value, &s)
	return s
}

func (e *Entry) String() string {
	return "{ " + e.Hash.String() + ": " + string(e.Value) + " " + e.Prev.String() + " }"
}

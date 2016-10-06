package appendonlylog

import "container/list"

type Entry struct {
	list.Element
	Key   string
	Value []byte
	Next  []*Entry
}

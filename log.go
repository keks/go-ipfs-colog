package appendonlylog

import (
  "log"
  "container/list"
  "github.com/haadcode/go-ipfs-log/immutabledb/interface"
)

type AppendOnlyLog struct {
  Id    string
  items *list.List
  db immutabledb.ImmutableDB
}

func New(id string, db immutabledb.ImmutableDB) *AppendOnlyLog {
  return new(AppendOnlyLog).Init(id, db)
}

func (l *AppendOnlyLog) Init(id string, db immutabledb.ImmutableDB) *AppendOnlyLog {
  l.Id = id
  l.db = db
  l.items = list.New()
  return l
}

func (l *AppendOnlyLog) Add(data []byte) *Entry {
  hash := l.db.Put(data)

  var next []*Entry
  n := l.items.Back()

  if n != nil {
    next = []*Entry{n.Value.(*Entry)}
  }

  e := &Entry{
    Key: hash,
    Value: data,
    Next: next,
  }

  l.items.PushBack(e)

  return e
}

func (l *AppendOnlyLog) Join(other *AppendOnlyLog) *AppendOnlyLog {
  var res = New(l.Id, l.db)

  items := l.Items()
  s := other.Items()

  // TODO:
  // 1) add _currentBatch for tracking "local entries while offline"
  // 2) follow JS impl:
  //    const diff   = differenceWith(other.items, this.items, Entry.equals)
  //    // TODO: need deterministic sorting of the union
  //    const final  = unionWith(this._currentBatch, diff, Entry.equals)
  //    this._items  = this._items.concat(final)
  // 3) fetch history (ie. entries that are not in local log)
  // 4) add support for multiple heads, implement find and update heads logic

  for i := len(items) - 1; i >= 0; i-- {
    res.items.PushBack(items[i])
  }

  for i := len(s) - 1; i >= 0; i-- {
    if (!contains(items, s[i])) {
      res.items.PushBack(s[i])
    }
  }

  return res
}

func contains(s []*Entry, e *Entry) bool {
  for _, a := range s {
    if a.Key == e.Key {
      return true
    }
  }
  return false
}

func (l *AppendOnlyLog) Head() *Entry {
  head := l.items.Back().Value.(*Entry)
  return head
}

func (l *AppendOnlyLog) Items() []*Entry {
  var more []*Entry

  for e := l.items.Front(); e != nil; e = e.Next() {
    entry := e.Value.(*Entry)
    more = append(more, entry)
  }

  return more
}

func (l *AppendOnlyLog) Print() {
  for _, e := range l.Items() {
    log.Println("Entry:", e.Key)

    data := l.db.Get(e.Key)

    log.Println("Data:", string(data))

    if (len(e.Next) > 0) {
      for _, n := range e.Next {
        log.Println("Next:", n.Key)
      }
    }

    log.Println()
  }
}

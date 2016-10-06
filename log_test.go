package colog

import (
	"bytes"
	"fmt"
	db "github.com/keks/go-ipfs-colog/immutabledb/ipfs"
	"strings"
	"testing"
)

var dataDirectory = "/tmp/go-ipfs-colog-dev"
var ipfsdb = db.Open(dataDirectory)

var id = "abc"
var value1 = []byte("Hello1")
var value2 = []byte("Hello2")
var value3 = []byte("Hello3")
var hash1 = "QmX96xhp6cUB1YE5nqZsmKHbZFiAEderPc3gapGdwAoEod"

/* Create */

func TestNew(t *testing.T) {
	var log1 = New(id, ipfsdb)

	if log1 == nil {
		t.Fatalf("Couldn't create a log")
	}

	if log1.Id != id {
		t.Fatalf("Id not set")
	}

	if log1.db == nil {
		t.Fatalf("DB not set")
	}

	if len(log1.Items()) != 0 {
		t.Fatalf("Items not empty")
	}
}

/* Add */

func TestAdd(t *testing.T) {
	var log1 = New(id, ipfsdb)

	one := log1.Add(value1)

	if one == nil {
		t.Fatal("Entry was not added")
	}

	if strings.Compare(string(one.Hash), hash1) != 0 {
		t.Fatalf("Wrong key: %s", one.Hash)
	}

	if bytes.Compare(one.Value, value1) != 0 {
		t.Fatalf("Wrong key: %s", one.Hash)
	}

	if len(one.Prev) != 1 && one.Prev.Sorted()[0] == "" {
		t.Fatalf("Wrong next reference: %s", one.Prev)
	}

	t.Logf("%#v\n", one)

	if len(log1.Items()) != 1 {
		t.Fatalf("Wrong items count: %d", len(log1.Items()))
	}
}

func ExampleAdd_one() {
	var log1 = New(id, ipfsdb)

	one := log1.Add(value1)

	fmt.Println(one.Hash)
	fmt.Println(string(one.Value))
	fmt.Println(one.Prev)
	// Output:
	// QmX96xhp6cUB1YE5nqZsmKHbZFiAEderPc3gapGdwAoEod
	// Hello1
	// map[:{}]
}

func ExampleAdd_two() {
	var log1 = New(id, ipfsdb)

	log1.Add(value1)
	log1.Add(value2)

	items := log1.Items()
	fmt.Println(len(items))
	fmt.Println(items[1].Hash)
	fmt.Println(string(items[1].Value))
	fmt.Println(log1.EntryFromHash(items[1].Prev.Sorted()[0]).Hash)
	// Output:
	// 2
	// Qme39B2h1QTDYAwCa4gXa6DB6R3TAFaG2Z8HF48U1wkKE6
	// Hello2
	// QmX96xhp6cUB1YE5nqZsmKHbZFiAEderPc3gapGdwAoEod
}

func ExampleAdd_three() {
	var log1 = New(id, ipfsdb)

	log1.Add(value1)
	log1.Add(value2)
	log1.Add(value3)

	items := log1.Items()
	fmt.Println(len(items))
	fmt.Println(string(items[0].Value))
	fmt.Println(string(items[1].Value))
	fmt.Println(string(items[2].Value))
	fmt.Println(items[0].Prev)
	fmt.Println(log1.EntryFromHash(items[1].Prev.Sorted()[0]).Hash)
	fmt.Println(log1.EntryFromHash(items[2].Prev.Sorted()[0]).Hash)
	// Output:
	// 3
	// Hello1
	// Hello2
	// Hello3
	// map[:{}]
	// QmX96xhp6cUB1YE5nqZsmKHbZFiAEderPc3gapGdwAoEod
	// Qme39B2h1QTDYAwCa4gXa6DB6R3TAFaG2Z8HF48U1wkKE6
}

func BenchmarkAdd(b *testing.B) {
	var log1 = New(id, ipfsdb)

	for i := 0; i < b.N; i++ {
		log1.Add(value1)
	}
}

/* Join */

func TestJoin(t *testing.T) {
	var log1 = New(id, ipfsdb)
	var log2 = New(id, ipfsdb)

	log1.Add(value1)
	log2.Add(value2)

	log1.Join(log2)
	items := log1.Items()

	if len(items) != 2 {
		t.Fatalf("Wrong number of entries: %i", len(items))
	}

	// Make sure the joined log doesn't have pointers to the joined logs
	log1.Add(value1)
	log2.Add(value2)

}

func ExampleJoin_one() {
	var log1 = New(id, ipfsdb)
	var log2 = New(id, ipfsdb)

	log1.Add(value1)
	log2.Add(value2)

	log1.Join(log2)

	items := log1.Items()
	first := items[0]
	second := items[1]

	fmt.Println(len(items))
	fmt.Println(first.Hash)
	fmt.Println(second.Hash)
	fmt.Println(string(first.Value))
	fmt.Println(string(second.Value))
	// Output:
	// 2
	// Qme39B2h1QTDYAwCa4gXa6DB6R3TAFaG2Z8HF48U1wkKE6
	// QmX96xhp6cUB1YE5nqZsmKHbZFiAEderPc3gapGdwAoEod
	// Hello2
	// Hello1
}

func BenchmarkJoin(b *testing.B) {
	var log1 = New(id, ipfsdb)
	var log2 = New(id, ipfsdb)

	log1.Add(value1)
	log2.Add(value2)

	for i := 0; i < b.N; i++ {
		log1.Join(log2)
	}
}

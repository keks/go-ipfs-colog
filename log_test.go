package appendonlylog

import (
	"bytes"
	"fmt"
	db "github.com/haadcode/go-ipfs-log/immutabledb/ipfs"
	"strings"
	"testing"
)

var dataDirectory = "/tmp/go-ipfs-log-dev"
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
		t.Errorf("Couldn't create a log")
	}

	if log1.Id != id {
		t.Errorf("Id not set")
	}

	if log1.db == nil {
		t.Errorf("DB not set")
	}

	if len(log1.Items()) != 0 {
		t.Errorf("Items not empty")
	}
}

/* Add */

func TestAdd(t *testing.T) {
	var log1 = New(id, ipfsdb)

	one := log1.Add(value1)

	if one == nil {
		t.Errorf("Entry was not added")
	}

	if strings.Compare(one.Key, hash1) != 0 {
		t.Errorf("Wrong key: %s", one.Key)
	}

	if bytes.Compare(one.Value, value1) != 0 {
		t.Errorf("Wrong key: %s", one.Key)
	}

	if len(one.Next) != 0 {
		t.Errorf("Wrong next reference: %s", one.Next)
	}

	if len(log1.Items()) != 1 {
		t.Errorf("Wrong items count: %i", len(log1.Items()))
	}
}

func ExampleAdd_one() {
	var log1 = New(id, ipfsdb)

	one := log1.Add(value1)

	fmt.Println(one.Key)
	fmt.Println(string(one.Value))
	fmt.Println(one.Next)
	// Output:
	// QmX96xhp6cUB1YE5nqZsmKHbZFiAEderPc3gapGdwAoEod
	// Hello1
	// []
}

func ExampleAdd_two() {
	var log1 = New(id, ipfsdb)

	log1.Add(value1)
	log1.Add(value2)

	items := log1.Items()
	fmt.Println(len(items))
	fmt.Println(items[1].Key)
	fmt.Println(string(items[1].Value))
	fmt.Println(items[1].Next[0].Key)
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
	fmt.Println(items[0].Next)
	fmt.Println(items[1].Next[0].Key)
	fmt.Println(items[2].Next[0].Key)
	// Output:
	// 3
	// Hello1
	// Hello2
	// Hello3
	// []
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

	log3 := log1.Join(log2)
	items := log3.Items()
	first := items[0]
	second := items[1]

	if len(items) != 2 {
		t.Errorf("Wrong number of entries: %i", len(items))
	}

	if log3.Id != log1.Id {
		t.Errorf("Wrong id: %s", log3.Id)
	}

	if bytes.Compare(first.Value, value1) != 0 {
		t.Errorf("Wrong value: %s", string(first.Value))
	}

	if bytes.Compare(second.Value, value2) != 0 {
		t.Errorf("Wrong value: %s", string(second.Value))
	}

	// Make sure the joined log doesn't have pointers to the joined logs
	log1.Add(value1)
	log2.Add(value2)

	if len(log3.Items()) != 2 {
		t.Errorf("Wrong number of entries: %i", len(log3.Items()))
	}
}

func ExampleJoin_one() {
	var log1 = New(id, ipfsdb)
	var log2 = New(id, ipfsdb)

	log1.Add(value1)
	log2.Add(value2)

	log3 := log1.Join(log2)

	items := log3.Items()
	first := items[0]
	second := items[1]

	fmt.Println(len(items))
	fmt.Println(first.Key)
	fmt.Println(second.Key)
	fmt.Println(string(first.Value))
	fmt.Println(string(second.Value))
	// Output:
	// 2
	// QmX96xhp6cUB1YE5nqZsmKHbZFiAEderPc3gapGdwAoEod
	// Qme39B2h1QTDYAwCa4gXa6DB6R3TAFaG2Z8HF48U1wkKE6
	// Hello1
	// Hello2
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

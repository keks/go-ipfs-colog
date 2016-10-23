package colog

import (
	"fmt"
	db "github.com/keks/go-ipfs-colog/immutabledb/ipfs-api"
	"strings"
	"testing"
)

var ipfsdb = db.New()

var value1 = "Hello1"
var value2 = "Hello2"
var value3 = "Hello3"
var hash1 = "QmbiruS6UMT6gT3JBHtZNWKitEssyUQzYu8k4gGd6rhzNc"

/* Create */

func TestNew(t *testing.T) {
	var log1 = New(ipfsdb)

	if log1 == nil {
		t.Fatalf("Couldn't create a log")
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
	var log1 = New(ipfsdb)

	one, err := log1.Add(value1)
	if err != nil {
		t.Fatal(err)
	}

	if one == nil {
		t.Fatal("Entry was not added")
	}

	if strings.Compare(string(one.Hash), hash1) != 0 {
		t.Fatalf("Wrong key: %s", one.Hash)
	}

	if strings.Compare(one.GetString(), value1) != 0 {
		t.Fatalf("Wrong value: %s", one.GetString())
	}

	if len(one.Prev) != 1 && one.Prev.Sorted()[0] == "" {
		t.Fatalf("Wrong next reference: %s", one.Prev)
	}

	if len(log1.Items()) != 1 {
		t.Fatalf("Wrong items count: %d", len(log1.Items()))
	}
}

func ExampleAdd_one() {
	var log1 = New(ipfsdb)

	one, err := log1.Add(value1)
	if err != nil {
		panic(err)
	}

	fmt.Println(one.Hash)
	fmt.Println(one.GetString())
	fmt.Println(one.Prev)
	// Output:
	// QmbiruS6UMT6gT3JBHtZNWKitEssyUQzYu8k4gGd6rhzNc
	// Hello1
	// { null }
}

func ExampleAdd_two() {
	var log1 = New(ipfsdb)

	log1.Add(value1)
	log1.Add(value2)

	items := log1.Items()
	fmt.Println(len(items))
	fmt.Println(items[1].Hash)
	fmt.Println(string(items[1].GetString()))
	e, err := log1.Get(items[1].Prev.Sorted()[0])
	if err != nil {
		panic(err)
	}

	fmt.Println(e.Hash)
	// Output:
	// 2
	// QmdezppSoeGZmyYQZEuiUTU4cxTqyVfiwBct7D5iosY6zN
	// Hello2
	// QmbiruS6UMT6gT3JBHtZNWKitEssyUQzYu8k4gGd6rhzNc
}

func ExampleAdd_three() {
	var log1 = New(ipfsdb)

	log1.Add(value1)
	log1.Add(value2)
	log1.Add(value3)

	items := log1.Items()
	fmt.Println(len(items))
	fmt.Println(string(items[0].GetString()))
	fmt.Println(string(items[1].GetString()))
	fmt.Println(string(items[2].GetString()))
	fmt.Println(items[0].Prev)
	e1, err := log1.Get(items[1].Prev.Sorted()[0])
	if err != nil {
		panic(err)
	}
	fmt.Println(e1.Hash)
	e2, err := log1.Get(items[2].Prev.Sorted()[0])
	if err != nil {
		panic(err)
	}
	fmt.Println(e2.Hash)
	// Output:
	// 3
	// Hello1
	// Hello2
	// Hello3
	// { null }
	// QmbiruS6UMT6gT3JBHtZNWKitEssyUQzYu8k4gGd6rhzNc
	// QmdezppSoeGZmyYQZEuiUTU4cxTqyVfiwBct7D5iosY6zN
}

func BenchmarkAdd(b *testing.B) {
	var log1 = New(ipfsdb)

	for i := 0; i < b.N; i++ {
		log1.Add(value1)
	}
}

/* Join */

func TestJoin(t *testing.T) {
	var log1 = New(ipfsdb)
	var log2 = New(ipfsdb)

	log1.Add(value1)
	log2.Add(value2)

	log1.Join(log2)
	items := log1.Items()

	if len(items) != 2 {
		t.Log("items:", items)
		t.Fatalf("Wrong number of entries: %d", len(items))
	}

	// Make sure the joined log doesn't have pointers to the joined logs
	log1.Add(value1)
	log2.Add(value2)

}

func ExampleJoin_one() {
	var log1 = New(ipfsdb)
	var log2 = New(ipfsdb)

	log1.Add(value1)
	log2.Add(value2)

	log1.Join(log2)

	items := log1.Items()
	first := items[0]
	second := items[1]

	fmt.Println(len(items))
	fmt.Println(first.Hash)
	fmt.Println(second.Hash)
	fmt.Println(string(first.GetString()))
	fmt.Println(string(second.GetString()))
	// Output:
	// 2
	// QmbiruS6UMT6gT3JBHtZNWKitEssyUQzYu8k4gGd6rhzNc
	// QmXgHKWQEG6NwpqJMvvvBdMrzWqPNo6cLtZr4BpZHFvrXV
	// Hello1
	// Hello2
}

func BenchmarkJoin(b *testing.B) {
	var log1 = New(ipfsdb)
	var log2 = New(ipfsdb)

	log1.Add(value1)
	log2.Add(value2)

	for i := 0; i < b.N; i++ {
		log1.Join(log2)
	}
}

package colog

import (
	"fmt"
)

func ExampleWatch() {
	l1 := New(ipfsdb)
	l2 := New(ipfsdb)
	l3 := New(ipfsdb)
	l4 := New(ipfsdb)

	ch := l1.Watch()
	wait := make(chan struct{})

	go func() {
		for e := range ch {
			fmt.Println(e.GetString())
			wait <- struct{}{}
		}
	}()

	l1.Add("abcdef")

	<-wait

	l2.Join(l1)
	l2.Add("ghijk")
	l1.Join(l2)

	<-wait

	l3.Join(l1)
	e3 := checkadd(l3.Add("lmnop"))

	l1.FetchFromHead(e3.Hash)

	<-wait

	e4 := checkadd(l4.Add("qrstu"))

	l1.FetchFromHead(e4.Hash)

	<-wait

	//Output:
	//abcdef
	//ghijk
	//lmnop
	//qrstu
}

func ExampleWatch_close() {
	l1 := New(ipfsdb)
	l2 := New(ipfsdb)

	ch := l1.Watch()
	wait := make(chan struct{})

	go func() {
		for e := range ch {
			fmt.Println(e.GetString())
			wait <- struct{}{}
		}
	}()

	l1.Add("abcdef")

	<-wait

	ch2 := l1.Watch()
	l1.Unwatch(ch)

	l2.Join(l1)
	l2.Add("ghijk")
	l1.Join(l2)

	<-ch2

	fmt.Println(l1.chans.Count())

	//Output:
	// abcdef
	// 1
}

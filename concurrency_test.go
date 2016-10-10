package colog

import (
	"encoding/json"
	"fmt"
	"sync"
)

func ExampleConcurrentAccess() {
	var wg sync.WaitGroup

	l := New(ipfsdb)

	n := 50

	sum := 0

	wg.Add(n)

	for i := 0; i < n; i++ {
		go func(j int) {
			l.Add(j)
			wg.Done()
		}(i)
	}

	wg.Wait()

	for _, e := range l.Items() {
		var j int
		err := json.Unmarshal(e.Value, &j)
		if err != nil {
			panic(err)
		}

		sum += j
	}

	// 1+2+3+4+...+n   = n*(n+1)/2
	// 0+1+2+3+...+n-1 = (n-1)*n/2
	//								 = n*(n-1)/2
	if sum != n*(n-1)/2 {
		panic("wrong sum")
	}

	fmt.Println("ok")

	//Output:
	//ok
}

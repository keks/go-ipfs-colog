package main

import (
	"github.com/haadcode/go-ipfs-log/appendonlylog"
	db "github.com/haadcode/go-ipfs-log/immutabledb/ipfs"
	"log"
)

var dataDirectory = "/tmp/go-ipfs-log-dev"
var ipfsdb = db.Open(dataDirectory)

var log1 = appendonlylog.New("abc", ipfsdb)
var log2 = appendonlylog.New("def", ipfsdb)

func printLog(l *appendonlylog.AppendOnlyLog) {
	log.Println()
	log.Println("--------------------")
	log.Println("Log Id:", l.Id)
	log.Println("Head:", l.Head().Key)
	log.Println("Items:", len(l.Items()))
	log.Println("--------------------")
	log.Println()
	l.Print()
}

func main() {
	log.Println("-- go-ipfs-log --")
	log.Println()

	one := log1.Add([]byte("Hallo welt!"))
	log.Println("Added one entry:", one.Key)

	two := log2.Add([]byte("Hello world!"))
	log.Println("Added one entry:", two.Key)

	log1.Add([]byte("Data datA"))
	log1.Add([]byte("12345"))

	printLog(log1)
	printLog(log2)

	log2.Add([]byte("12345")) // add double entry to second log, this should not be in log3 twice

	log3 := log1.Join(log2)
	printLog(log3)

	log3.Add([]byte("88"))
	log2.Add([]byte("777"))
	log4 := log3.Join(log2)
	printLog(log4)

	ipfsdb.Close()
}

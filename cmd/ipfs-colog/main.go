package main

import (
	"github.com/keks/go-ipfs-colog"
	db "github.com/keks/go-ipfs-colog/immutabledb/ipfs-api"
	//"gx/ipfs/QmVcLF2CgjQb5BWmYFWsDfxDjbzBfcChfdHRedxeL3dV4K/cli"

	"log"
)

var ipfsdb = db.New()

var log1 = colog.New("abc", ipfsdb)
var log2 = colog.New("def", ipfsdb)

func printLog(l *colog.CoLog) {
	log.Println()
	log.Println("--------------------")
	log.Println("Log Id:", l.Id)
	log.Println("Heads:", l.Heads())
	log.Println("Items:", len(l.Items()))
	log.Println("--------------------")
	log.Println()
	l.Print()
}

func main() {
	log.Println("-- go-ipfs-colog --")
	log.Println()

	one, err := log1.Add("Hallo welt!")
	log.Println("Added one entry:", one.Hash, "err:", err)

	two, err := log2.Add("Hello world!")
	log.Println("Added one entry:", two.Hash, "err:", err)

	log1.Add("Data datA")
	log1.Add("12345")

	printLog(log1)
	printLog(log2)

	log2.Add("12345") // add double entry to second log, this should not be in log3 twice

	log1.Join(log2)
	printLog(log1)

	log1.Add([]byte("88"))
	log2.Add([]byte("777"))
	log1.Join(log2)
	printLog(log1)

	ipfsdb.Close()
}

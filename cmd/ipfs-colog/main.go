package main

import (
	"github.com/keks/go-ipfs-colog"
	db "github.com/keks/go-ipfs-colog/immutabledb/ipfs-api"
	"gx/ipfs/QmVcLF2CgjQb5BWmYFWsDfxDjbzBfcChfdHRedxeL3dV4K/cli"

	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path"
)

var headFilePath string
var l *colog.CoLog

func fail(reasons ...interface{}) {
	fmt.Print("error: ")
	fmt.Println(reasons...)

	os.Exit(-1)
}

func init() {
	headFilePath = path.Join(os.Getenv("HOME"), ".colog-heads")
}

func updateHeadFile(e *colog.Entry) error {
	f, err := os.Create(headFilePath)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.WriteString(string(e.Hash) + "\n")
	return err
}

func prepareLog() (*colog.CoLog, error) {
	ipfsdb := db.New()
	l = colog.New("abc", ipfsdb)

	f, err := os.Open(headFilePath)
	if err != nil { // no state to recover; all good
		log.Println("couldn't recover state; open err:", err)
		return l, nil
	}

	defer f.Close()

	bf := bufio.NewReader(f)

	for {
		line, err := bf.ReadString('\n')
		if len(line) == 0 || err != nil {
			break
		}

		err = l.FetchFromHead(colog.Hash(line))
		if err != nil {
			continue
		}
	}

	return l, nil
}

var putCmd = cli.Command{
	Name:      "add",
	ShortName: "a",
	Usage:     "add a value to the colog",
	Category:  "simple",
	Flags:     []cli.Flag{cli.BoolFlag{Name: "s"}},
	Action: func(c *cli.Context) error {
		var data interface{}

		if c.Bool("s") {
			data = c.Args()[0]
		} else {
			b := bytes.Buffer{}

			io.Copy(&b, os.Stdin)
			data = b.String()

		}
		e, err := l.Add(data)
		if err != nil {
			return err
		}

		fmt.Println(e)
		return updateHeadFile(e)
	},
}

func main() {
	var err error
	l, err = prepareLog()
	if err != nil {
		fail(err)
	}

	app := cli.NewApp()
	app.Name = "ipfs-colog"
	app.Usage = "work with cologs"
	app.Commands = []cli.Command{putCmd}
	app.Run(os.Args)
}

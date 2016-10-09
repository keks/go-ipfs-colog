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

func updateHeadFile(l *colog.CoLog) error {
	f, err := os.Create(headFilePath)
	if err != nil {
		return err
	}

	defer f.Close()

	for _, h := range l.Heads() {
		_, err = f.WriteString(string(h) + "\n")
		if err != nil {
			return err
		}
	}
	return nil
}

var ipfsdb = db.New()

func prepareLog() (*colog.CoLog, error) {
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

		line = line[:len(line)-1]

		err = l.FetchFromHead(colog.Hash(line))
		if err != nil {
			continue
		}
	}

	return l, nil
}

var (
	addCmd = cli.Command{
		Name:      "add",
		ShortName: "a",
		Usage:     "add a value to the colog",
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "r"},
			cli.BoolFlag{Name: "s"},
			cli.BoolFlag{Name: "v"},
		},
		Action: func(c *cli.Context) error {
			var data interface{}

			l_old := l

			if c.Bool("r") {
				l = colog.New("abc", ipfsdb)
			}

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

			if c.Bool("v") {
				fmt.Println("l:")
				for _, e := range l.Items() {
					fmt.Println(e)
				}
				fmt.Println("l_old:")
				for _, e := range l_old.Items() {
					fmt.Println(e)
				}
			}

			if c.Bool("r") {
				l_old.Join(l)
				l = l_old
			}

			if !c.Bool("v") {
				fmt.Println(e)
			}
			return updateHeadFile(l)
		},
	}

	printCmd = cli.Command{
		Name:      "print",
		ShortName: "p",
		Usage:     "print the log",
		Action: func(c *cli.Context) error {
			for _, e := range l.Items() {
				fmt.Println(e)
			}

			return nil
		},
	}

	headsCmd = cli.Command{
		Name:      "heads",
		ShortName: "hs",
		Usage:     "print the heads",
		Action: func(c *cli.Context) error {
			for _, h := range l.Heads() {
				fmt.Println(h)
			}

			return nil
		},
	}
)

func main() {
	var err error
	l, err = prepareLog()
	if err != nil {
		fail(err)
	}

	app := cli.NewApp()
	app.Name = "ipfs-colog"
	app.Usage = "work with cologs"
	app.Commands = []cli.Command{addCmd, printCmd, headsCmd}
	app.Run(os.Args)
}

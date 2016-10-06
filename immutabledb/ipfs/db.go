package immutabledb

import (
	"github.com/keks/go-ipfs-colog/immutabledb"
	"gx/ipfs/QmQs9UguUVkFC3hXsS6MyGu377GzUfZfraddPDazsfmc6t/go-ipfs/core"
	dag "gx/ipfs/QmQs9UguUVkFC3hXsS6MyGu377GzUfZfraddPDazsfmc6t/go-ipfs/merkledag"
	path "gx/ipfs/QmQs9UguUVkFC3hXsS6MyGu377GzUfZfraddPDazsfmc6t/go-ipfs/path"
	repo "gx/ipfs/QmQs9UguUVkFC3hXsS6MyGu377GzUfZfraddPDazsfmc6t/go-ipfs/repo"
	fsrepo "gx/ipfs/QmQs9UguUVkFC3hXsS6MyGu377GzUfZfraddPDazsfmc6t/go-ipfs/repo/fsrepo"

	"context"
	"log"
)

// Trick to make sure ImmutableIPFS implements ImmutableDB
var _ immutabledb.ImmutableDB = ImmutableIPFS{}

type ImmutableIPFS struct {
	Repo repo.Repo
	Node *core.IpfsNode
}

func Open(path string) ImmutableIPFS {
	r, err := fsrepo.Open(path)
	if err != nil {
		log.Fatal("Can't open data repository at %s: %s", path, err)
	}

	cfg := &core.BuildCfg{
		Repo:   r,
		Online: false,
	}

	nd, err := core.NewNode(context.Background(), cfg)
	if err != nil {
		log.Fatal("Can't create IPFS node: %s", err)
	}

	return ImmutableIPFS{
		Repo: r,
		Node: nd,
	}
}

func (db ImmutableIPFS) Close() error {
	db.Repo.Close()
	db.Node.Close()
	return nil
}

func (db ImmutableIPFS) Put(data []byte) string {
	obj := dag.NodeWithData(data)

	k, err := db.Node.DAG.Add(obj)
	if err != nil {
		log.Fatal(err)
	}

	return k.B58String()
}

func (db ImmutableIPFS) Get(key string) []byte {
	ctx := context.Background()
	fpath := path.Path(key)

	object, err := core.Resolve(ctx, db.Node, fpath)
	if err != nil {
		log.Fatal(err)
	}

	node := dag.NodeWithData(object.Data())

	return node.Data()
}

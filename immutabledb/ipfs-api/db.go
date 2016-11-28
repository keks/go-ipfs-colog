package idb

import (
	shell "github.com/ipfs/go-ipfs-api"
)

type DB struct {
	sh *shell.Shell
}

func New() (*DB, error) {
	sh, err := getLocalApiShell()
	return &DB{sh}, err
}

func (db *DB) Put(data []byte) (string, error) {
	return db.sh.ObjectPut(
		&shell.IpfsObject{
			Data: string(data),
		})
}

func (db *DB) Get(addr string) ([]byte, error) {
	o, err := db.sh.ObjectGet(addr)
	if err != nil {
		return nil, err
	}

	return []byte(o.Data), nil
}

func (db *DB) Close() error {
	return nil
}

// the contents of this file are taken from https://github.com/whyrusleeping/gx/blob/master/gxutil/shell.go
// written by @whyrusleeping for the IPFS project
// published under the MIT license as stated at https://github.com/whyrusleeping/gx/blob/master/LICENSE

package idb

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	shell "github.com/ipfs/go-ipfs-api"
	homedir "github.com/mitchellh/go-homedir"
	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr-net"
)

func getLocalApiShell() (*shell.Shell, error) {
	ipath := os.Getenv("IPFS_PATH")
	if ipath == "" {
		home, err := homedir.Dir()
		if err != nil {
			return nil, err
		}

		ipath = filepath.Join(home, ".ipfs")
	}

	apifile := filepath.Join(ipath, "api")

	data, err := ioutil.ReadFile(apifile)
	if err != nil {
		return nil, err
	}

	addr := strings.Trim(string(data), "\n\t ")

	host, err := multiaddrToNormal(addr)
	if err != nil {
		return nil, err
	}

	local := shell.NewShell(host)

	_, _, err = local.Version()
	if err != nil {
		return nil, err
	}

	return local, nil
}

// same as getLocalApiShell
func multiaddrToNormal(addr string) (string, error) {
	maddr, err := ma.NewMultiaddr(addr)
	if err != nil {
		return "", err
	}

	_, host, err := manet.DialArgs(maddr)
	if err != nil {
		return "", err
	}

	return host, nil
}

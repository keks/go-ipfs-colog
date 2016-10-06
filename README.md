# go-ipfs-colog

colog is a **co**ncurrent **log**. In the JS version of Orbit it is currently calle log, but in many contexts a log is a linear sequence of events. To avoid confusion, the Go version was named colog.

## Installing

Note: You need [gx] installed.

First fetch the package:
```
go get -d github.com/keks/go-ipfs-colog
```
(`-d` means only download, don't build. That probably would work.)

Then, `cd` into the directory
```
cd $GOPATH/src/github.com/keks/go-ipfs-colog
```
and use gx to install the deps:
```
gx i
```
then build and install the package:
```
go install . ./immutabledb/...
```

The executable located in `cmd/ipfs-colog` is not done yet. That is also why the immutabledb for ipfs is not build by default.

## Test
Once successfully installed

```
go test github.com/keks/go-ipfs-colog
```
does the trick.

## Run
```
ipfs-colog
```

gx: https://github.com/whyrysleeping/gx

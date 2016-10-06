# go-ipfs-colog

colog is a *co*ncurrent *log*. In the JS version of Orbit it is currently calle log, but in many contexts a log is a linear sequence of events. To avoid confusion, the Go version was named colog.

## Installing
I wish `go install github.com/keks/go-ipfs-colog` would suffice, but we're not quite there yet. Also this package is not gx'd yet.

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

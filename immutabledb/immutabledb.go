package immutabledb

type ImmutableDB interface {
	Put([]byte) string
	Get(string) []byte
	Close() error
}

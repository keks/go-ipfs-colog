package immutabledb

type ImmutableDB interface {
	Put([]byte) (string, error)
	Get(string) ([]byte, error)
	Close() error
}

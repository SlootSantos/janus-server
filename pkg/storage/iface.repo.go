//go:generate mockgen -source=$GOFILE -destination=mock.$GOFILE -package=$GOPACKAGE

package storage

// RepoIface ...
type RepoIface interface {
	Get(key string) (string, error)
	Set(key string, value []byte) (string, error)
}

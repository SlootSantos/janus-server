//go:generate mockgen -source=$GOFILE -destination=mock.$GOFILE -package=$GOPACKAGE

package storage

// UserIface ...
type UserIface interface {
	Get(key string) (*UserModel, error)
	Set(key string, value *UserModel) error
}

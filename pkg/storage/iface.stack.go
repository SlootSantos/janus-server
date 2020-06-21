//go:generate mockgen -source=$GOFILE -destination=mock.$GOFILE -package=$GOPACKAGE

package storage

// StackIface ...
type StackIface interface {
	Get(key string) ([]StackModel, error)
	Set(key string, value []StackModel) error
}

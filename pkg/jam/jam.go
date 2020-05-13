//go:generate mockgen -source=$GOFILE -destination=mock.$GOFILE -package=$GOPACKAGE

package jam

import (
	"context"

	"github.com/SlootSantos/janus-server/pkg/storage"
)

// Stack represents an entire JAM Stack
type Stack = storage.StackModel

// StackCDN contains all stack relevant information about the CDN
type StackCDN = storage.StackCDNModel

// StackRepo contains all stack relevant information about the git repository
type StackRepo = storage.RepoModel

// CreationParam contains a compound parameter type for all strackResource creations
type CreationParam struct {
	ID     string
	Bucket struct {
		ID string
	}
	Repo struct {
		Name string
	}
}

// OutputParam contains a compound parameter type for every created stackresource
type OutputParam Stack

// DeletionParam contains a compound parameter type for all strackResource deletions
type DeletionParam Stack

// all resources used to create a JAM-Stack shall have these methods
type stackResource interface {
	Create(context.Context, *CreationParam, *OutputParam) (string, error)
	Destroy(context.Context, *DeletionParam) error
	List(context.Context) string
}

// Creator contains all resources to create a JAM-Stack
type Creator struct {
	resources []stackResource
}

// New creates a new JAM-Stack creator
func New(resources ...stackResource) *Creator {
	return &Creator{
		resources: resources,
	}
}

package userrepo

import (
	"sync"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/xdorro/golang-grpc-base-project/pkg/repo"
)

var _ IRepo = (*Repo)(nil)

// IRepo is the interface that must be implemented by a repository.
type IRepo interface {
}

// Repo is a repository struct.
type Repo struct {
	mu         sync.Mutex
	collection *mongo.Collection
}

// Option service option.
type Option struct {
	Repo repo.IRepo
}

// NewRepo creates a new repository.
func NewRepo(opt *Option) IRepo {
	r := &Repo{
		collection: opt.Repo.Collection("users"),
	}

	return r
}

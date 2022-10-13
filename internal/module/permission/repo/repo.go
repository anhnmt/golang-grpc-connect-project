package permissionrepo

import (
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	permissionmodel "github.com/xdorro/golang-grpc-base-project/internal/module/permission/model"
	"github.com/xdorro/golang-grpc-base-project/pkg/repo"
)

var _ IRepo = (*Repo)(nil)

// IRepo is the interface that must be implemented by a repository.
type IRepo interface {
	CountDocuments(filter any) (int64, error)
	Find(filter any, opt ...*options.FindOptions) ([]*permissionmodel.Permission, error)
	FindOne(filter any, opt ...*options.FindOneOptions) (*permissionmodel.Permission, error)
	InsertOne(data any, opt ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	UpdateOne(filter, data any, opt ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	SoftDeleteOne(filter any, opt ...*options.UpdateOptions) (*mongo.UpdateResult, error)
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
		collection: opt.Repo.CollectionModel(&permissionmodel.Permission{}),
	}

	return r
}

func (r *Repo) CountDocuments(filter any) (int64, error) {
	return repo.CountDocuments(r.collection, filter)
}

func (r *Repo) Find(filter any, opt ...*options.FindOptions) ([]*permissionmodel.Permission, error) {
	return repo.Find[permissionmodel.Permission](r.collection, filter, opt...)
}

func (r *Repo) FindOne(filter any, opt ...*options.FindOneOptions) (*permissionmodel.Permission, error) {
	return repo.FindOne[permissionmodel.Permission](r.collection, filter, opt...)
}

func (r *Repo) InsertOne(data any, opt ...*options.InsertOneOptions) (
	*mongo.InsertOneResult, error,
) {
	return repo.InsertOne(r.collection, data, opt...)
}

func (r *Repo) UpdateOne(filter, data any, opt ...*options.UpdateOptions) (
	*mongo.UpdateResult, error,
) {
	return repo.UpdateOne(r.collection, filter, data, opt...)
}

func (r *Repo) SoftDeleteOne(filter any, opt ...*options.UpdateOptions) (
	*mongo.UpdateResult, error,
) {
	return repo.SoftDeleteOne(r.collection, filter, opt...)
}

package repo

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ IRepo = (*Repo)(nil)

// IRepo is the interface that must be implemented by a repository.
type IRepo interface {
	Close() error
	Collection(name string) *mongo.Collection
}

// Repo is a repository struct.
type Repo struct {
	mu     sync.Mutex
	ctx    context.Context
	dbURL  string
	dbName string

	client *mongo.Client
}

// NewRepo creates a new repository.
func NewRepo() IRepo {
	r := &Repo{
		ctx:    context.Background(),
		dbURL:  viper.GetString("DB_URL"),
		dbName: viper.GetString("DB_NAME"),
	}

	log.Info().
		Str("db_url", r.dbURL).
		Str("db_name", r.dbName).
		Msg("Connecting to MongoDB")

	// Set connect timeout to 15 seconds
	ctxConn, cancel := context.WithTimeout(r.ctx, 15*time.Second)
	defer cancel()

	clientOpts := options.Client().
		ApplyURI(r.dbURL)

	// Create new client and connect to MongoDB
	client, err := mongo.Connect(ctxConn, clientOpts)
	if err != nil {
		log.Panic().Err(err).Msg("Connecting to MongoDB failed")
	}

	// Ping the primary
	if err = client.Ping(ctxConn, nil); err != nil {
		log.Panic().Err(err).Msg("Ping to MongoDB failed")
	}

	// Add client to repository
	r.setClient(client)

	log.Info().Msg("Connecting to MongoDB successfully.")

	return r
}

// Close closes the repository.
func (r *Repo) Close() error {
	ctx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()

	if err := r.client.Disconnect(ctx); err != nil {
		log.Err(err).Msg("Failed to disconnect from MongoDB")
		return err
	}

	return nil
}

// Collection returns the mongo collection by Name.
func (r *Repo) Collection(name string) *mongo.Collection {
	return r.client.Database(r.dbName).Collection(name)
}

// setClient adds a new client to the repository.
func (r *Repo) setClient(client *mongo.Client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.client = client
}

package casbin

import (
	"sync"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	mongodbadapter "github.com/casbin/mongodb-adapter/v3"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"github.com/xdorro/golang-grpc-base-project/pkg/repo"
)

var _ ICasbin = (*Casbin)(nil)

// ICasbin is the interface that must be implemented by a casbin.
type ICasbin interface {
	Enforcer() *casbin.CachedEnforcer
}

// Option casbin option.
type Option struct {
	Repo repo.IRepo
}

// Casbin is a casbin struct.
type Casbin struct {
	mu          sync.Mutex
	dbName      string
	casbinModel string
	casbinName  string
	enforcer    *casbin.CachedEnforcer

	// options
	repo repo.IRepo
}

// NewCasbin creates a new casbin.
func NewCasbin(opt *Option) ICasbin {
	c := &Casbin{
		repo:        opt.Repo,
		dbName:      viper.GetString("database.name"),
		casbinModel: viper.GetString("casbin.model"),
		casbinName:  viper.GetString("casbin.name"),
	}

	log.Info().Msg("Connecting to Casbin")

	config := &mongodbadapter.AdapterConfig{
		DatabaseName:   c.dbName,
		CollectionName: c.casbinName,
	}
	adapter, err := mongodbadapter.NewAdapterByDB(c.repo.Client(), config)
	if err != nil {
		log.Panic().Err(err).Msg("Failed to create mongodb adapter")
	}

	m, _ := model.NewModelFromString(c.casbinModel)

	enforcer, err := casbin.NewCachedEnforcer(m, adapter)
	if err != nil {
		log.Panic().Err(err).Msg("Failed to create casbin enforcer")
	}

	// Load the policy from DB.
	if err = enforcer.LoadPolicy(); err != nil {
		log.Panic().Err(err).Msg("Failed to load policy")
	}

	// Add enforcer to Casbin.
	c.setClient(enforcer)

	log.Info().Msg("Loaded casbin successfully")

	return c
}

// Client adds a new client to the repository.
func (c *Casbin) setClient(enforcer *casbin.CachedEnforcer) {
	c.mu.Lock()
	c.enforcer = enforcer
	c.mu.Unlock()
}

// Enforcer return cache enforcer
func (c *Casbin) Enforcer() *casbin.CachedEnforcer {
	return c.enforcer
}

package usermodel

import (
	"github.com/rs/zerolog/log"
	userv1 "github.com/xdorro/proto-base-project/proto-gen-go/user/v1"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	"github.com/xdorro/golang-grpc-base-project/internal/model"
	"github.com/xdorro/golang-grpc-base-project/pkg/utils"
)

var _ IUser = &User{}

// IUser is the interface for a user
type IUser interface {
	model.IBaseModel

	HashPassword() error
	ComparePassword(password string) bool
}

// User is a user struct.
type User struct {
	model.BaseModel `bson:",inline"`

	Name     string `json:"name,omitempty" bson:"name,omitempty"`
	Email    string `json:"email,omitempty" bson:"email,omitempty"`
	Password string `json:"password" bson:"password,omitempty"`
	Role     string `json:"role" bson:"role,omitempty"`
	Status   int32  `json:"status,omitempty" bson:"status,omitempty"`
}

// CollectionName returns the name of the collection from struct name
func (m *User) CollectionName() string {
	return utils.CollectionName(m)
}

// GetIndexModels returns the index models
func (m *User) GetIndexModels() []mongo.IndexModel {
	return []mongo.IndexModel{}
}

// PreCreate is a callback that gets called before creating a models.
func (m *User) PreCreate() {
	m.BaseModel.PreCreate()
}

// PreUpdate is a callback that gets called before updating a models.
func (m *User) PreUpdate() {
	m.BaseModel.PreUpdate()
}

// HashPassword hashes a password
func (m *User) HashPassword() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(m.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Err(err).Msg("Error hash password")
		return err
	}

	m.Password = string(bytes)
	return nil
}

// ComparePassword compares a password with a hash
func (m *User) ComparePassword(password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(m.Password), []byte(password)); err != nil {
		log.Err(err).Msg("Error compare hash and password")
		return false
	}

	return true
}

// UserToProto converts a user to a proto
func UserToProto(m *User) *userv1.User {
	return &userv1.User{
		Id:    m.ID,
		Name:  m.Name,
		Email: m.Email,
		Role:  m.Role,
	}
}

// UsersToProto converts a slice of users to a slice of proto
func UsersToProto(list []*User) []*userv1.User {
	return utils.ToProto[User, userv1.User](list, UserToProto)
}

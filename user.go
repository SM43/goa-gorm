package userapi

import (
	"context"
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	user "github.com/sm43/goa-gorm/gen/user"
)

var userStore = make([]*user.StoredUser, 0)

// user service example implementation.
// The example methods log the requests and return zero values.
type usersrvc struct {
	db     *gorm.DB
	logger *log.Logger
}

// NewUser returns the user service implementation.
func NewUser(db *gorm.DB, logger *log.Logger) user.Service {
	return &usersrvc{db, logger}
}

// Add new user and return its ID.
func (s *usersrvc) Add(ctx context.Context, p *user.User) (res *user.User, err error) {
	res = &user.User{}
	s.logger.Print("user.add")

	item := user.StoredUser{ID: *p.ID, Name: *p.Name}
	userStore = append(userStore, &item)

	res = (&user.User{ID: p.ID, Name: p.Name})
	err = s.db.Create(&User{Name: *p.Name}).Error

	if err != nil {
		return nil, user.MakeDbError(fmt.Errorf(err.Error()))
	}

	s.logger.Print("Array Size - ", len(userStore))
	return
}

// List all users
func (s *usersrvc) List(ctx context.Context) (res []*user.StoredUser, err error) {
	s.logger.Print("user.list")

	var all []User
	err = s.db.Find(&all).Error
	if err != nil {
		return nil, user.MakeDbError(fmt.Errorf(err.Error()))
	}

	ret := make([]*user.StoredUser, len(all))
	for i, r := range all {
		ret[i] = Init(r)
	}
	return ret, nil
}

// Init Converts
func Init(u User) *user.StoredUser {
	return &user.StoredUser{
		ID:   uint64(u.ID),
		Name: u.Name,
	}
}

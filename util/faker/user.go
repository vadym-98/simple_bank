package faker

import (
	db "github.com/vadym-98/simple_bank/db/sqlc"
	"github.com/vadym-98/simple_bank/util"
)

type UserBuilder struct {
	user db.User
}

func (ub *UserBuilder) Get() db.User {
	return ub.user
}

func NewUser() *UserBuilder {
	return &UserBuilder{
		user: db.User{
			Username: util.RandomOwner(),
			FullName: util.RandomOwner(),
			Email:    util.RandomEmail(),
		},
	}
}

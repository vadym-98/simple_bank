package faker

import (
	db "github.com/vadym-98/simple_bank/db/sqlc"
	"github.com/vadym-98/simple_bank/util"
)

type accountBuilder struct {
	account db.Account
}

func (ab *accountBuilder) WithBalance(b int64) *accountBuilder {
	ab.account.Balance = b
	return ab
}

func (ab *accountBuilder) WithCurrency(c string) *accountBuilder {
	ab.account.Currency = c
	return ab
}

func (ab *accountBuilder) Get() db.Account {
	return ab.account
}

func NewAccount() *accountBuilder {
	return &accountBuilder{
		account: db.Account{
			ID:       util.RandomInt(1, 1000),
			Owner:    util.RandomOwner(),
			Balance:  util.RandomMoney(),
			Currency: util.RandomCurrency(),
		},
	}
}

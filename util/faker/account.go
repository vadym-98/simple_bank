package faker

import (
	db "github.com/vadym-98/simple_bank/db/sqlc"
	"github.com/vadym-98/simple_bank/util"
)

type AccountBuilder struct {
	account db.Account
}

func (ab *AccountBuilder) WithBalance(b int64) *AccountBuilder {
	ab.account.Balance = b
	return ab
}

func (ab *AccountBuilder) WithCurrency(c string) *AccountBuilder {
	ab.account.Currency = c
	return ab
}

func (ab *AccountBuilder) Get() db.Account {
	return ab.account
}

func NewAccount() *AccountBuilder {
	return &AccountBuilder{
		account: db.Account{
			ID:       util.RandomInt(1, 1000),
			Owner:    util.RandomOwner(),
			Balance:  util.RandomMoney(),
			Currency: util.RandomCurrency(),
		},
	}
}

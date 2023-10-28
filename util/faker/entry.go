package faker

import (
	db "github.com/vadym-98/simple_bank/db/sqlc"
	"github.com/vadym-98/simple_bank/util"
)

type entryBuilder struct {
	entry db.Entry
}

func (eb *entryBuilder) WithAccountID(id int64) *entryBuilder {
	eb.entry.AccountID = id
	return eb
}

func (eb *entryBuilder) WithAmount(amount int64) *entryBuilder {
	eb.entry.Amount = amount
	return eb
}

func (eb *entryBuilder) Get() db.Entry {
	return eb.entry
}

func NewEntry() *entryBuilder {
	return &entryBuilder{
		entry: db.Entry{
			ID:        util.RandomInt(1, 1000),
			AccountID: util.RandomInt(1, 1000),
			Amount:    util.RandomMoney(),
		},
	}
}

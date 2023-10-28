package faker

import (
	db "github.com/vadym-98/simple_bank/db/sqlc"
	"github.com/vadym-98/simple_bank/util"
)

type EntryBuilder struct {
	entry db.Entry
}

func (eb *EntryBuilder) WithAccountID(id int64) *EntryBuilder {
	eb.entry.AccountID = id
	return eb
}

func (eb *EntryBuilder) WithAmount(amount int64) *EntryBuilder {
	eb.entry.Amount = amount
	return eb
}

func (eb *EntryBuilder) Get() db.Entry {
	return eb.entry
}

func NewEntry() *EntryBuilder {
	return &EntryBuilder{
		entry: db.Entry{
			ID:        util.RandomInt(1, 1000),
			AccountID: util.RandomInt(1, 1000),
			Amount:    util.RandomMoney(),
		},
	}
}

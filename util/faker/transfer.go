package faker

import (
	db "github.com/vadym-98/simple_bank/db/sqlc"
	"github.com/vadym-98/simple_bank/util"
)

type transferBuilder struct {
	transfer db.Transfer
}

func (tb *transferBuilder) WithFromAccountID(id int64) *transferBuilder {
	tb.transfer.FromAccountID = id
	return tb
}

func (tb *transferBuilder) WithToAccountID(id int64) *transferBuilder {
	tb.transfer.ToAccountID = id
	return tb
}

func (tb *transferBuilder) Get() db.Transfer {
	return tb.transfer
}

func NewTransfer() *transferBuilder {
	return &transferBuilder{
		transfer: db.Transfer{
			ID:            util.RandomInt(1, 1000),
			FromAccountID: util.RandomInt(1, 1000),
			ToAccountID:   util.RandomInt(1, 1000),
			Amount:        util.RandomMoney(),
		},
	}
}

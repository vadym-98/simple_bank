package faker

import (
	db "github.com/vadym-98/simple_bank/db/sqlc"
	"github.com/vadym-98/simple_bank/util"
)

type TransferBuilder struct {
	transfer db.Transfer
}

func (tb *TransferBuilder) WithFromAccountID(id int64) *TransferBuilder {
	tb.transfer.FromAccountID = id
	return tb
}

func (tb *TransferBuilder) WithToAccountID(id int64) *TransferBuilder {
	tb.transfer.ToAccountID = id
	return tb
}

func (tb *TransferBuilder) Get() db.Transfer {
	return tb.transfer
}

func NewTransfer() *TransferBuilder {
	return &TransferBuilder{
		transfer: db.Transfer{
			ID:            util.RandomInt(1, 1000),
			FromAccountID: util.RandomInt(1, 1000),
			ToAccountID:   util.RandomInt(1, 1000),
			Amount:        util.RandomMoney(),
		},
	}
}

package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/vadym-98/simple_bank/util"
)

var validCurrency validator.Func = func(f validator.FieldLevel) bool {
	if currency, ok := f.Field().Interface().(string); ok {
		return util.IsSupportedCurrency(currency)
	}

	return false
}

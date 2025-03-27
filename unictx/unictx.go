package unictx

import (
	"fmt"
	"github.com/hootuu/eggcone/unictx/modelx"
	"github.com/hootuu/gelato/errors"
	"github.com/hootuu/gelato/strs"
)

func Set(domain string, key string, value string) *errors.Error {
	return modelx.Set(domain, key, value)
}

func Get(domain string, key string, dfVal string) (string, *errors.Error) {
	return modelx.Get(domain, key, dfVal)
}

func SetInt64(domain string, key string, value int64) *errors.Error {
	strVal := fmt.Sprintf("%d", value)
	return modelx.Set(domain, key, strVal)
}

func GetInt64(domain string, key string, dfVal int64) (int64, *errors.Error) {
	strVal, err := modelx.Get(domain, key, fmt.Sprintf("%d", dfVal))
	if err != nil {
		return 0, err
	}
	intVal := strs.ToInt64(strVal)
	return intVal, nil
}

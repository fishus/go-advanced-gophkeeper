package domain

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/theplant/luhn"
)

type CardExpDate struct {
	Month uint
	Year  uint
}

type VaultDataCard struct {
	Info       string      `json:"info"`
	Number     string      `json:"number"`
	HolderName string      `json:"holder_name"`
	ExpDate    CardExpDate `json:"exp_date"`
	CvcCode    string      `json:"cvc_code"`
}

var _ IVaultRecordData = (*VaultDataCard)(nil)

func (v VaultDataCard) GetInfo() string {
	return v.Info
}

// MarshalJSON implements the json.Marshaler interface.
func (v VaultDataCard) MarshalJSON() ([]byte, error) {
	type VaultDataAlias VaultDataCard

	aliasValue := struct {
		VaultDataAlias
		ExpDate string `json:"exp_date"`
	}{
		VaultDataAlias: VaultDataAlias(v),
		ExpDate:        fmt.Sprintf("%02d/%02d", v.ExpDate.Month, v.ExpDate.Year),
	}

	return json.Marshal(aliasValue)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (v *VaultDataCard) UnmarshalJSON(data []byte) (err error) {
	type VaultDataAlias VaultDataCard

	aliasValue := &struct {
		*VaultDataAlias
		ExpDate string `json:"exp_date"`
	}{
		VaultDataAlias: (*VaultDataAlias)(v),
	}

	if err = json.Unmarshal(data, aliasValue); err != nil {
		return
	}

	var expDateMonth uint64
	var expDateYear uint64
	if len(aliasValue.ExpDate) == 5 && strings.Contains(aliasValue.ExpDate, "/") {
		expDateMonth, err = strconv.ParseUint(aliasValue.ExpDate[:2], 10, 0)
		if err != nil {
			return ErrInvalidCardExpDate
		}

		expDateYear, err = strconv.ParseUint(aliasValue.ExpDate[3:], 10, 0)
		if err != nil {
			return ErrInvalidCardExpDate
		}
	} else {
		return ErrInvalidCardExpDate
	}

	v.ExpDate.Month = uint(expDateMonth)
	v.ExpDate.Year = uint(expDateYear)

	return
}

func (v VaultDataCard) Validate() error {
	// Validate card number
	num, err := strconv.Atoi(strings.ReplaceAll(v.Number, " ", ""))
	if err != nil {
		return ErrIncorrectCardNumber
	}
	if !luhn.Valid(num) {
		return ErrIncorrectCardNumber
	}

	// Validate exp date (MM/DD)
	// Valid exp.date from 00/00 to 12/99
	if v.ExpDate.Month > 12 || v.ExpDate.Year > 99 {
		return ErrIncorrectCardExpDate
	}

	// Validate cvc code
	if len(v.CvcCode) > 3 || strings.ContainsFunc(v.CvcCode, func(c rune) bool {
		return c < '0' || c > '9'
	}) {
		return ErrIncorrectCardCvcCode
	}

	return nil
}

package domain

import (
	"encoding/json"
	"fmt"
	"strconv"
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

	expDateMonth, err := strconv.ParseUint(aliasValue.ExpDate[:2], 10, 0)
	if err != nil {
		return
	}

	expDateYear, err := strconv.ParseUint(aliasValue.ExpDate[3:], 10, 0)
	if err != nil {
		return
	}

	v.ExpDate.Month = uint(expDateMonth)
	v.ExpDate.Year = uint(expDateYear)

	return
}

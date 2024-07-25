package domain

import (
	"encoding/json"
	"strconv"
)

const VaultFileMaxFilesize = 1024 * 1024

type VaultDataFile struct {
	Info     string `json:"info"`
	Filename string `json:"filename"`
	MimeType string `json:"mime_type"`
	Filesize uint64 `json:"filesize"`
	Data     []byte `json:"data,omitempty"`
}

var _ IVaultRecordData = (*VaultDataFile)(nil)

func (v VaultDataFile) GetInfo() string {
	return v.Info
}

// MarshalJSON implements the json.Marshaler interface.
func (v VaultDataFile) MarshalJSON() ([]byte, error) {
	type VaultDataAlias VaultDataFile

	aliasValue := struct {
		VaultDataAlias
		Filesize string `json:"filesize"`
		Data     string `json:"data,omitempty"`
	}{
		VaultDataAlias: VaultDataAlias(v),
		Filesize:       strconv.FormatUint(v.Filesize, 10),
		Data:           string(v.Data),
	}

	return json.Marshal(aliasValue)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (v *VaultDataFile) UnmarshalJSON(data []byte) (err error) {
	type VaultDataAlias VaultDataFile

	aliasValue := &struct {
		*VaultDataAlias
		Filesize string `json:"filesize"`
		Data     string `json:"data,omitempty"`
	}{
		VaultDataAlias: (*VaultDataAlias)(v),
	}

	if err = json.Unmarshal(data, aliasValue); err != nil {
		return
	}

	v.Filesize, err = strconv.ParseUint(aliasValue.Filesize, 10, 0)
	if err != nil {
		return
	}
	v.Data = []byte(aliasValue.Data)

	return
}

func (v VaultDataFile) Validate() error {
	if v.Filesize > VaultFileMaxFilesize {
		return ErrVaultMaxFilesize
	}

	return nil
}

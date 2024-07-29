package domain

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

func (v VaultDataFile) Validate() error {
	if v.Filesize > VaultFileMaxFilesize {
		return ErrVaultMaxFilesize
	}

	return nil
}

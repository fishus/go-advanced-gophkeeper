package domain

type VaultDataNote struct {
	Info    string `json:"info"`
	Content string `json:"content"`
}

var _ IVaultRecordData = (*VaultDataNote)(nil)

func (r VaultDataNote) GetInfo() string {
	return r.Info
}

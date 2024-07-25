package domain

type VaultDataCreds struct {
	Info     string `json:"info"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

var _ IVaultRecordData = (*VaultDataCreds)(nil)

func (v VaultDataCreds) GetInfo() string {
	return v.Info
}

func (v VaultDataCreds) Validate() error {
	return nil
}

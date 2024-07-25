package domain

const VaultNoteMaxLen = 1024 * 1024

type VaultDataNote struct {
	Info    string `json:"info"`
	Content string `json:"content"`
}

var _ IVaultRecordData = (*VaultDataNote)(nil)

func (v VaultDataNote) GetInfo() string {
	return v.Info
}

func (v VaultDataNote) Validate() error {
	if len(v.Content) > VaultNoteMaxLen {
		return ErrVaultNoteMaxLen
	}

	return nil
}

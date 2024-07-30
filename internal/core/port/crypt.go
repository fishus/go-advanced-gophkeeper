package port

import (
	"context"
)

type CryptAdapter interface {
	// EncryptSymmetric encrypts data with a symmetric key
	EncryptSymmetric(context.Context, []byte) ([]byte, error)

	// DecryptSymmetric decrypts data with a symmetric key
	DecryptSymmetric(context.Context, []byte) ([]byte, error)
}

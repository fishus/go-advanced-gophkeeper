package gsm

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"

	"github.com/fishus/go-advanced-gophkeeper/internal/core/port"
)

// CryptAdapter implements port.CryptAdapter interface
type CryptAdapter struct {
	key [32]byte
}

// New creates a new crypt adapter instance
func New(secret []byte) (port.CryptAdapter, error) {
	key := sha256.Sum256(secret)

	return &CryptAdapter{
		key,
	}, nil
}

// EncryptSymmetric encrypts data with a symmetric key
func (c *CryptAdapter) EncryptSymmetric(ctx context.Context, src []byte) (encrypted []byte, err error) {
	aesblock, err := aes.NewCipher(c.key[:])
	if err != nil {
		return
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return
	}

	nonce := c.key[(len(c.key) - aesgcm.NonceSize()):]
	encrypted = aesgcm.Seal(nil, nonce, src, nil)
	return
}

// DecryptSymmetric encrypts data with a symmetric key
func (c *CryptAdapter) DecryptSymmetric(ctx context.Context, encrypted []byte) (decrypted []byte, err error) {
	aesblock, err := aes.NewCipher(c.key[:])
	if err != nil {
		return
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return
	}

	nonce := c.key[(len(c.key) - aesgcm.NonceSize()):]
	decrypted, err = aesgcm.Open(nil, nonce, encrypted, nil)
	return
}

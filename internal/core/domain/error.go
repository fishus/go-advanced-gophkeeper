package domain

import "errors"

var (
	ErrInvalidArgument = errors.New("invalid argument")
	ErrAlreadyExists   = errors.New("already exists")
	ErrNotFound        = errors.New("not found")
	ErrUserNotFound    = errors.New("user not found")
	ErrUserIDNotSet    = errors.New("user id not set")
	ErrTimeout         = errors.New("timeout exceeded")

	ErrTokenDuration      = errors.New("invalid token duration format")
	ErrTokenCreation      = errors.New("error creating token")
	ErrExpiredToken       = errors.New("access token has expired")
	ErrInvalidToken       = errors.New("access token is invalid")
	ErrInvalidCredentials = errors.New("incorrect login or password")

	ErrUndefinedVaultKind   = errors.New("vault kind is not defined")
	ErrInvalidVaultKind     = errors.New("incorrect vault kind")
	ErrInvalidCardExpDate   = errors.New("invalid card expiration date (MM/YY)")
	ErrVaultMaxFilesize     = errors.New("maximum file size exceeded")
	ErrVaultNoteMaxLen      = errors.New("maximum length of note content exceeded")
	ErrIncorrectCardNumber  = errors.New("incorrect card number")
	ErrIncorrectCardExpDate = errors.New("incorrect card exp. date")
	ErrIncorrectCardCvcCode = errors.New("incorrect card cvc code")

	ErrVaultRecordNotCreated = errors.New("vault record was not created")
)

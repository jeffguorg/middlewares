package signature

import (
	"crypto"
	"crypto/hmac"
	"encoding/base64"
	"github.com/go-errors/errors"
)

var (
	ErrHashUnavailable  = errors.New("the requested hash function is unavailable")
	ErrSignatureInvalid = errors.New("signature is invalid")
)

type SigningMethodHMAC struct {
	Key        []byte
	HashMethod crypto.Hash
}

func (method SigningMethodHMAC) Verify(signingString, signature string) error {
	if calculated, err := method.Sign(signingString); err != nil {
		return err
	} else if calculated != signature {
		return ErrSignatureInvalid
	}
	return nil
}

func (method SigningMethodHMAC) Sign(signingString string) (string, error) {
	if !method.HashMethod.Available() {
		return "", ErrHashUnavailable
	}

	hasher := hmac.New(method.HashMethod.New, method.Key)
	hasher.Write([]byte(signingString))

	return base64.RawURLEncoding.EncodeToString(hasher.Sum(nil)), nil
}

var (
	_ SigningMethod = SigningMethodHMAC{}
)

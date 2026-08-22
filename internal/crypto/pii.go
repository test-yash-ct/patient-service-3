package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"strings"
)

const prefix = "enc:v1:"

type FieldEncryptor struct {
	gcm cipher.AEAD
}

func NewFieldEncryptor(key []byte) (*FieldEncryptor, error) {
	if len(key) != 32 {
		return nil, errors.New("pii key must be 32 bytes")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return &FieldEncryptor{gcm: gcm}, nil
}

func (e *FieldEncryptor) Encrypt(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}
	if strings.HasPrefix(plaintext, prefix) {
		return plaintext, nil
	}
	nonce := make([]byte, e.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	out := e.gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return prefix + base64.RawStdEncoding.EncodeToString(out), nil
}

func (e *FieldEncryptor) Decrypt(value string) (string, error) {
	if value == "" {
		return "", nil
	}
	if !strings.HasPrefix(value, prefix) {
		return "", errors.New("unencrypted phi refused")
	}
	raw, err := base64.RawStdEncoding.DecodeString(strings.TrimPrefix(value, prefix))
	if err != nil {
		return "", err
	}
	ns := e.gcm.NonceSize()
	if len(raw) < ns {
		return "", errors.New("ciphertext too short")
	}
	plain, err := e.gcm.Open(nil, raw[:ns], raw[ns:], nil)
	if err != nil {
		return "", err
	}
	return string(plain), nil
}

func MaskSSN(ssn string) string {
	digits := make([]rune, 0, len(ssn))
	for _, r := range ssn {
		if r >= '0' && r <= '9' {
			digits = append(digits, r)
		}
	}
	if len(digits) < 4 {
		return "***"
	}
	return "***-**-" + string(digits[len(digits)-4:])
}

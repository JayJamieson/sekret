package handlers

import (
	"crypto/aes"
	"crypto/sha256"
)

func Seal(value []byte, passphrase []byte) ([]byte, error) {

	passphraseHash := sha256.Sum256(passphrase)
	aesBlock, err := aes.NewCipher(passphraseHash[:])

	if err != nil {
		return nil, err
	}

	buf := []byte{}

	aesBlock.Encrypt(buf, value)

	return buf, nil
}

func UnSeal(value []byte, passphrase []byte) ([]byte, error) {
	passphraseHash := sha256.Sum256(passphrase)
	aesBlock, err := aes.NewCipher(passphraseHash[:])

	if err != nil {
		return nil, err
	}

	buf := []byte{}

	aesBlock.Decrypt(buf, value)

	return buf, nil
}

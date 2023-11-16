package sec

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

func Encrypt(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, 12) // creating 12 random octets used once
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	Ctext := aesgcm.Seal(nil, nonce, data, nil) // Ciphering using cipher block chaining  methode
	Ctext = append(nonce, Ctext...)             // Concatenating generated nonce with ciphered text
	return Ctext, nil
}

func Decrypt(Ctext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(Ctext) < aes.BlockSize {
		return nil, errors.New("invalid Ctext")
	}

	nonce := Ctext[:12] // retrievieng nonce from ciphered text
	Ctext = Ctext[12:]

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err := aesgcm.Open(nil, nonce, Ctext, nil) // unciphering using cipher block chaining method
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

package converter

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

type Encryptor struct {
	Converter Converter
	key       []byte
}

// in -> data -> out

func (e *Encryptor) Encode(in []byte) ([]byte, error) {
	data, err := e.Converter.Encode(in)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return aesGCM.Seal(nonce, nonce, data, nil), nil
}

func (e *Encryptor) Decode(in []byte) ([]byte, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := aesGCM.NonceSize()
	nonce, cipherData := in[:nonceSize], in[nonceSize:]

	data, err := aesGCM.Open(nil, nonce, cipherData, nil)
	if err != nil {
		return nil, err
	}

	out, err := e.Converter.Decode(data)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (e *Encryptor) SetKey(key []byte) {
	e.key = key
}

package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

type CryptoService struct {
	gcm cipher.AEAD
}

func NewCryptoService(key []byte) (cs *CryptoService, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return &CryptoService{
		gcm: gcm,
	}, nil
}

func (cs CryptoService) Encrypt(plain []byte) (cipher []byte, err error) {
	nonce := make([]byte, cs.gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}

	cipher = cs.gcm.Seal(nil, nonce, plain, nil)
	cipher = append(nonce, cipher...)

	return cipher, nil
}

func (cs CryptoService) EncryptBase64(plain []byte) (base64Cipher []byte, err error) {
	cipher, err := cs.Encrypt(plain)
	if err != nil {
		return nil, err
	}

	base64Cipher = make([]byte, base64.RawStdEncoding.EncodedLen(len(cipher)))
	base64.RawStdEncoding.Encode(base64Cipher, cipher)

	return
}

func (cs CryptoService) Decrypt(cipher []byte) (plain []byte, err error) {
	nonce := cipher[0:cs.gcm.NonceSize()]
	ciphertext := cipher[cs.gcm.NonceSize():]

	plain, err = cs.gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return
}

func (cs CryptoService) DecryptBase64(base64Cipher []byte) (plain []byte, err error) {
	cipher := make([]byte, base64.RawStdEncoding.DecodedLen(len(base64Cipher)))
	_, err = base64.RawStdEncoding.Decode(cipher, base64Cipher)
	if err != nil {
		return nil, err
	}

	return cs.Decrypt(cipher)
}

package main

import (
	"code.google.com/p/go.crypto/scrypt"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

func EncryptAESCFB(dst, src, key, iv []byte) error {
	aesBlockEncrypter, err := aes.NewCipher([]byte(key))
	if err != nil {
		return err
	}
	aesEncrypter := cipher.NewCFBEncrypter(aesBlockEncrypter, iv)
	aesEncrypter.XORKeyStream(dst, src)
	return nil
}

func DecryptAESCFB(dst, src, key, iv []byte) error {
	aesBlockDecrypter, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil
	}
	aesDecrypter := cipher.NewCFBDecrypter(aesBlockDecrypter, iv)
	aesDecrypter.XORKeyStream(dst, src)
	return nil
}

func Encrypt(b []byte, password string) []byte {
	key := generateKey(password)
	iv := generateIV()
	encrypted := make([]byte, len(b))
	err := EncryptAESCFB(encrypted, b, key, iv)
	if err != nil {
		panic(err)
	}
	result := make([]byte, len(iv)+len(encrypted))

	resultIndex := 0
	for _, i := range iv {
		result[resultIndex] = i
		resultIndex += 1
	}
	for _, i := range encrypted {
		result[resultIndex] = i
		resultIndex += 1
	}
	return result
}

func Decrypt(b []byte, password string) ([]byte, bool) {
	key := generateKey(password)
	iv := b[:16]
	encrypted := b[16:]
	result := make([]byte, len(encrypted))
	err := DecryptAESCFB(result, encrypted, key, iv)
	if err != nil {
		return nil, false
	}
	return result, true
}

func generateIV() []byte {
	iv := make([]byte, aes.BlockSize)
	rand.Read(iv)
	return iv
}

func generateKey(password string) []byte {
	salt := []byte("bla bla bla")
	key, err := scrypt.Key([]byte(password), salt, 16384, 8, 1, 32)
	if err != nil {
		panic(err)
	}
	return key
}

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
		return err
	}
	aesDecrypter := cipher.NewCFBDecrypter(aesBlockDecrypter, iv)
	aesDecrypter.XORKeyStream(dst, src)
	return nil
}

func Encrypt(b []byte, password string) []byte {
	salt := generateSalt(256)
	key := generateKey(password, salt)
	iv := generateIV(aes.BlockSize)

	encrypted := make([]byte, len(b))
	err := EncryptAESCFB(encrypted, b, key, iv)
	if err != nil {
		panic(err)
	}
	result := []byte{}

	for _, i := range salt {
		result = append(result, i)
	}
	for _, i := range iv {
		result = append(result, i)
	}
	for _, i := range encrypted {
		result = append(result, i)
	}
	return result
}

func Decrypt(b []byte, password string) ([]byte, bool) {
	salt := b[:256]
	iv := b[256:(256 + aes.BlockSize)]
	key := generateKey(password, salt)
	encrypted := b[(256 + aes.BlockSize):]
	result := make([]byte, len(encrypted))
	err := DecryptAESCFB(result, encrypted, key, iv)
	if err != nil {
		return nil, false
	}
	return result, true
}

func generateIV(size int) []byte {
	iv := make([]byte, size)
	rand.Read(iv)
	return iv
}

func generateSalt(size int) []byte {
	iv := make([]byte, size)
	rand.Read(iv)
	return iv
}

func generateKey(password string, salt []byte) []byte {
	key, err := scrypt.Key([]byte(password), salt, 16384, 8, 1, 32)
	if err != nil {
		panic(err)
	}
	return key
}

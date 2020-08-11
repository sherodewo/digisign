package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
)

func pkcs7Unpad(data []byte) ([]byte, error) {
	length := len(data)
	unPadding := int(data[length-1])
	return data[:(length - unPadding)], nil
}

func decrypt(key []byte, iv []byte, encrypted string) ([]byte, error) {
	data, err := base64.RawStdEncoding.DecodeString(encrypted)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 || len(data)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("bad blocksize(%v), aes.BlockSize = %v\n", len(data), aes.BlockSize)
	}
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	cbc := cipher.NewCBCDecrypter(c, iv)
	cbc.CryptBlocks(data, data)
	out, err := pkcs7Unpad(data)
	if err != nil {
		return out, err
	}
	return out, nil
}

func DecryptCredential(key string, encryptedText string) (string, error) {
	s := strings.Split(encryptedText, ":")
	src, iv := s[0], s[1]

	h := sha256.New()
	h.Write([]byte(key))
	keyEncrypted := h.Sum(nil)

	decodeIv, err := hex.DecodeString(iv)
	decryptedText, err := decrypt(keyEncrypted, decodeIv, src)
	return string(decryptedText), err
}

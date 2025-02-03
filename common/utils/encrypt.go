package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

func AesEncrypt(key, originData []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	originData = PKCS5Padding(originData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	encryptedData := make([]byte, len(originData))
	blockMode.CryptBlocks(encryptedData, originData)
	return encryptedData, nil
}

func PKCS5Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

func AesDecrypt(key, encryptedData []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	originData := make([]byte, len(encryptedData))
	blockMode.CryptBlocks(originData, encryptedData)
	originData = PKCS5UnPadding(originData)
	return originData, nil
}

func PKCS5UnPadding(originData []byte) []byte {
	l := len(originData)
	unPadding := int(originData[l-1])
	if unPadding < 1 || unPadding > 32 {
		unPadding = 0
	}
	return originData[:(l - unPadding)]
}

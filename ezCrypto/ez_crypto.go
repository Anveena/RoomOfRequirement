package ezCrypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/md5"
	"fmt"
)

func freedomPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
func freedomUnPadding(plainText []byte) []byte {
	length := len(plainText)
	number := int((plainText)[length-1])
	return (plainText)[:length-number]
}
func TripleDESEncrypt(origData []byte, keyStr string) (_ []byte, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()
	key := []byte(keyStr)
	block, err := des.NewTripleDESCipher(key[:24])
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	tmpData := freedomPadding(origData, blockSize)
	tbc := make([]byte, len(tmpData))
	for i := 0; i+blockSize <= len(tmpData); i += blockSize {
		block.Encrypt(tbc[i:i+blockSize], tmpData[i:i+blockSize])
	}
	return tbc, nil
}
func MakeMD5Key(aesKey string, times int64) []byte {
	key := []byte(aesKey)
	for i := int64(0); i < times; i++ {
		md5Maker := md5.New()
		md5Maker.Write(key)
		key = md5Maker.Sum(nil)
	}
	return key
}
func AESCBCEncrypt(origData []byte, key []byte) (_ []byte, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	tmpData := freedomPadding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, (key)[:blockSize])
	crypted := make([]byte, len(tmpData))
	blockMode.CryptBlocks(crypted, tmpData)
	return crypted, nil
}
func AESCBCDecrypt(encData []byte, key []byte) (_ []byte, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, (key)[:blockSize])
	origData := make([]byte, len(encData))
	blockMode.CryptBlocks(origData, encData)
	origData = freedomUnPadding(origData)
	return origData, nil
}
func AESEncryptForJava(origData []byte, key []byte) (_ []byte, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := c.BlockSize()
	tmpData := freedomPadding(origData, blockSize)

	result := make([]byte, len(tmpData))
	for i := 0; i < len(tmpData); i += aes.BlockSize {
		c.Encrypt(result[i:], tmpData[i:])
	}
	return result, nil
}

func AESDecryptForJava(encData []byte, key []byte) (_ []byte, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	result := make([]byte, len(encData))
	for i := 0; i < len(encData); i += aes.BlockSize {
		c.Decrypt(result[i:], encData[i:])
	}
	return freedomUnPadding(result), nil
}
func EZEncrypt(origData []byte, ezKey string, salt uint64) (_ []byte, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()
	aesKey := MakeMD5Key(ezKey, int64(salt%251))
	return AESCBCEncrypt(origData, aesKey)
}
func EZDecrypt(encData []byte, ezKey string, salt uint64) (_ []byte, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()
	aesKey := MakeMD5Key(ezKey, int64(salt%251))
	return AESCBCDecrypt(encData, aesKey)
}

package ezHash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

func GetMD5Base64String(s string) string {
	md5Maker := md5.New()
	md5Maker.Write([]byte(s))
	rsData := md5Maker.Sum(nil)
	return base64.StdEncoding.EncodeToString(rsData)
}
func GetSHA1Base64String(s string) (string, error) {
	sha1Maker := sha1.New()
	_, err := sha1Maker.Write([]byte(s))
	if err != nil {
		return "", err
	}
	rsData := sha1Maker.Sum(nil)
	return base64.StdEncoding.EncodeToString(rsData), nil
}
func GetSHA256Base64String(s string) (string, error) {
	sha256Maker := sha256.New()
	_, err := sha256Maker.Write([]byte(s))
	if err != nil {
		return "", err
	}
	rsData := sha256Maker.Sum(nil)
	return base64.StdEncoding.EncodeToString(rsData), nil
}

func GetMD5HexString(s string) string {
	md5Maker := md5.New()
	md5Maker.Write([]byte(s))
	rsData := md5Maker.Sum(nil)
	return hex.EncodeToString(rsData)
}
func GetSHA1HexString(s string) (string, error) {
	sha1Maker := sha1.New()
	_, err := sha1Maker.Write([]byte(s))
	if err != nil {
		return "", err
	}
	rsData := sha1Maker.Sum(nil)
	return hex.EncodeToString(rsData), nil
}
func GetSHA256HexString(s string) (string, error) {
	sha256Maker := sha256.New()
	_, err := sha256Maker.Write([]byte(s))
	if err != nil {
		return "", err
	}
	rsData := sha256Maker.Sum(nil)
	return hex.EncodeToString(rsData), nil
}

package ezPasswordEncoder

import (
	"encoding/base64"
	"github.com/Anveena/RoomOfRequirement/ezCrypto"
)

func EncodePassword(origPwd string) (string, error) {
	origPwdData := []byte(origPwd)
	encData, err := ezCrypto.EZEncrypt(&origPwdData, "this code may be not working", 9458)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(*encData), nil
}
func GetPasswordFromEncodedStr(base64Str string) (string, error) {
	encData, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return "", err
	}
	origData, err := ezCrypto.EZDecrypt(&encData, "this code may be not working", 9458)
	if err != nil {
		return "", err
	}
	return string(*origData), nil
}

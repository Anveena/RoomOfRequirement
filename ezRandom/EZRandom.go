package ezRandom

import "math/rand"

const octNumberBytes = "0123456789"
const hexNumberBytes = "0123456789abcdef"
const lowercaseLetterBytes = "qwertyuiopasdfghjklzxcvbnm"
const capitalLetterBytes = "QWERTYUIOPASDFGHJKLZXCVBNM"
const (
	OctNumberOnly            = 1
	HexNumberOnly            = 2
	LowercaseLetterOnly      = 3
	CapitalLetterOnly        = 4
	NumberAndLowercaseLetter = 5
	NumberAndCapitalLetter   = 6
	NumberAndAllLetter       = 7
)

func RandomString(charType int, length int) string {
	var elementArr string
	switch charType {
	case OctNumberOnly:
		elementArr = octNumberBytes
		break
	case HexNumberOnly:
		elementArr = hexNumberBytes
		break
	case LowercaseLetterOnly:
		elementArr = lowercaseLetterBytes
		break
	case CapitalLetterOnly:
		elementArr = capitalLetterBytes
		break
	case NumberAndLowercaseLetter:
		elementArr = octNumberBytes + lowercaseLetterBytes
		break
	case NumberAndCapitalLetter:
		elementArr = octNumberBytes + capitalLetterBytes
		break
	case NumberAndAllLetter:
		elementArr = octNumberBytes + lowercaseLetterBytes + capitalLetterBytes
		break
	default:
		elementArr = octNumberBytes + lowercaseLetterBytes + capitalLetterBytes
		break
	}
	elementsLen := len(elementArr)
	b := make([]byte, length)
	for i := range b {
		b[i] = elementArr[rand.Intn(elementsLen)]
	}
	return string(b)
}

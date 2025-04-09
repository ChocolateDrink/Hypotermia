package utils_crypto_crypt

import (
	"math"
	"strings"
)

func EncryptBasic(text string) string {
	text = strings.ReplaceAll(text, ".", "\\fgfgff55555")
	text = strings.ReplaceAll(text, "-", "\\555ffgfgfff")
	text = strings.ReplaceAll(text, "_", "\\55ffgff5f55")

	encrypted := []rune{}
	shift := getShift(1, 3, 5, 6)

	for _, char := range text {
		newChar := rune(int32(char) + int32(shift))
		encrypted = append(encrypted, newChar)
		shift++
	}

	return string(encrypted)
}

func DecryptBasic(text string) string {
	decrypted := []rune{}
	shift := getShift(1, 3, 5, 6)

	for _, char := range text {
		newChar := rune(int32(char) - int32(shift))
		decrypted = append(decrypted, newChar)
		shift++
	}

	decryptedStr := string(decrypted)

	decryptedStr = strings.ReplaceAll(decryptedStr, "\\fgfgff55555", ".")
	decryptedStr = strings.ReplaceAll(decryptedStr, "\\555ffgfgfff", "-")
	decryptedStr = strings.ReplaceAll(decryptedStr, "\\55ffgff5f55", "_")

	return decryptedStr
}

func getShift(n1 int, n2 int, n3 int, n4 int) int {
	a := (n1 * n1) + (n2 * n2) + (n3 * n3) + (n4 * n4)
	b := n3 * n4 * n1 * n2
	c := getFac(n2) + getFac(n4) + getFac(n3) + getFac(n1)
	d := int(math.Log2(float64(b)) * 10)
	e := int(math.Sin(float64(n2+n1)) * math.Cos(float64(n3+n4)) * 100)
	f := getCom(n1+n2, n3) * getCom(n3+n4, n1)

	return (a * 2) + (b / 10) + (c / 4) - d + e + f
}

func getFac(n int) int {
	if n <= 1 {
		return 1
	}

	return n * getFac(n-1)
}

func getCom(n1 int, n2 int) int {
	return getFac(n1) / (getFac(n2) * getFac(n1-n2))
}

package utils_crypto

import (
	"math"
	"strings"
)

func EncryptBasic(text string) string {
	text = strings.ReplaceAll(text, ".", "\\fgfgff55555")
	text = strings.ReplaceAll(text, "-", "\\555ffgfgfff")
	text = strings.ReplaceAll(text, "_", "\\55ffgff5f55")

	encrypted := []rune{}
	shift := getShift(1, 5, 4, 9)

	for _, char := range text {
		newChar := rune(int32(char) + int32(shift))
		encrypted = append(encrypted, rotate(int(newChar), false))
		shift++
	}

	return permute(string(encrypted))
}

func DecryptBasic(text string) string {
	decrypted := []rune{}
	shift := getShift(1, 5, 4, 9)

	for _, char := range permute(text) {
		newChar := rune(rotate(int(char), true))
		newChar -= rune(shift)

		decrypted = append(decrypted, newChar)
		shift++
	}

	decryptedStr := string(decrypted)

	decryptedStr = strings.ReplaceAll(decryptedStr, "\\fgfgff55555", ".")
	decryptedStr = strings.ReplaceAll(decryptedStr, "\\555ffgfgfff", "-")
	decryptedStr = strings.ReplaceAll(decryptedStr, "\\55ffgff5f55", "_")

	return decryptedStr
}

func permute(text string) string {
	runes := []rune(text)
	shift := getShift(5, 0, 2, 1)

	for i := 0; i < len(runes); i += shift {
		end := min(i+shift, len(runes))

		for j, k := i, end-1; j < k; j, k = j+1, k-1 {
			runes[j], runes[k] = runes[k], runes[j]
		}
	}

	return string(runes)
}

func rotate(char int, back bool) int32 {
	var c int

	if back {
		a := char & 0x00FF
		b := (char >> 8) & 0x00FF
		c = (a << 8) | b
		c = (c - 328 + 0x10FFFF) % 0x10FFFF
		c ^= 0x5A5A
	} else {
		c = char ^ 0x5A5A
		c = (c + 328) % 0x10FFFF
		a := c & 0x00FF
		b := (c >> 8) & 0x00FF
		c = (a << 8) | b
	}

	return int32(c)
}

func getShift(n1 int, n2 int, n3 int, n4 int) int {
	a := (n1) + (n3+n3)*getCom(n3+n4, n1)
	b := n3 * n4 * n2 * int(math.Cos(float64(n3+n4))) * 92
	c := getFac(n2) + getFac(n3) + getFac(n1) + 3
	d := (n4 * n4) + int(math.Log2(float64(b))*10)
	e := int(math.Sin(float64(n2+n1)))*n1 - b
	f := n1*getCom(n1+n2, n3) + getFac(n4) + a
	g := int(math.Erfcinv(float64(a-n3))) - 9
	h := (n1*4 + int(math.Asin(float64(e))) + int(e))

	return (d * 2) + int((float32(f)/10)+(float32(c)))/int(math.Trunc(4)) - d + (g - h) + n2
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

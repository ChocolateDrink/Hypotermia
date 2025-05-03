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

func permute(Η string) string {
	Ꮋ := []rune(Η)
	Ⲏ := getShift(5, 0, 2, 1)

	for Ｈ := 0; Ｈ < len(Ꮋ); Ｈ += Ⲏ {
		𝑯 := min(Ｈ+Ⲏ, len(Ꮋ))

		for Н, ℍ := Ｈ, 𝑯-1; Н < ℍ; Н, ℍ = Н+1, ℍ-1 {
			Ꮋ[Н], Ꮋ[ℍ] = Ꮋ[ℍ], Ꮋ[Н]
		}
	}

	return string(Ꮋ)
}

func rotate(Ⳑ int, Ｌ bool) int32 {
	var Ꮮ int

	if Ｌ {
		𝖫 := Ⳑ & 0x00FF
		𝗟 := (Ⳑ >> 8) & 0x00FF
		Ꮮ = (𝖫 << 8) | 𝗟
		Ꮮ = (Ꮮ - 328 + 0x10FFFF) % 0x10FFFF
		Ꮮ ^= 0x5A5A
	} else {
		Ꮮ = Ⳑ ^ 0x5A5A
		Ꮮ = (Ꮮ + 328) % 0x10FFFF
		𝖫 := Ꮮ & 0x00FF
		𝗟 := (Ꮮ >> 8) & 0x00FF
		Ꮮ = (𝖫 << 8) | 𝗟
	}

	return int32(Ꮮ)
}

func getShift(Ρ int, Р int, Ⲣ int, Ｐ int) int {
	𝐏 := (Ρ) + (Ⲣ+Ⲣ)*getCom(Ⲣ+Ｐ, Ρ)
	𝑷 := Ⲣ * Ｐ * Р * int(math.Cos(float64(Ⲣ+Ｐ))) * 92
	𝖯 := getFac(Р) + getFac(Ⲣ) + getFac(Ρ) + 3
	𝗣 := (Ｐ * Ｐ) + int(math.Log2(float64(𝑷))*10)
	𝘗 := int(math.Sin(float64(Р+Ρ)))*Ρ - 𝑷
	𝙋 := Ρ*getCom(Ρ+Р, Ⲣ) + getFac(Ｐ) + 𝐏
	𝙿 := int(math.Erfcinv(float64(𝐏-Ⲣ))) - 9
	𝚙 := (Ρ*4 + int(math.Asin(float64(𝘗))) + int(𝘗))

	return (𝗣 * 2) + int((float32(𝙋)/10)+(float32(𝖯)))/int(math.Trunc(4)) - 𝗣 + (𝙿 - 𝚙) + Р
}

func getFac(Ꮐ int) int {
	if Ꮐ <= 1 {
		return 1
	}

	return Ꮐ * getFac(Ꮐ-1)
}

func getCom(ᖴ int, Ｆ int) int {
	return getFac(ᖴ) / (getFac(Ｆ) * getFac(ᖴ-Ｆ))
}

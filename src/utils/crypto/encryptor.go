package main

import (
	"fmt"

	"Hypotermia/src/utils/crypto/crypt"
)

func main() {
	var token string

	fmt.Print("Data: ")
	_, err := fmt.Scanf("%s", &token)
	if err != nil {
		fmt.Println("Error reading:", err)
		return
	}

	fmt.Printf("Encrypting data: \"%s\"\n", token)

	encrypted := utils_crypto_crypt.EncryptBasic(token)
	fmt.Println("Encrypted data:", encrypted)

	decrypted := utils_crypto_crypt.DecryptBasic(encrypted)
	fmt.Println("Decrypted data:", decrypted)

	fmt.Printf("Does it match: %v", token == decrypted)
}

package main

import (
	"fmt"

	"Hypothermia/src/utils/crypto"
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

	encrypted := utils_crypto.EncryptBasic(token)
	fmt.Println("Encrypted data:", encrypted)

	decrypted := utils_crypto.DecryptBasic(encrypted)
	fmt.Println("Decrypted data:", decrypted)

	fmt.Printf("Does it match: %v", token == decrypted)
}

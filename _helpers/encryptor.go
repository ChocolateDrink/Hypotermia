package main

import (
	"fmt"

	"Hypothermia/src/utils/crypto"
)

func main() {
	var data string

	fmt.Print("Data: ")
	_, err := fmt.Scanf("%s", &data)
	if err != nil {
		fmt.Println("Error reading:", err)
		return
	}

	fmt.Printf("Encrypting data: \"%s\"\n", data)

	encrypted := utils_crypto.EncryptBasic(data)
	fmt.Println("Encrypted data:", encrypted)

	decrypted := utils_crypto.DecryptBasic(encrypted)
	fmt.Println("Decrypted data:", decrypted)

	fmt.Printf("Does it match: %v\n", data == decrypted)
}

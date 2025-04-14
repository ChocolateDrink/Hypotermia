package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"os"
	"slices"
	"unsafe"

	"golang.org/x/sys/windows"
)

func GetPassword(path string) ([]byte, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var localState map[string]any

	err = json.Unmarshal(data, &localState)
	if err != nil {
		return nil, err
	}

	osCrypt, ok := localState["os_crypt"].(map[string]any)
	if !ok {
		return nil, err
	}

	encryptedKey, ok := osCrypt["encrypted_key"].(string)
	if !ok {
		return nil, err
	}

	decodedKey, err := base64.StdEncoding.DecodeString(encryptedKey)
	if err != nil {
		return nil, err
	}

	decodedKey = decodedKey[5:]

	key, err := decryptPassword(decodedKey)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func DecryptWithPassword(data []byte, password []byte) (string, error) {
	iv := data[3:15]
	encoded := data[15:]

	block, err := aes.NewCipher(password)
	if err != nil {
		return "", err
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	decoded, err := aead.Open(nil, iv, encoded, nil)
	if err != nil {
		return "", err
	}

	return string(decoded), nil
}

func decryptPassword(data []byte) ([]byte, error) {
	var in windows.DataBlob
	var out windows.DataBlob

	in.Size = uint32(len(data))
	in.Data = &data[0]

	err := windows.CryptUnprotectData(&in, nil, nil, 0, nil, 0, &out)
	if err != nil {
		return nil, err
	}

	defer windows.LocalFree(windows.Handle(unsafe.Pointer(out.Data)))

	slice := unsafe.Slice(out.Data, out.Size)
	return slices.Clone(slice), nil
}

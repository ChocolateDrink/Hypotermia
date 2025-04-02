package utils

import (
	"fmt"

	"golang.org/x/sys/windows/registry"
)

func SetRegistry(path string, key string, value string) error {
	if len(key) == 0 {
		return fmt.Errorf("key is empty")
	}

	if len(value) == 0 {
		return fmt.Errorf("value is empty")
	}

	handle, err := registry.OpenKey(
		registry.CURRENT_USER,
		path, registry.ALL_ACCESS,
	)

	if err != nil {
		return err
	}

	defer handle.Close()

	err = handle.SetStringValue(key, value)
	if err != nil {
		return err
	}

	return nil
}

func GetRegistry(path string, key string, value string) (string, error) {
	if len(key) == 0 {
		return "", fmt.Errorf("key is empty")
	}

	if len(value) == 0 {
		return "", fmt.Errorf("value is empty")
	}

	handle, err := registry.OpenKey(
		registry.CURRENT_USER,
		path, registry.ALL_ACCESS,
	)

	if err != nil {
		return "", err
	}

	defer handle.Close()

	val, _, err := handle.GetStringValue(key)
	if err != nil {
		return "", err
	}

	return val, nil
}

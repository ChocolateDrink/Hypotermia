package utils

import (
	"fmt"

	"golang.org/x/sys/windows/registry"
)

func MakeRegistryKey(path string) error {
	if len(path) == 0 {
		return fmt.Errorf("path is empty")
	}

	key, _, err := registry.CreateKey(
		registry.CURRENT_USER,
		path, registry.ALL_ACCESS,
	)

	if err != nil {
		return err
	}

	defer key.Close()

	return nil
}

func DelRegistryKey(path string) error {
	if len(path) == 0 {
		return fmt.Errorf("path is empty")
	}

	err := registry.DeleteKey(registry.CURRENT_USER, path)
	if err != nil {
		return err
	}

	return nil
}

func SetRegistryVal(path string, key string, value string) error {
	if len(path) == 0 {
		return fmt.Errorf("path is empty")
	}

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

func GetRegistryVal(path string, key string) (string, error) {
	if len(path) == 0 {
		return "", fmt.Errorf("path is empty")
	}

	if len(key) == 0 {
		return "", fmt.Errorf("key is empty")
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

func DelRegistryVal(path string, key string) error {
	if len(path) == 0 {
		return fmt.Errorf("path is empty")
	}

	if len(key) == 0 {
		return fmt.Errorf("key is empty")
	}

	handle, err := registry.OpenKey(
		registry.CURRENT_USER,
		path, registry.ALL_ACCESS,
	)

	if err != nil {
		return err
	}

	defer handle.Close()

	err = handle.DeleteValue(key)
	if err != nil {
		return err
	}

	return nil
}

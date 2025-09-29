package config

import (
	"errors"
	"os"
	"path/filepath"
)

func DirExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, err
	} else if !info.IsDir() {
		return false, errors.New("path exists but is not a directory")
	}
	return true, nil
}

func GetUserHomeDirectory() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, AppName), nil
}

func GetSecretFilePath() (string, error) {
	homeDir, err := GetUserHomeDirectory()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, SecretFileName), nil
}

func GetAPIKey() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	path := filepath.Join(homeDir, AppName, SecretFileName)

	apiKey, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(apiKey), nil
}

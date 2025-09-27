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

func GetHomeDirectory() (string, error) {
	homeDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, AppName), nil
}

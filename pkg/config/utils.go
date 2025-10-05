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

func GetTreeFromDir(root string) ([]string, error) {
	var files []string

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			relPath, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			files = append(files, relPath)
		}
		return nil
	})

	return files, err
}

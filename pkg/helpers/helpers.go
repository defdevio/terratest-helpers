package helpers

import (
	"os"
	"path"
	"testing"
)

// Creates a file if does not already exist on the specified path
func CreateFile(path string, content string) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			bytes := []byte(content)
			err = os.WriteFile(path, bytes, 0644)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Creates the required provider file on the system
func CreateAzureProviderFile(providerFilePath string, t *testing.T) error {
	providerContent := `
provider "azurerm" {
	features {}
}`

	err := CreateFile(providerFilePath, providerContent)
	return err
}

// Cleans up the common files created by terraform during init, plan, and apply
func CleanUpTestFiles(t *testing.T, files []string, workDir string) error {
	for _, file := range files {
		filePath := path.Join(workDir, file)
		err := os.RemoveAll(filePath)
		if err != nil {
			return err
		}
	}

	return nil
}

package helpers

import (
	"context"
	"os"
	"path"
	"testing"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/resources/mgmt/resources"
	"github.com/gruntwork-io/terratest/modules/azure"
)

var (
	ctx = context.Background()
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

func CreateAzureResourceGroup(t *testing.T, subscriptionID string, resourceGroup string, location *string) error {
	resourceGroupClient, err := azure.CreateResourceGroupClientE(subscriptionID)
	if err != nil {
		t.Fatal(err)
	}

	// Create the resource group using the resourceGroupClient
	resp, err := resourceGroupClient.CreateOrUpdate(ctx, resourceGroup, resources.Group{
		Location: location,
	})

	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode == 201 {
		t.Logf("Created resource group '%s'", *resp.Name)
	}

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

func DeleteAzureResourceGroup(t *testing.T, subscriptionID string, resourceGroup string) error {
	resourceGroupClient, err := azure.CreateResourceGroupClientE(subscriptionID)
	if err != nil {
		t.Fatal(err)
	}

	_, err = resourceGroupClient.Delete(ctx, resourceGroup)
	if err != nil {
		t.Fatal(err)
	}

	return err
}

// Gets the keys of a terraform map and returns a slice of strings
func GetMapKeys(t *testing.T, terraformMap map[string]any) []string {
	keys := make([]string, len(terraformMap))
	i := 0
	for key := range terraformMap {
		keys[i] = key
		i++
	}

	return keys
}

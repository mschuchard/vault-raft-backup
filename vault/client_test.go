package vault

import (
	"testing"
)

func TestNewVaultConfig(test *testing.T) {
	// test with defaults
	vaultConfigDefault, err := NewVaultConfig()
	if err != nil {
		test.Error("vault config constructor failed default initialization")
		test.Error(err)
	}

	if vaultConfigDefault.address != "http://127.0.0.1:8200" || vaultConfigDefault.insecure != false || vaultConfigDefault.engine != awsIam || len(vaultConfigDefault.token) != 0 || vaultConfigDefault.awsMountPath != "aws" || len(vaultConfigDefault.awsRole) != 0 || len(vaultConfigDefault.snapshotPath) != 0 {
		test.Error("vault config default constructor did not initialize with expected values")
		test.Errorf("address expected: http://127.0.0.1:8200, actual: %s", vaultConfigDefault.address)
		test.Errorf("insecure expected: false, actual: %t", vaultConfigDefault.insecure)
		test.Errorf("engine expected: aws, actual: %v", vaultConfigDefault.engine)
		test.Errorf("token expected: (empty), actual: %s", vaultConfigDefault.token)
		test.Errorf("aws mount path expected: aws, actual: %s", vaultConfigDefault.awsMountPath)
		test.Errorf("aws role expected: (empty), actual: %s", vaultConfigDefault.awsRole)
		test.Errorf("snapshot path expected: (empty), actual: %s", vaultConfigDefault.snapshotPath)
	}

	// setup env for custom constructor inputs

	// test errors
}

func TestNewVaultClient(test *testing.T) {

}

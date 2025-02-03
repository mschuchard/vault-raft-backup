package vault

import (
	"os"
	"strings"
	"testing"

	"github.com/mschuchard/vault-raft-backup/util"
)

var (
	backupVaultConfig     = &util.VaultConfig{}
	expectedDefaultConfig = vaultConfig{
		address:      "http://127.0.0.1:8200",
		insecure:     true,
		engine:       awsIam,
		token:        "",
		awsMountPath: "aws",
		awsRole:      "",
	}
	backupVaultTokenConfig = &util.VaultConfig{
		Address:      "https://127.0.0.1:8234",
		Engine:       "token",
		Token:        util.VaultToken,
		SnapshotPath: "/tmp/my_vault.backup",
	}
	expectedTokenConfig = vaultConfig{
		address:      "https://127.0.0.1:8234",
		insecure:     false,
		engine:       vaultToken,
		token:        util.VaultToken,
		awsMountPath: "",
		awsRole:      "",
	}
	backupVaultAWSConfig = &util.VaultConfig{
		Address:      "https://127.0.0.1:8234",
		Insecure:     true,
		AWSMountPath: "gcp",
		AWSRole:      "my_role",
	}
	expectedAWSConfig = vaultConfig{
		address:      "https://127.0.0.1:8234",
		insecure:     true,
		engine:       awsIam,
		token:        "",
		awsMountPath: "gcp",
		awsRole:      "my_role",
	}
)

func TestNewVaultConfig(test *testing.T) {
	// test with defaults
	vaultConfigDefault, err := NewVaultConfig(backupVaultConfig)
	if err != nil {
		test.Error("vault config constructor failed default initialization")
		test.Error(err)
	}

	if *vaultConfigDefault != expectedDefaultConfig {
		test.Error("vault config default constructor did not initialize with expected values")
		test.Errorf("expected vault config values: %v", expectedDefaultConfig)
		test.Errorf("actual vault config values: %v", *vaultConfigDefault)
	}

	// test with token
	vaultConfigToken, err := NewVaultConfig(backupVaultTokenConfig)
	if err != nil {
		test.Error("vault config constructor failed custom token initialization")
		test.Error(err)
	}

	if *vaultConfigToken != expectedTokenConfig {
		test.Error("vault config token constructor did not initialize with expected values")
		test.Errorf("expected vault config values: %v", expectedTokenConfig)
		test.Errorf("actual vault config values: %v", *vaultConfigToken)
	}

	// test with aws
	vaultConfigAWS, err := NewVaultConfig(backupVaultAWSConfig)
	if err != nil {
		test.Error("vault config constructor custom failed aws initialization")
		test.Error(err)
	}

	if *vaultConfigAWS != expectedAWSConfig {
		test.Error("vault config aws constructor did not initialize with expected values")
		test.Errorf("expected vault config values: %v", expectedAWSConfig)
		test.Errorf("actual vault config values: %v", *vaultConfigAWS)
	}

	// test errors in reverse validation order
	backupVaultConfig.Token = "1234"
	if _, err = NewVaultConfig(backupVaultConfig); err == nil || err.Error() != "invalid vault token" {
		test.Errorf("expected error: invalid vault token, actual: %v", err)
	}

	backupVaultConfig.Engine = "kubernetes"
	if _, err = NewVaultConfig(backupVaultConfig); err == nil || err.Error() != "invalid Vault authentication engine" {
		test.Errorf("expected error: invalid Vault authentication engine, actual: %v", err)
	}

	backupVaultConfig.Engine = ""
	backupVaultConfig.AWSMountPath = "azure"
	if _, err = NewVaultConfig(backupVaultConfig); err == nil || err.Error() != "unable to deduce authentication engine" {
		test.Errorf("expected error: unable to deduce authentication engine, actual: %v", err)
	}
	backupVaultConfig.Token = ""

	backupVaultConfig.Address = "file:///foo"
	if _, err = NewVaultConfig(backupVaultConfig); err == nil || err.Error() != "invalid Vault server address" {
		test.Error("expected error for invalid Vault server address, but none was returned")
	}
	os.Setenv("VAULT_ADDR", "")
}

func TestNewVaultClient(test *testing.T) {
	// test client with aws iam auth
	expectedAWSConfig.address = "http://127.0.0.1:8200"
	if _, err := NewVaultClient(&expectedAWSConfig); err == nil || !strings.Contains(err.Error(), "NoCredentialProviders: no valid providers in chain") {
		test.Errorf("expected error (contains): NoCredentialProviders: no valid providers in chain, actual: %v", err)
	}

	// test client with token auth
	expectedTokenConfig.address = "http://127.0.0.1:8200"
	if _, err := NewVaultClient(&expectedTokenConfig); err != nil {
		test.Error("client failed to initialize with basic token auth config information")
		test.Error(err)
	}
}

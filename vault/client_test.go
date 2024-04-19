package vault

import (
	"os"
	"strings"
	"testing"

	"github.com/mitodl/vault-raft-backup/util"
)

var (
	expectedDefaultConfig = vaultConfig{
		address:      "http://127.0.0.1:8200",
		insecure:     true,
		engine:       awsIam,
		token:        "",
		awsMountPath: "aws",
		awsRole:      "",
		snapshotPath: "/tmp/vault.bak",
	}
	expectedTokenConfig = vaultConfig{
		address:      "https://127.0.0.1:8234",
		insecure:     false,
		engine:       vaultToken,
		token:        util.VaultToken,
		awsMountPath: "",
		awsRole:      "",
		snapshotPath: "/tmp/my_vault.backup",
	}
	expectedAWSConfig = vaultConfig{
		address:      "https://127.0.0.1:8234",
		insecure:     true,
		engine:       awsIam,
		token:        "",
		awsMountPath: "gcp",
		awsRole:      "my_role",
		snapshotPath: "/tmp/vault.bak",
	}
)

func TestNewVaultConfig(test *testing.T) {
	// test with defaults
	vaultConfigDefault, err := NewVaultConfig()
	if err != nil {
		test.Error("vault config constructor failed default initialization")
		test.Error(err)
	}

	if *vaultConfigDefault != expectedDefaultConfig {
		test.Error("vault config default constructor did not initialize with expected values")
		test.Errorf("expected vault config values: %v", expectedDefaultConfig)
		test.Errorf("actual vault config values: %v", *vaultConfigDefault)
	}

	// setup env for custom constructor inputs with token
	os.Setenv("VAULT_ADDR", "https://127.0.0.1:8234")
	os.Setenv("VAULT_AUTH_ENGINE", "token")
	os.Setenv("VAULT_TOKEN", util.VaultToken)
	os.Setenv("VAULT_SNAPSHOT_PATH", "/tmp/my_vault.backup")
	vaultConfigToken, err := NewVaultConfig()
	if err != nil {
		test.Error("vault config constructor failed custom token initialization")
		test.Error(err)
	}

	if *vaultConfigToken != expectedTokenConfig {
		test.Error("vault config token constructor did not initialize with expected values")
		test.Errorf("expected vault config values: %v", expectedTokenConfig)
		test.Errorf("actual vault config values: %v", *vaultConfigToken)
	}
	os.Setenv("VAULT_TOKEN", "")
	os.Setenv("VAULT_AUTH_ENGINE", "")
	os.Setenv("VAULT_SNAPSHOT_PATH", "")

	// setup env for custom constructor inputs with aws
	os.Setenv("VAULT_SKIP_VERIFY", "true")
	os.Setenv("VAULT_AWS_MOUNT", "gcp")
	os.Setenv("VAULT_AWS_ROLE", "my_role")
	vaultConfigAWS, err := NewVaultConfig()
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
	os.Setenv("VAULT_TOKEN", "1234")
	os.Setenv("VAULT_AUTH_ENGINE", "token")
	if _, err = NewVaultConfig(); err == nil || err.Error() != "invalid vault token" {
		test.Errorf("expected error: invalid vault token, actual: %v", err)
	}

	os.Setenv("VAULT_AUTH_ENGINE", "kubernetes")
	if _, err = NewVaultConfig(); err == nil || err.Error() != "invalid Vault authentication engine" {
		test.Errorf("expected error: invalid Vault authentication engine, actual: %v", err)
	}

	os.Setenv("VAULT_AUTH_ENGINE", "")
	if _, err = NewVaultConfig(); err == nil || err.Error() != "unable to deduce authentication engine" {
		test.Errorf("expected error: unable to deduce authentication engine, actual: %v", err)
	}
	os.Setenv("VAULT_TOKEN", "")

	os.Setenv("VAULT_SKIP_VERIFY", "not a boolean")
	if _, err = NewVaultConfig(); err == nil || err.Error() != "invalid VAULT_SKIP_VERIFY value" {
		test.Errorf("expected error: invalid VAULT_SKIP_VERIFY value, actual: %v", err)
	}
	os.Setenv("VAULT_SKIP_VERIFY", "")

	os.Setenv("VAULT_ADDR", "file:///foo")
	if _, err = NewVaultConfig(); err == nil || err.Error() != "invalid Vault server address" {
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

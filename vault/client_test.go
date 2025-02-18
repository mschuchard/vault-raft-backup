package vault

import (
	"strings"
	"testing"

	"github.com/mschuchard/vault-raft-backup/util"
)

var (
	backupVaultConfig      = &util.VaultConfig{}
	backupVaultTokenConfig = &util.VaultConfig{
		Address:      "https://127.0.0.1:8234",
		Engine:       "token",
		Token:        util.VaultToken,
		SnapshotPath: "/tmp/my_vault.backup",
	}
	backupVaultAWSConfig = &util.VaultConfig{
		Address:      "https://127.0.0.1:8234",
		Insecure:     true,
		AWSMountPath: "gcp",
		AWSRole:      "my_role",
	}
)

func TestNewVaultClient(test *testing.T) {
	// test with defaults
	_, err := NewVaultClient(backupVaultConfig)
	if err == nil || !strings.Contains(err.Error(), "NoCredentialProviders: no valid providers in chain") {
		test.Errorf("expected error (contains): NoCredentialProviders: no valid providers in chain, actual: %v", err)
	}

	/*if vaultClientDefault.Address() != "http://127.0.0.1:8200" || len(vaultClientDefault.Token()) > 0 {
		test.Error("vault client default constructor did not initialize with expected values")
		test.Error("expected default vault client values: http://127.0.0.1:8200 and empty string")
		test.Errorf("actual vault client values: %v", *vaultClientDefault)
	}*/

	// test with token
	/*vaultClientToken, err := NewVaultClient(backupVaultTokenConfig)
	if err != nil {
		test.Error("client failed to initialize with basic token auth config information")
		test.Error(err)
	}

	if vaultClientToken.Address() != "https://127.0.0.1:8234" || vaultClientToken.Token() != util.VaultToken {
		test.Error("vault client token constructor did not initialize with expected values")
		test.Errorf("expected vault client values: %s, %s", backupVaultTokenConfig.Address, backupVaultTokenConfig.Token)
		test.Errorf("actual vault client values: %v", *vaultClientToken)
	}

	// test with aws
	_, err = NewVaultClient(backupVaultAWSConfig)
	if err == nil || !strings.Contains(err.Error(), "NoCredentialProviders: no valid providers in chain") {
		test.Errorf("expected error (contains): NoCredentialProviders: no valid providers in chain, actual: %v", err)
	}

	if vaultClientAWS.Address() != "https://127.0.0.1:8234" || len(vaultClientAWS.Token()) > 0 {
		test.Error("vault client aws constructor did not initialize with expected values")
		test.Errorf("expected vault client values: %s, %s", backupVaultAWSConfig.Address, backupVaultAWSConfig.Token)
		test.Errorf("actual vault client values: %v", *vaultClientAWS)
	}*/

	// test errors in reverse validation order
	backupVaultConfig.Token = "1234"
	if _, err = NewVaultClient(backupVaultConfig); err == nil || err.Error() != "invalid vault token" {
		test.Errorf("expected error: invalid vault token, actual: %v", err)
	}

	backupVaultConfig.Engine = "kubernetes"
	if _, err = NewVaultClient(backupVaultConfig); err == nil || err.Error() != "invalid Vault authentication engine" {
		test.Errorf("expected error: invalid Vault authentication engine, actual: %v", err)
	}

	backupVaultConfig.Engine = ""
	backupVaultConfig.AWSMountPath = "azure"
	if _, err = NewVaultClient(backupVaultConfig); err == nil || err.Error() != "unable to deduce authentication engine" {
		test.Errorf("expected error: unable to deduce authentication engine, actual: %v", err)
	}
	backupVaultConfig.Token = ""

	backupVaultConfig.Address = "file:///foo"
	if _, err = NewVaultClient(backupVaultConfig); err == nil || err.Error() != "invalid Vault server address" {
		test.Error("expected error for invalid Vault server address, but none was returned")
	}
}

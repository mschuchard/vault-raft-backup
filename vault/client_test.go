package vault

import (
	"errors"
	"os"
	"testing"

	"github.com/mitodl/vault-raft-backup/util"
)

func TestNewVaultConfig(test *testing.T) {
	// test with defaults
	vaultConfigDefault, err := NewVaultConfig()
	if err != nil {
		test.Error("vault config constructor failed default initialization")
		test.Error(err)
	}

	if vaultConfigDefault.address != "http://127.0.0.1:8200" || vaultConfigDefault.insecure == false || vaultConfigDefault.engine != awsIam || len(vaultConfigDefault.token) != 0 || vaultConfigDefault.awsMountPath != "aws" || len(vaultConfigDefault.awsRole) != 0 || len(vaultConfigDefault.snapshotPath) != 0 {
		test.Error("vault config default constructor did not initialize with expected values")
		test.Errorf("address expected: http://127.0.0.1:8200, actual: %s", vaultConfigDefault.address)
		test.Errorf("insecure expected: true, actual: %t", vaultConfigDefault.insecure)
		test.Errorf("engine expected: aws, actual: %v", vaultConfigDefault.engine)
		test.Errorf("token expected: (empty), actual: %s", vaultConfigDefault.token)
		test.Errorf("aws mount path expected: aws, actual: %s", vaultConfigDefault.awsMountPath)
		test.Errorf("aws role expected: (empty), actual: %s", vaultConfigDefault.awsRole)
		test.Errorf("snapshot path expected: (empty), actual: %s", vaultConfigDefault.snapshotPath)
	}

	// setup env for custom constructor inputs with token
	os.Setenv("VAULT_ADDR", "https://127.0.0.1:8234")
	os.Setenv("VAULT_AUTH_ENGINE", "token")
	os.Setenv("VAULT_TOKEN", util.VaultToken)
	vaultConfigToken, err := NewVaultConfig()
	if err != nil {
		test.Error("vault config constructor failed custom token initialization")
		test.Error(err)
	}

	if vaultConfigToken.address != "https://127.0.0.1:8234" || vaultConfigToken.insecure == true || vaultConfigToken.engine != vaultToken || vaultConfigToken.token != util.VaultToken || len(vaultConfigToken.awsMountPath) != 0 || len(vaultConfigToken.awsRole) != 0 || len(vaultConfigToken.snapshotPath) != 0 {
		test.Error("vault config token constructor did not initialize with expected values")
		test.Errorf("address expected: https://127.0.0.1:8234, actual: %s", vaultConfigToken.address)
		test.Errorf("insecure expected: false, actual: %t", vaultConfigToken.insecure)
		test.Errorf("engine expected: token, actual: %v", vaultConfigToken.engine)
		test.Errorf("token expected: %s, actual: %s", util.VaultToken, vaultConfigToken.token)
		test.Errorf("aws mount path expected: (empty), actual: %s", vaultConfigToken.awsMountPath)
		test.Errorf("aws role expected: (empty), actual: %s", vaultConfigToken.awsRole)
		test.Errorf("snapshot path expected: (empty), actual: %s", vaultConfigToken.snapshotPath)
	}

	// setup env for custom constructor inputs with aws
	os.Setenv("VAULT_SKIP_VERIFY", "true")
	os.Setenv("VAULT_AUTH_ENGINE", "aws")
	os.Setenv("VAULT_AWS_MOUNT", "gcp")
	os.Setenv("VAULT_AWS_ROLE", "my_role")
	vaultConfigAWS, err := NewVaultConfig()
	if err != nil {
		test.Error("vault config constructor custom failed aws initialization")
		test.Error(err)
	}

	if vaultConfigAWS.address != "https://127.0.0.1:8234" || vaultConfigAWS.insecure == false || vaultConfigAWS.engine != awsIam || vaultConfigAWS.token != util.VaultToken || vaultConfigAWS.awsMountPath != "gcp" || vaultConfigAWS.awsRole != "my_role" || len(vaultConfigAWS.snapshotPath) != 0 {
		test.Error("vault config aws constructor did not initialize with expected values")
		test.Errorf("address expected: https://127.0.0.1:8234, actual: %s", vaultConfigAWS.address)
		test.Errorf("insecure expected: true, actual: %t", vaultConfigAWS.insecure)
		test.Errorf("engine expected: aws, actual: %v", vaultConfigAWS.engine)
		test.Errorf("token expected: %s, actual: %s", util.VaultToken, vaultConfigAWS.token)
		test.Errorf("aws mount path expected: gcp, actual: %s", vaultConfigAWS.awsMountPath)
		test.Errorf("aws role expected: my_role, actual: %s", vaultConfigAWS.awsRole)
		test.Errorf("snapshot path expected: (empty), actual: %s", vaultConfigAWS.snapshotPath)
	}

	// test errors in reverse validation order
	os.Setenv("VAULT_TOKEN", "1234")
	if _, err = NewVaultConfig(); err != errors.New("invalid vault token") {
		test.Errorf("expected error: invalid vault token, actual: %v", err)
	}

	os.Setenv("VAULT_AUTH_ENGINE", "kubernetes")
	if _, err = NewVaultConfig(); err != errors.New("invalid Vault authentication engine") {
		test.Errorf("expected error: invalid Vault authentication engine, actual: %v", err)
	}

	os.Setenv("VAULT_AUTH_ENGINE", "")
	if _, err = NewVaultConfig(); err != errors.New("unable to deduce authentication engine") {
		test.Errorf("expected error: unable to deduce authentication engine, actual: %v", err)
	}

	os.Setenv("VAULT_SKIP_VERFIY", "not a boolean")
	if _, err = NewVaultConfig(); err != errors.New("invalid VAULT_SKIP_VERIFY value") {
		test.Errorf("expected error: invalid VAULT_SKIP_VERIFY value, actual: %v", err)
	}

	os.Setenv("VAULT_ADDR", "file:///foo")
	if _, err = NewVaultConfig(); err != nil {
		test.Error("expected error for invalid Vault server address, but none was returned")
	}
}

func TestNewVaultClient(test *testing.T) {
	// test client with aws iam auth
	vaultAWSConfig, _ := NewVaultConfig()
	if _, err := NewVaultClient(vaultAWSConfig); err != errors.New("unable to login to AWS IAM auth method") {
		test.Errorf("expected error: unable to login to AWS IAM auth method, actual: %v", err)
	}

	// test client with token auth
	os.Setenv("VAULT_ADDR", "http://127.0.0.1:8234")
	os.Setenv("VAULT_TOKEN", util.VaultToken)
	vaultTokenConfig, _ := NewVaultConfig()
	if _, err := NewVaultClient(vaultTokenConfig); err != nil {
		test.Error("client failed to initialize with basic token auth config information")
	}
}

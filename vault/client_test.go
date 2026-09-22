package vault

import (
	"strings"
	"testing"

	"github.com/mschuchard/vault-raft-backup/enum"
	"github.com/mschuchard/vault-raft-backup/util"
)

var (
	basicBackupConfig = &util.VaultConfig{
		Address: util.VaultAddress,
		Engine:  enum.VaultToken,
		Token:   util.VaultToken,
	}
	awsBackupConfig = &util.VaultConfig{
		Address:   util.VaultAddress,
		Engine:    enum.AWSIAM,
		VaultRole: "myIAMRole",
	}
	azBackupConfig = &util.VaultConfig{
		Address:    util.VaultAddress,
		Engine:     enum.AzureIMDS,
		VaultRole:  "myAzureRole",
		AzResource: "https://management.azure.com/",
	}
	kubeBackupConfig = &util.VaultConfig{
		Address:   util.VaultAddress,
		Engine:    enum.KubernetesSA,
		VaultRole: "mySARole",
	}
	approleBackupConfig = &util.VaultConfig{
		Address: util.VaultAddress,
		Engine:  enum.AppRole,
	}
)

// test client constructor
func TestNewVaultClient(test *testing.T) {
	basicClient, err := NewVaultClient(basicBackupConfig)
	if err != nil {
		test.Error("authenticating a vault client with a basic token config errored")
		test.Error(err)
	}
	if basicClient.Address() != basicBackupConfig.Address || basicClient.Token() != basicBackupConfig.Token {
		test.Error("the authenticated Vault client return failed basic validation")
		test.Errorf("expected Vault token: %s, actual: %s", basicBackupConfig.Token, basicClient.Token())
		test.Errorf("expected Vault address: %s, actual: %s", basicBackupConfig.Address, basicClient.Address())
	}

	// test errors
	invalidServerConfig := &util.VaultConfig{Address: "https//:foo.com"}
	if _, err := NewVaultClient(invalidServerConfig); err == nil || err.Error() != "parse \"https//:foo.com\": invalid URI for request" {
		test.Errorf("expected error: parse \"https//:foo.com\": invalid URI for request, actual: %s", err)
	}
}

// test client auth
func TestAuthClient(test *testing.T) {
	if err := authClient(awsBackupConfig, util.VaultClient); err == nil || !strings.Contains(err.Error(), "NoCredentialProviders: no valid providers in chain") {
		test.Error("authenticating a vault client with aws did not error in the expected manner")
		test.Errorf("expected error (contains): NoCredentialProviders: no valid providers in chain, actual: %v", err)
	}

	awsBackupConfig.VaultRole = ""
	if err := authClient(awsBackupConfig, util.VaultClient); err == nil || !strings.Contains(err.Error(), "NoCredentialProviders: no valid providers in chain") {
		test.Error("authenticating a vault client with aws did not error in the expected manner")
		test.Errorf("expected error (contains): NoCredentialProviders: no valid providers in chain, actual: %v", err)
	}

	if err := authClient(kubeBackupConfig, util.VaultClient); err == nil || !strings.Contains(err.Error(), "error reading service account token from default location") {
		test.Error("authenticating a vault client with kubernetes did not error in the expected manner")
		test.Errorf("expected error (contains): error reading service account token from default location, actual: %v", err)
	}

	if err := authClient(azBackupConfig, util.VaultClient); err == nil || !strings.Contains(err.Error(), "error calling Azure token endpoint") {
		test.Error("authenticating a vault client with azure did not error in the expected manner")
		test.Errorf("expected error (contains): error calling Azure token endpoint, actual: %v", err)
	}

	approleBackupConfig.VaultRole = util.RoleID
	approleBackupConfig.SecretID = util.SecretID
	if err := authClient(approleBackupConfig, util.VaultClient); err != nil {
		test.Error("authenticating a vault client with approle config errored")
		test.Error(err)
	}

	// reset client auth to prep for next test
	util.VaultClient.SetToken(util.VaultToken)

	// retrieve a wrapped secret id for testing approle auth in "pull" mode
	util.VaultClient.SetWrappingLookupFunc(func(operation, path string) string {
		if path == "auth/approle/role/myAppRole/secret-id" {
			return "60s"
		}
		return ""
	})
	wrappedSecretID, err := util.VaultClient.Logical().Write("auth/approle/role/myAppRole/secret-id", nil)
	// reset the wrapping lookup func immediately so it does not affect subsequent requests
	util.VaultClient.SetWrappingLookupFunc(nil)
	if err != nil {
		test.Error("failed to retrieve wrapped secret ID for approle pull auth")
		test.Error(err)
	}
	// access wrapping token from wrapped secret id and assign to source config
	if wrappedSecretID == nil || wrappedSecretID.WrapInfo == nil || len(wrappedSecretID.WrapInfo.Token) == 0 {
		test.Error("the secret id write did not return a wrapped response")
	}
	approleBackupConfig.WrapToken = wrappedSecretID.WrapInfo.Token

	if err := authClient(approleBackupConfig, util.VaultClient); err != nil {
		test.Error("authenticating a vault client with approle pull (wrapping token) config errored")
		test.Error(err)
	}

	// this needs to be last to ensure the client is authenticated with a root token for all other tests
	if err := authClient(basicBackupConfig, util.VaultClient); err != nil {
		test.Error("authenticating a vault client with a basic token config errored")
		test.Error(err)
	}

	// test errors
	invalidAuth := &util.VaultConfig{Engine: "does not exist"}
	if err := authClient(invalidAuth, util.VaultClient); err == nil || err.Error() != "invalid authengine enum" {
		test.Errorf("expected error: invalid authengine enum, actual: %s", err)
	}

	invalidToken := &util.VaultConfig{Engine: enum.VaultToken, Token: "foobarbaz123!"}
	if err := authClient(invalidToken, util.VaultClient); err == nil || err.Error() != "invalid vault token" {
		test.Errorf("expected error: invalid vault token, actual: %s", err)
	}

	kubeBackupConfig.VaultRole = ""
	if err := authClient(kubeBackupConfig, util.VaultClient); err == nil || err.Error() != "no kubernetes vault role specified" {
		test.Errorf("expected error: no kubernetes vault role specified, actual: %s", err)
	}

	azBackupConfig.VaultRole = ""
	if err := authClient(azBackupConfig, util.VaultClient); err == nil || err.Error() != "no azure vault role specified" {
		test.Errorf("expected error: no azure vault role specified, actual: %s", err)
	}

	approleBackupConfig.VaultRole = ""
	if err := authClient(approleBackupConfig, util.VaultClient); err == nil || err.Error() != "approle credentials absent" {
		test.Errorf("expected error: approle credentials absent, actual: %s", err)
	}
}

// test default mount
func TestCheckAuthParams(test *testing.T) {
	if mount := checkAuthParams("", "", enum.KubernetesSA); mount != "kubernetes" {
		test.Errorf("expected default mount: kubernetes, actual: %s", mount)
	}

	if mount := checkAuthParams("gcp", "", enum.AWSIAM); mount != "gcp" {
		test.Errorf("expected mount input param: gcp, actual: %s", mount)
	}
}

// test vault authenticate with authentication method
func TestLoginWithMethod(test *testing.T) {
	// for now this is encapsulated by TestAuthClient, but in the future it may be useful to here also
}

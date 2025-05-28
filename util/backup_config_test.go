package util

import (
	"os"
	"regexp"

	"testing"

	"github.com/mschuchard/vault-raft-backup/enum"
)

func TestHclDecodeConfig(test *testing.T) {
	config, err := NewBackupConfig("fixtures/valid.hcl")
	if err != nil {
		test.Error("the valid hcl file did not decode properly")
		test.Error(err)
	}
	vaultConfig := *config.VaultConfig
	cloudConfig := *config.CloudConfig
	expectedVaultConfig := VaultConfig{
		Address:      "https://127.0.0.1",
		Insecure:     true,
		Engine:       "token",
		Token:        "foobar",
		AWSMountPath: "azure",
		AWSRole:      "me",
		SnapshotPath: "/path/to/vault.bak",
	}
	expectedCloudConfig := CloudConfig{
		AZAccountURL: "https://foo.com",
		Container:    Container,
		Platform:     enum.AWS,
		Prefix:       Prefix,
	}

	if vaultConfig != expectedVaultConfig || cloudConfig != expectedCloudConfig || !config.SnapshotCleanup {
		test.Error("decoded config struct did not contain expected values")
		test.Errorf("expected vault: %v", expectedVaultConfig)
		test.Errorf("actual vault: %v", vaultConfig)
		test.Errorf("expected cloud: %v", expectedCloudConfig)
		test.Errorf("actual cloud: %v", cloudConfig)
		test.Error("expected snapshot cleanup: true")
		test.Errorf("actual snapshot cleanup: %t", config.SnapshotCleanup)
	}

	_, err = hclDecodeConfig("fixtures/invalid.hcl")
	if err == nil || err.Error() != "fixtures/invalid.hcl:2,3-11: Unsupported argument; An argument named \"does_not\" is not expected here." {
		test.Error("the invalid hcl file did not error, or errored unexpectedly")
		test.Error(err)
	}

	_, err = hclDecodeConfig("fixtures/no_cloud_config.hcl")
	if err == nil || err.Error() != "cloud_config block absent" {
		test.Error("the no_cloud_config hcl file did not error, or errored unexpectedly")
		test.Error(err)
	}
}

func TestOSImportConfig(test *testing.T) {
	// source of truth for values
	const (
		azAccountURL string          = "https://foo.com"
		platform     enum.Platform   = enum.GCP
		addr         string          = "https://127.0.0.1:8234"
		skipVerify   string          = "false"
		authEngine   enum.AuthEngine = enum.VaultToken
		token        string          = "abcdefg"
		awsMount     string          = "gcp"
		awsRole      string          = "my_role"
		snapshotPath string          = "/tmp/my_vault.backup"
	)

	os.Setenv("AZ_ACCOUNT_URL", azAccountURL)
	os.Setenv("CONTAINER", Container)
	os.Setenv("PLATFORM", string(platform))
	os.Setenv("PREFIX", Prefix)
	os.Setenv("VAULT_ADDR", addr)
	os.Setenv("VAULT_SKIP_VERIFY", skipVerify)
	os.Setenv("VAULT_AUTH_ENGINE", string(authEngine))
	os.Setenv("VAULT_TOKEN", token)
	os.Setenv("VAULT_AWS_MOUNT", awsMount)
	os.Setenv("VAULT_AWS_ROLE", awsRole)
	os.Setenv("VAULT_SNAPSHOT_PATH", snapshotPath)
	os.Setenv("SNAPSHOT_CLEANUP", "false")

	config, err := NewBackupConfig("")
	if err != nil {
		test.Error("vault raft backup config failed to construct from environment variables")
		test.Error(err)
	}

	vaultConfig := config.VaultConfig
	cloudConfig := config.CloudConfig
	expectedCloudConfig := CloudConfig{
		AZAccountURL: azAccountURL,
		Container:    Container,
		Platform:     platform,
		Prefix:       Prefix,
	}
	expectedVaultConfig := VaultConfig{
		Address:      addr,
		Insecure:     false,
		Engine:       authEngine,
		Token:        token,
		AWSMountPath: awsMount,
		AWSRole:      awsRole,
		SnapshotPath: snapshotPath,
	}
	if *cloudConfig != expectedCloudConfig || *vaultConfig != expectedVaultConfig || config.SnapshotCleanup {
		test.Error("imported config struct(s) did not initialize with expected values")
		test.Errorf("expected vault: %v", expectedVaultConfig)
		test.Errorf("actual vault: %v", *vaultConfig)
		test.Errorf("expected cloud: %v", expectedCloudConfig)
		test.Errorf("actual cloud: %v", *cloudConfig)
		test.Error("expected snapshot cleanup: false")
		test.Errorf("actual snapshot cleanup: %t", config.SnapshotCleanup)
	}

	// test errors in reverse order
	os.Setenv("PLATFORM", "azure")
	os.Unsetenv("AZ_ACCOUNT_URL")
	if _, err := envImportConfig(); err == nil || err.Error() != "az_account_url environment variable absent" {
		test.Errorf("expected error: az_account_url environment variable absent, actual: %s", err)
	}

	os.Setenv("VAULT_AUTH_ENGINE", "kubernetes")
	if _, err := envImportConfig(); err == nil || err.Error() != "invalid authengine enum" {
		test.Errorf("expected error: invalid authengine enum, actual: %s", err)
	}

	os.Setenv("PLATFORM", "foo")
	if _, err := envImportConfig(); err == nil || err.Error() != "invalid platform enum" {
		test.Errorf("expected error: invalid platform enum, actual: %s", err)
	}

	os.Unsetenv("CONTAINER")
	if _, err = envImportConfig(); err == nil || err.Error() != "environment variable absent" {
		test.Errorf("expected error: environment variable absent, actual: %s", err)
	}

	os.Setenv("SNAPSHOT_CLEANUP", "not a boolean")
	if _, err = envImportConfig(); err == nil || err.Error() != "invalid SNAPSHOT_CLEANUP value" {
		test.Errorf("expected error: invalid SNAPSHOT_CLEANUP value, actual: %s", err)
	}

	os.Setenv("VAULT_SKIP_VERIFY", "not a boolean")
	if _, err = envImportConfig(); err == nil || err.Error() != "invalid VAULT_SKIP_VERIFY value" {
		test.Errorf("expected error: invalid VAULT_SKIP_VERIFY value, actual: %s", err)
	}
}

func TestValidateParameters(test *testing.T) {
	snapshotPath, err := defaultSnapshotPath("")
	if err != nil {
		test.Error("errored with valid input parameters")
		test.Error(err)
	}
	// regexp match for random vault raft snapshot tmp file
	if matched, _ := regexp.MatchString(`/tmp/vault-\d{4}-\d{2}-\d{2}-\d{6}-\d+\.bak`, snapshotPath); !matched {
		test.Error("default snapshot path is not of expected format")
		test.Errorf("expected default snapshot path: /tmp/vault-<datetime>.bak, actual: %s", snapshotPath)
	}
}

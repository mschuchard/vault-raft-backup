package util

import (
	"os"

	"testing"
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
		Container: Container,
		Platform:  AWS,
		Prefix:    Prefix,
	}

	if vaultConfig != expectedVaultConfig || cloudConfig != expectedCloudConfig || !config.SnapshotCleanup {
		test.Error("decoded config struct did not contain expected values")
		test.Errorf("expected vault: %v", expectedVaultConfig)
		test.Errorf("actual vault: %v", vaultConfig)
		test.Errorf("expected aws: %v", expectedCloudConfig)
		test.Errorf("actual aws: %v", cloudConfig)
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
		platform     platform = GCP
		addr         string   = "https://127.0.0.1:8234"
		skipVerify   string   = "false"
		authEngine   string   = "token"
		token        string   = "abcdefg"
		awsMount     string   = "gcp"
		awsRole      string   = "my_role"
		snapshotPath string   = "/tmp/my_vault.backup"
	)

	os.Setenv("CONTAINER", Container)
	os.Setenv("PLATFORM", string(platform))
	os.Setenv("PREFIX", Prefix)
	os.Setenv("VAULT_ADDR", addr)
	os.Setenv("VAULT_SKIP_VERIFY", skipVerify)
	os.Setenv("VAULT_AUTH_ENGINE", authEngine)
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
		Container: Container,
		Platform:  platform,
		Prefix:    Prefix,
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
		test.Errorf("expected aws: %v", expectedCloudConfig)
		test.Errorf("actual aws: %v", *cloudConfig)
		test.Error("expected snapshot cleanup: false")
		test.Errorf("actual snapshot cleanup: %t", config.SnapshotCleanup)
	}

	os.Setenv("PLATFORM", "foo")
	if _, err := envImportConfig(); err == nil || err.Error() != "unsupported platform" {
		test.Errorf("expected error: unsupported platform, actual: %s", err)
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

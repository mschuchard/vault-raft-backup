package util

import (
	"os"

	"testing"
)

func TestHclDecodeConfig(test *testing.T) {
	config, err := hclDecodeConfig("fixtures/valid.hcl")
	if err != nil {
		test.Error("the valid hcl file did not decode properly")
		test.Error(err)
	}
	vaultConfig := *config.VaultConfig
	awsConfig := *config.AWSConfig
	expectedVaultConfig := VaultConfig{
		Address:      "https://127.0.0.1",
		Insecure:     true,
		Engine:       "token",
		Token:        "foobar",
		AWSMountPath: "azure",
		AWSRole:      "me",
		SnapshotPath: "/path/to/vault.bak",
	}
	expectedAWSConfig := AWSConfig{
		S3Bucket: "my_bucket",
		S3Prefix: "prefix",
	}

	if vaultConfig != expectedVaultConfig || awsConfig != expectedAWSConfig || !config.SnapshotCleanup {
		test.Error("decoded config struct did not contain expected values")
		test.Errorf("expected vault: %v", expectedVaultConfig)
		test.Errorf("actual vault: %v", vaultConfig)
		test.Errorf("expected aws: %v", expectedAWSConfig)
		test.Errorf("actual aws: %v", awsConfig)
		test.Error("expected snapshot cleanup: true")
		test.Errorf("actual snapshot cleanup: %t", config.SnapshotCleanup)
	}

	_, err = hclDecodeConfig("fixtures/invalid.hcl")
	if err == nil || err.Error() != "fixtures/invalid.hcl:2,3-11: Unsupported argument; An argument named \"does_not\" is not expected here." {
		test.Error("the invalid hcl file did not error, or errored unexpectedly")
		test.Error(err)
	}
}

func TestOSImportConfig(test *testing.T) {
	os.Setenv("S3_BUCKET", "my_bucket")
	os.Setenv("S3_PREFIX", "my_prefix")
	os.Setenv("VAULT_ADDR", "https://127.0.0.1:8234")
	os.Setenv("VAULT_SKIP_VERIFY", "false")
	os.Setenv("VAULT_AUTH_ENGINE", "token")
	os.Setenv("VAULT_TOKEN", "abcdefg")
	os.Setenv("VAULT_AWS_MOUNT", "gcp")
	os.Setenv("VAULT_AWS_ROLE", "my_role")
	os.Setenv("VAULT_SNAPSHOT_PATH", "/tmp/my_vault.backup")
	os.Setenv("SNAPSHOT_CLEANUP", "true")
	config, err := osImportConfig()
	awsConfig := config.AWSConfig
	vaultConfig := config.VaultConfig
	expectedAWSConfig := AWSConfig{
		S3Bucket: os.Getenv("S3_BUCKET"),
		S3Prefix: os.Getenv("S3_PREFIX"),
	}
	expectedVaultConfig := VaultConfig{
		Address:      os.Getenv("VAULT_ADDR"),
		Insecure:     false,
		Engine:       os.Getenv("VAULT_AUTH_ENGINE"),
		Token:        os.Getenv("VAULT_TOKEN"),
		AWSMountPath: os.Getenv("VAULT_AWS_MOUNT"),
		AWSRole:      os.Getenv("VAULT_AWS_ROLE"),
		SnapshotPath: os.Getenv("VAULT_SNAPSHOT_PATH"),
	}
	if err != nil {
		test.Error("vault raft backup config failed to construct from environment variables")
		test.Error(err)
	}
	if *awsConfig != expectedAWSConfig || *vaultConfig != expectedVaultConfig || !config.SnapshotCleanup {
		test.Error("imported config struct did not initialize with expected values")
		test.Errorf("expected vault: %v", expectedVaultConfig)
		test.Errorf("actual vault: %v", *vaultConfig)
		test.Errorf("expected aws: %v", expectedAWSConfig)
		test.Errorf("actual aws: %v", *awsConfig)
		test.Error("expected snapshot cleanup: true")
		test.Errorf("actual snapshot cleanup: %t", config.SnapshotCleanup)
	}

	os.Setenv("VAULT_SKIP_VERIFY", "not a boolean")
	if _, err = osImportConfig(); err == nil || err.Error() != "invalid VAULT_SKIP_VERIFY value" {
		test.Errorf("expected error: invalid VAULT_SKIP_VERIFY value, actual: %s", err)
	}
	os.Setenv("VAULT_SKIP_VERIFY", "true")

	os.Setenv("SNAPSHOT_CLEANUP", "not a boolean")
	if _, err = osImportConfig(); err == nil || err.Error() != "invalid SNAPSHOT_CLEANUP value" {
		test.Errorf("expected error: invalid SNAPSHOT_CLEANUP value, actual: %s", err)
	}
}

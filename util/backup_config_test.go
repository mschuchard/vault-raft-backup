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
	awsConfig := *config.AWSConfig
	gcpConfig := *config.GCPConfig
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
	expectedGCPConfig := GCPConfig{
		CSBucket: "my_bucket",
		CSPrefix: "prefix",
	}

	if vaultConfig != expectedVaultConfig || awsConfig != expectedAWSConfig || gcpConfig != expectedGCPConfig || !config.SnapshotCleanup {
		test.Error("decoded config struct did not contain expected values")
		test.Errorf("expected vault: %v", expectedVaultConfig)
		test.Errorf("actual vault: %v", vaultConfig)
		test.Errorf("expected aws: %v", expectedAWSConfig)
		test.Errorf("actual aws: %v", awsConfig)
		test.Errorf("expected gcp: %v", expectedGCPConfig)
		test.Errorf("actual gcp: %v", gcpConfig)
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
	// source of truth for values
	const (
		bucket       string = "my_bucket"
		prefix       string = "my_prefix"
		addr         string = "https://127.0.0.1:8234"
		skipVerify   string = "false"
		authEngine   string = "token"
		token        string = "abcdefg"
		awsMount     string = "gcp"
		awsRole      string = "my_role"
		snapshotPath string = "/tmp/my_vault.backup"
	)

	os.Setenv("S3_BUCKET", bucket)
	os.Setenv("S3_PREFIX", prefix)
	os.Setenv("CS_BUCKET", bucket)
	os.Setenv("CS_PREFIX", prefix)
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

	awsConfig := config.AWSConfig
	vaultConfig := config.VaultConfig
	gcpConfig := config.GCPConfig
	expectedAWSConfig := AWSConfig{
		S3Bucket: bucket,
		S3Prefix: prefix,
	}
	expectedGCPConfig := GCPConfig{
		CSBucket: bucket,
		CSPrefix: prefix,
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
	if *awsConfig != expectedAWSConfig || *gcpConfig != expectedGCPConfig || *vaultConfig != expectedVaultConfig || config.SnapshotCleanup {
		test.Error("imported config struct(s) did not initialize with expected values")
		test.Errorf("expected vault: %v", expectedVaultConfig)
		test.Errorf("actual vault: %v", *vaultConfig)
		test.Errorf("expected aws: %v", expectedAWSConfig)
		test.Errorf("actual aws: %v", *awsConfig)
		test.Error("expected snapshot cleanup: false")
		test.Errorf("actual snapshot cleanup: %t", config.SnapshotCleanup)
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

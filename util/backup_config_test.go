package util

import (
	"regexp"

	"testing"

	"github.com/mschuchard/vault-raft-backup/enum"
)

func TestNewBackupConfig(test *testing.T) {
	config, err := NewBackupConfig("fixtures/valid.hcl")
	if err != nil {
		test.Error("the valid hcl file did not decode properly")
		test.Error(err)
	}
	vaultConfig := *config.VaultConfig
	cloudConfig := *config.CloudConfig
	snapshotConfig := *config.SnapshotConfig
	expectedVaultConfig := VaultConfig{
		Address:    "https://127.0.0.1",
		Insecure:   true,
		Engine:     "token",
		Token:      "foobar",
		SecretID:   "abcdef-123456",
		WrapToken:  "abcdef.ghijkl",
		AzResource: "https://management.azure.com/",
		AuthMount:  "azure",
		VaultRole:  "myRole",
		Namespace:  "root",
	}
	expectedCloudConfig := CloudConfig{
		AZAccountURL: "https://foo.com",
		Container:    Container,
		Platform:     enum.AWS,
		Prefix:       Prefix,
	}
	expectedSnapshotConfig := SnapshotConfig{
		Cleanup:          true,
		CompressionLevel: 1,
		Path:             "/path/to/vault.bak",
		Restore:          true,
	}

	if vaultConfig != expectedVaultConfig || cloudConfig != expectedCloudConfig || snapshotConfig != expectedSnapshotConfig {
		test.Error("decoded config struct did not contain expected values")
		test.Errorf("expected vault: %v", expectedVaultConfig)
		test.Errorf("actual vault: %v", vaultConfig)
		test.Errorf("expected cloud: %v", expectedCloudConfig)
		test.Errorf("actual cloud: %v", cloudConfig)
		test.Errorf("expected snapshot: %v", expectedSnapshotConfig)
		test.Errorf("actual snapshot: %v", snapshotConfig)
	}

	_, err = NewBackupConfig("fixtures/invalid.hcl")
	if err == nil || err.Error() != "fixtures/invalid.hcl:1,14-14: Missing required argument; The argument \"auth_engine\" is required, but no definition was found., and 1 other diagnostic(s)" {
		test.Error("the invalid hcl file did not error, or errored unexpectedly")
		test.Error(err)
	}

	_, err = NewBackupConfig("fixtures/no_cloud_config.hcl")
	if err == nil || err.Error() != "cloud_config block absent" {
		test.Error("the no_cloud_config hcl file did not error, or errored unexpectedly")
		test.Error(err)
	}

	_, err = NewBackupConfig("fixtures/no_vault_config.hcl")
	if err == nil || err.Error() != "vault_config block absent" {
		test.Error("the no_vault_config hcl file did not error, or errored unexpectedly")
		test.Error(err)
	}
}

func TestValidateParams(test *testing.T) {
	if err := validateParams(enum.AZ, enum.VaultToken, "https://foo.com", nil); err == nil || err.Error() != "invalid az_account_url value" {
		test.Errorf("expected error: invalid az_account_url value, actual: %s", err)
	}

	if err := validateParams(enum.AZ, enum.VaultToken, "", nil); err == nil || err.Error() != "az_account_url value absent" {
		test.Errorf("expected error: az_account_url value absent, actual: %s", err)
	}

	if err := validateParams(enum.AWS, enum.AuthEngine("foo"), "", nil); err == nil || err.Error() != "invalid authengine enum" {
		test.Errorf("expected error: invalid authengine enum, actual: %s", err)
	}

	if err := validateParams(enum.Platform("foo"), enum.VaultToken, "", nil); err == nil || err.Error() != "invalid platform enum" {
		test.Errorf("expected error: invalid platform enum, actual: %s", err)
	}

	if err := validateParams(enum.GCP, enum.VaultToken, "", &SnapshotConfig{Restore: true, Cleanup: true}); err != nil {
		test.Error("errored during co-specification of snapshot restore and cleanup which is only a warning")
	}

	if err := validateParams(enum.GCP, enum.VaultToken, "", &SnapshotConfig{CompressionLevel: 4}); err == nil || err.Error() != "invalid snapshot compression level" {
		test.Errorf("expected error: invalid snapshot compression level, actual: %s", err)
	}
}

func TestDefaultSnapshotPath(test *testing.T) {
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

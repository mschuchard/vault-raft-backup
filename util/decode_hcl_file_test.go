package util

import "testing"

func TestInvalidHcl(test *testing.T) {
	config, err := HclDecodeConfig("fixtures/valid.hcl")
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

	_, err = HclDecodeConfig("fixtures/invalid.hcl")
	if err == nil || err.Error() != "fixtures/invalid.hcl:2,3-11: Unsupported argument; An argument named \"does_not\" is not expected here." {
		test.Error("the invalid hcl file did not error, or errored unexpectedly")
		test.Error(err)
	}
}

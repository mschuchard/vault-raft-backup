package util

import "testing"

func TestInvalidHcl(test *testing.T) {

	_, err := HclDecodeConfig("fixtures/invalid.hcl")
	if err == nil || err.Error() != "fixtures/invalid.hcl:2,3-11: Unsupported argument; An argument named \"does_not\" is not expected here." {
		test.Error("the invalid hcl file did not error, or errored unexpectedly")
	}
}

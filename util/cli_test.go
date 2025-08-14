package util

import (
	"os"
	"testing"
)

func TestCliConfig(test *testing.T) {
	os.Args[1] = "-c"
	os.Args[2] = "fixtures/valid.hcl"

	if configPath := Cli(); configPath != "fixtures/valid.hcl" {
		test.Error("the config path was not parsed correctly")
		test.Errorf("expected: fixtures/valid.hcl, actual: %s", configPath)
	}
}

// flag redefined issue; needs better workaround
/*func TestCliVersion(test *testing.T) {
	os.Args[1] = "version"
	_ = Cli()
}*/

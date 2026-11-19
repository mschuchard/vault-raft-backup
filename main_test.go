package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mschuchard/vault-raft-backup/testutil"
)

func Example() {
	// write a temporary acceptance test config file
	configContents := fmt.Sprintf(`cloud {
  container = "/tmp"
  platform  = "local"
}

vault {
  auth_engine = "token"
  token       = "%s"
}

snapshot {
  cleanup = true
}
`, testutil.VaultToken)
	if err := os.WriteFile("/tmp/accept.hcl", []byte(configContents), 0o600); err != nil {
		log.Fatalf("failed to write config file for test: %s", err)
	}

	// testutil.Cli reads the config path from os.Args, so replace and later restore them
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()
	os.Args = []string{originalArgs[0], "-c", "/tmp/accept.hcl"}

	main()
	// Output:
}

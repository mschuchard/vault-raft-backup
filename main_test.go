package main

import (
	"os"

	"github.com/mschuchard/vault-raft-backup/util"
)

func ExampleMain() {
	// test that main behaves as expected before snapshot as raft not supported with vault server
	os.Setenv("VAULT_TOKEN", util.VaultToken)
	os.Setenv("S3_BUCKET", "bucket")

	main()
	// Output: foo
}

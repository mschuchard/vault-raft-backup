package main

import (
	"testing"

	"github.com/mschuchard/vault-raft-backup/util"
)

func Example() {
	// test that main behaves as expected before snapshot as raft not supported with vault dev mode server
	test := testing.T{}
	test.Setenv("VAULT_TOKEN", util.VaultToken)
	test.Setenv("CONTAINER", "/tmp")
	test.Setenv("PLATFORM", "local")

	main()
	// Output:
}

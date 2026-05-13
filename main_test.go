package main

import (
	"log"
	"testing"
	"time"

	vault "github.com/hashicorp/vault/api"

	"github.com/mschuchard/vault-raft-backup/util"
)

func Example() {
	// ensure vault server is available (it is possible it has not finished restarting from another test e.g. snapshotrestore)
	client, err := vault.NewClient(&vault.Config{Address: util.VaultAddress})
	if err != nil {
		log.Fatalf("failed to create vault client for validating server: %s", err)
	}
	for i := range 16 {
		if _, err := client.Sys().SealStatus(); err == nil {
			break
		} else if i == 15 {
			log.Fatalf("vault server was not available after fifteen seconds")
		}
		time.Sleep(1 * time.Second)
	}

	// test that main behaves as expected before snapshot as raft not supported with vault dev mode server
	test := testing.T{}
	test.Setenv("VAULT_TOKEN", util.VaultToken)
	test.Setenv("CONTAINER", "/tmp")
	test.Setenv("PLATFORM", "local")

	main()
	// Output:
}

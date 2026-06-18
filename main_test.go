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
	// ensure vault server is healthy after snapshot restoration during unit tests
	for i := range 16 {
		// cluster healthy and unsealed?
		if health, err := client.Sys().Health(); err == nil && !health.Sealed {
			break
		} else if i == 15 {
			// check if error
			if err != nil {
				log.Print(err)
			}
			// for some reason the server never recovered
			log.Fatalf("vault server was not available after fifteen seconds")
		}
		// otherwise wait and try again
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

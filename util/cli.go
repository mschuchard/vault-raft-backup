package util

import (
	"flag"
	"log"
	"os"
)

// rudimentary cli for inputting hcl configuration file and emitting version
func Cli() *string {
	// cli flags for hcl config file path and version
	hclConfigPath := flag.String("c", "", "path to hcl file for backup configuration")
	version := flag.Bool("version", false, "display current version")
	flag.Parse()

	// version output
	if *version {
		log.Print("1.3.0")
		os.Exit(0)
	}

	// verify config file existence
	if len(*hclConfigPath) > 0 {
		if _, err := os.Stat(*hclConfigPath); err != nil {
			log.Fatalf("the config file at %s does not exist", *hclConfigPath)
		}
	}

	// return path to hcl config file
	return hclConfigPath
}

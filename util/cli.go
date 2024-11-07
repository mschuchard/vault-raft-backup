package util

import (
	"flag"
	"log"
	"os"
)

func Cli() *string {
	// cli flags for hcl config file path and version
	hclConfigPath := flag.String("c", "", "path to hcl file for backup configuration")
	version := flag.Bool("version", false, "display current version")
	flag.Parse()

	// version output
	if *version {
		log.Print("1.2.0")
		os.Exit(0)
	}

	// return path to hcl config file
	return hclConfigPath
}

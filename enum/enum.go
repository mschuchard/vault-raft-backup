package enum

import (
	"errors"
	"log"
	"slices"
)

// platform pseudo-enum
type Platform string

const (
	AWS   Platform = "aws"
	AZ    Platform = "azure"
	GCP   Platform = "gcp"
	LOCAL Platform = "local"
)

var platforms []Platform = []Platform{AWS, AZ, GCP, LOCAL}

// platform type conversion
func (p Platform) New() (Platform, error) {
	if !slices.Contains(platforms, p) {
		log.Printf("string %s could not be converted to Platform enum", p)
		return "", errors.New("invalid platform enum")
	}
	return p, nil
}

// authengine pseudo-enum
type AuthEngine string

const (
	AWSIAM     AuthEngine = "aws"
	VaultToken AuthEngine = "token"
	Default    AuthEngine = ""
)

var authEngines []AuthEngine = []AuthEngine{AWSIAM, VaultToken, Default}

// authengine type conversion
func (a AuthEngine) New() (AuthEngine, error) {
	if !slices.Contains(authEngines, a) {
		log.Printf("string %s could not be converted to AuthEngine enum", a)
		return "", errors.New("invalid authengine enum")
	}
	return a, nil
}

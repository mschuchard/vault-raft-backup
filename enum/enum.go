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

// platform type conversion
func (p Platform) New() (Platform, error) {
	if !slices.Contains([]Platform{AWS, AZ, GCP, LOCAL}, p) {
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

// authengine type conversion
func (a AuthEngine) New() (AuthEngine, error) {
	if !slices.Contains([]AuthEngine{AWSIAM, VaultToken, Default}, a) {
		log.Printf("string %s could not be converted to AuthEngine enum", a)
		return "", errors.New("invalid authengine enum")
	}
	return a, nil
}

package enum

import (
	"errors"
	"log"
)

// platform pseudo-enum
type Platform string

const (
	AWS   Platform = "aws"
	GCP   Platform = "gcp"
	LOCAL Platform = "local"
)

// platform type conversion
func (Platform) From(s string) (Platform, error) {
	if s != "aws" && s != "gcp" && s != "local" {
		log.Printf("string %s could not be converted to Platform enum", s)
		return "", errors.New("invalid platform enum")
	}
	return Platform(s), nil
}

// authengine pseudo-enum
type AuthEngine string

const (
	AWSIAM     AuthEngine = "aws"
	VaultToken AuthEngine = "token"
)

// authengine type conversion
func (AuthEngine) From(s string) (AuthEngine, error) {
	if s != "aws" && s != "token" {
		log.Printf("string %s could not be converted to AuthEngine enum", s)
		return "", errors.New("invalid authengine enum")
	}
	return AuthEngine(s), nil
}

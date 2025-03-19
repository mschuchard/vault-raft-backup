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
func (Platform) From(s Platform) (Platform, error) {
	if s != AWS && s != GCP && s != LOCAL {
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
func (AuthEngine) From(s AuthEngine) (AuthEngine, error) {
	if s != AWSIAM && s != VaultToken {
		log.Printf("string %s could not be converted to AuthEngine enum", s)
		return "", errors.New("invalid authengine enum")
	}
	return AuthEngine(s), nil
}

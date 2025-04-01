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
func (p Platform) New() (Platform, error) {
	if p != AWS && p != GCP && p != LOCAL {
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
)

// authengine type conversion
func (a AuthEngine) New() (AuthEngine, error) {
	if a != AWSIAM && a != VaultToken {
		log.Printf("string %s could not be converted to AuthEngine enum", a)
		return "", errors.New("invalid authengine enum")
	}
	return a, nil
}

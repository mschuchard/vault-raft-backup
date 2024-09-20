package storage

import "github.com/mschuchard/vault-raft-backup/util"

func NewConfig(backupConfig *util.AWSConfig) (*awsConfig, error) {
	return NewAWSConfig(backupConfig)
}

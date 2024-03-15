## vault-raft-backup

Vault Raft Backup is a lean tool for creating snapshots of the Raft integrated storage in Hashicorp Vault, and transferring those backups to AWS S3.

This repository and project is based on the work performed for [MITODL](https://github.com/mitodl/vault-raft-backup), and now serves as an upstream for the project hosted within that organization. Although the original work is unlicensed, this repository maintains the BSD-3 license with copyright notice on good faith.

### Prerequsities

If executing as an ad-hoc compile and run (i.e. `go run`), then the dependencies and requirements can be viewed in the [go.mod](go.mod) file. Additional setup requirements are as follows:

- Connectivity to a functioning Vault server cluster with Raft integrated storage.
- Authentication and authorization against the Vault server cluster for executing Raft snapshots.
  - Authentication can be input in general as a token.
  - Authentication can also specified as AWS IAM. In this situation, the Vault server cluster must have a role configured and mapped to an AWS IAM role. This AWS IAM role authorization must also be accessible by the Vault Raft Backup tool somehow (e.g. tool executed on EC2 instance with appropriate IAM Instance profile corresponding to AWS IAM role).
- A local filesystem with permissions and storage capable of staging the snapshot.
- Authentication and authorization against an AWS account for listing, reading, and writing objects to a S3 bucket.
- A S3 bucket capable of storing the snapshot.

The Vault policy for executing Raft snapshots appears like:

```hcl
path "sys/storage/raft/snapshot" {
  capabilities = ["read"]
}
```

### Usage

The following environment variables are read for configuration of the backup tool. This usage is due to the expectation that this tool will be executed as part of automation e.g. pipeline, service, orchestrator, etc. This is also because some inputs are sensitive, and therefore should be constrained to in-process memory.

```
# equivalent to VAULT_ADDR with vault cli executable
# default: http://127.0.0.1:8200
export VAULT_ADDR=<vault server cluster address>
# equivalent to VAULT_SKIP_VERIFY with vault cli executable
# default: false
export VAULT_SKIP_VERIFY=<boolean>
# default: determined based on other inputs
export VAULT_AUTH_ENGINE=<token | aws>
# equivalent to VAULT_TOKEN with vault cli executable
# default: empty
export VAULT_TOKEN=<vault authentication token>
# default: aws
export VAULT_AWS_MOUNT=<vault aws auth engine mount path>
# default: empty
export VAULT_AWS_ROLE=<vault aws authentication role>
# default: <tempdir>/vault.bak
# NOTE: if this file does not exist it will be created with 0600; if it does exist it will be completely overwritten
export VAULT_SNAPSHOT_PATH=<path to local filesystem for snapshot staging>
# required
export S3_BUCKET=<name of s3 bucket for snapshot transfer and storage>
# this is prepended to the base filename in VAULT_SNAPSHOT_PATH
# default: empty
export S3_PREFIX=<snapshot filename prefix during s3 transfer>
```

Additionally, AWS authentication and configuration must be provided with standard methods that do not require manual inputs. The AWS Golang SDK will automatically read authentication information as per normal (i.e. IAM instance profile, `AWS_SHARED_CREDENTIALS_FILE` credentials file, `AWS_PROFILE` config file, environment variables e.g. `AWS_SESSION_TOKEN` and `AWS_REGION`, etc.).

## Contributing
Code should pass all unit and acceptance tests. New features should involve new unit tests.

Please consult the GitHub Project for the current development roadmap.

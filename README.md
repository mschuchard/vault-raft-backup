## vault-raft-backup

Vault Raft Backup is a lean tool for creating snapshots of the Raft integrated storage in [Hashicorp Vault](https://www.vaultproject.io), and transferring those backups to AWS S3, Azure Blob, or GCP Cloud Storage. Additionally the backup can also be maintained on a locally accessible filesystem.

This repository and project is based on the work performed for [MITODL](https://github.com/mitodl/vault-raft-backup), and now serves as an upstream for the project hosted within that organization. Although the original work is unlicensed, this repository maintains the BSD-3 license with copyright notice on good faith.

Statically linked binaries for various operating systems and processor architectures are available at the Github releases page (see link in the right column). Note that the CLI flag `version` will output version information, and then promptly exit.

### Prerequsities
- Connectivity to a functioning Vault server cluster with Raft integrated storage.
- Authentication and authorization against the Vault server cluster for executing Raft snapshots.
  - Authentication can be input in general as a token.
  - Authentication can also specified as AWS IAM. In this situation, the Vault server cluster must have a role configured and mapped to an AWS IAM role. This AWS IAM role authorization must also be accessible by the Vault Raft Backup tool somehow (e.g. tool executed on EC2 instance with appropriate IAM Instance profile corresponding to AWS IAM role).
- A local filesystem with permissions and storage capable of staging the snapshot.
- Authentication and authorization against an AWS, Azure, or GCP account for listing, reading, and writing objects to a S3 bucket, Azure blob, or Cloud Storage bucket. Otherwise, authorization for writing objects to a local filesystem destination.
- A S3 or Cloud Storage bucket, Blob storage, or local filesystem, capable of storing the snapshot.

The Vault policy for authorizing the execution of Raft snapshots appears like:

```hcl
path "sys/storage/raft/snapshot" {
  capabilities = ["read"]
}
```

### Usage

Vault Raft Backup can be configured with either environment variables or a HCL2 config file. The environment variables method is deprecated as of version 1.5.0, and no new input parameters will be supported with that method.

Additionally, AWS, Azure, or GCP authentication and configuration must be provided with standard methods that do not require manual inputs. The AWS Golang SDK will automatically read authentication information as per normal (i.e. IAM instance profile, `AWS_SHARED_CREDENTIALS_FILE` credentials file, `AWS_PROFILE` config file, environment variables e.g. `AWS_SESSION_TOKEN` and `AWS_REGION`, etc.). The GCP and Azure Golang SDKs behave similarly for analogous authentication settings (note that this tool assumes Entra Security Principal authentication for Azure). If the snapshot is stored locally instead, then authentication and authorization will vary greatly based on your personal environment.

#### Environment Variables

The following environment variables are read for the configuration of the backup tool.

```
# VAULT
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

# CLOUD STORAGE
# default: empty (required if PLATFORM is "azure")
export AZ_ACCOUNT_URL=<azure account url>
# required
export CONTAINER=<name of cloud storage destination, or destination directory if storage platform is "local", for final snapshot transfer and storage>
# required
export PLATFORM=<aws, azure, gcp, or local>
# this is prepended to the base filename in VAULT_SNAPSHOT_PATH
# default: empty
export PREFIX=<snapshot filename prefix during storage transfer>

# SNAPSHOT
# determines whether or not the local staged snapshot file is removed after a successful transfer to the final storage location
# default: true
export SNAPSHOT_CLEANUP=<boolean>
# determines the level of compression for the snapshot
# 0: none, 1: fastest, 2: default, 3: most compressed
# default: 0
SNAPSHOT_COMPRESSION_LEVEL=<int>
# default: <tmpdir>/vault-YYYY-MM-DD-hhmmss-<\d+>.bak
# NOTE: if this file does not exist it will be created with 0600; if it does exist it will be completely overwritten
export VAULT_SNAPSHOT_PATH=<path to local filesystem for snapshot staging>
# whether to restore the snapshot instead of backing up (requires local filesystem snapshot storage and additional `write` permissions in the Vault policy)
# default: false
SNAPSHOT_RESTORE = <boolean>
```

#### HCL2 Config File

The HCL2 config file path is passed to the `vault-raft-backup` executable via the `-c` command line argument (e.g. `vault-raft-backup -c config.hcl`). The schema can be viewed below.

```hcl2
vault_config {
  # equivalent to VAULT_ADDR with vault cli executable
  # default: http://127.0.0.1:8200
  address        = <vault server cluster address>
  # equivalent to VAULT_SKIP_VERIFY with vault cli executable
  # default: false
  insecure       = <boolean>
  # default: determined based on other inputs
  auth_engine    = <token | aws>
  # equivalent to VAULT_TOKEN with vault cli executable
  # default: empty
  token          = <vault authentication token>
  # default: aws
  aws_mount_path = <vault aws auth engine mount path>
  # default: empty
  aws_role       = <vault aws authentication role>
}

cloud_config {
  # default: empty (required if platform is "azure")
  az_account_url = <azure account url>
  # required
  container = <name of cloud storage destination, or destination directory if storage platform is "local", for final snapshot transfer and storage>
  # required
  platform  = <aws, azure, gcp, or local>
  # this is prepended to the base filename in VAULT_SNAPSHOT_PATH
  # default: empty
  prefix    = <snapshot filename prefix during storage transfer>
}

snapshot_config {
  # determines whether or not the local staged snapshot file is removed after a successful transfer to the final storage location
  # default: true
  cleanup = <boolean>
  # determines the level of compression for the snapshot
  # 0: none, 1: fastest, 2: default, 3: most compressed
  # default: 0
  compression_level = <int>
  # default: <tmpdir>/vault-YYYY-MM-DD-hhmmss-<\d+>.bak
  # NOTE: if this file does not exist it will be created with 0600; if it does exist it will be completely overwritten
  path = <path on local filesystem for snapshot staging>
  # whether to restore the snapshot instead of backing up (requires local filesystem snapshot storage and additional `write` permissions in the Vault policy)
  # default: false
  restore = <boolean>
}
```

## Contributing
Code should pass all unit and acceptance tests. New features should involve new unit tests.

Please consult the GitHub Project for the current development roadmap.
vault_config {
  address        = "https://127.0.0.1"
  insecure       = true
  auth_engine    = "token"
  token          = "foobar"
  aws_mount_path = "azure"
  aws_role       = "me"
  snapshot_path  = "/path/to/vault.bak"
}

aws_config {
  s3_bucket = "my_bucket"
  s3_prefix = "prefix"
}

snapshot_cleanup = true

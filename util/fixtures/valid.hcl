vault_config {
  address        = "https://127.0.0.1"
  insecure       = true
  auth_engine    = "token"
  token          = "foobar"
  aws_mount_path = "azure"
  aws_role       = "me"
  snapshot_path  = "/path/to/vault.bak"
}

cloud_config {
  az_account_url = "https://foo.com"
  container      = "my_bucket"
  platform       = "aws"
  prefix         = "prefix"
}

snapshot_config {
  cleanup = true
  restore = true
}
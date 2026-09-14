cloud_config {
  az_account_url = "https://foo.com"
  container      = "my_bucket"
  platform       = "aws"
  prefix         = "prefix"
}

vault_config {
  address     = "https://127.0.0.1"
  insecure    = true
  auth_engine = "token"
  token       = "foobar"
  secret_id   = "abcdef-123456"
  wrap_token  = "abcdef.ghijkl"
  az_resource = "https://management.azure.com/"
  auth_mount  = "azure"
  vault_role  = "myRole"
  namespace   = "root"
}

snapshot_config {
  cleanup           = true
  compression_level = 1
  path              = "/path/to/vault.bak"
  restore           = true
}
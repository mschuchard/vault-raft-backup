### 1.5.0 (Next)
- Add snapshot restoration functionality.

### 1.4.2
- Validate GCP client closing.
- Improve authentication token format validation.
- Minor code execution optimization.
- Minor logging improvements.

### 1.4.1
- Improve and fix input parameter value validation.
- Improve snapshot cleanup.
- Fix AWS and Local storage returns.
- Fix error handling for general snapshot transfer.

### 1.4.0
- Add Azure storage support.
- Improve input parameter validation.
- Fix AWS mount path usage in Vault client authentication.

### 1.3.0
- Add local filesystem storage support.
- Improve Vault client configuration.

### 1.2.1
- Add timestamp suffix to default snapshot file name.
- Migrate AWS storage support to SDK v2.
- Remove `-` after `prefix` in storage snapshot name.

### 1.2.0
- Add GCP storage support.

### 1.1.2
- Add version flag
- Add random number string to default snapshot filename suffix for uniqueness.

### 1.1.1
- Fix invalid CLI argument character for HCL file value.
- Fix empty env settings for insecure and cleanup.

### 1.1.0
- Convert local snapshot cleanup from forced to optional.
- Enable configuration via HCL file.

### 1.0.1
- Improve AWS S3 config handling.
- Validate Vault server is unsealed.

### 1.0.0
- Initial release.

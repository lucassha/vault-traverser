# vault-traverser
Walk a given Vault path to check if a secret exists in the given path.

![Testing](https://github.com/lucassha/vault-traverser/actions/workflows/test.yml/badge.svg)

### Testing

A Makefile hosts the necessary components to spin up a test Vault server in Docker and automatically write some secrets into it for testing `traverse`. 

Spin up a test server: `make spinup`. This creates a Vault Dev server from the base Docker image and writes some sample secrets into the default Vault path `/secret`. 


### Usage

![Traverse Usage Example](https://github.com/lucassha/vault-traverser/blob/main/img/github_traverse.gif)

### Examples

```sh
# search the path /secret for the AWS key AKIA-12345678
traverse --path secret --secret AKIA-12345678

# search the path /secret for the AWS key AKIA-12345678. /secret is the default path
traverse --secret AKIA-12345678

# search the path /containers/production for the secret test_key
traverse --path containers/production --secret test_key
```

### TODO
- Add testing for SearchPath method in Vault package
- Add Terraform for S3 bucket
- Deploy to Homebrew
- Fully test kv v1
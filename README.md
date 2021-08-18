# vault-traverser
Walk a Vault path to check if a secret exists in the given path

![Testing](https://github.com/lucassha/vault-traverser/actions/workflows/test.yml/badge.svg)

### Dev

### Testing

A Makefile hosts the necessary components to spin up a test Vault server in Docker and automatically write some secrets into it for testing `traverse`. 

Spin up a test server: `make spinup`. This creates a Vault Dev server from the base Docker image and writes some sample secrets into the default Vault path `/secret`. 



### TODO
- Add testing for SearchPath method in Vault package
- Add Terraform for S3 bucket
- Deploy to Homebrew
- Fully test kv v1
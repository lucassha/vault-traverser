package vault

import (
	"errors"
	"strings"

	"github.com/hashicorp/vault/api"
)

const (
	data = "data"
	keys = "keys"
)

var (
	ErrSecretNotFound = errors.New("secret not found at given path")
)

// VaultClient is used to connect to the Vault API
// and stores a k/v engine version (either v1 or v2)
type VaultClient struct {
	client   *api.Client
	kvEngine string
}

// NewVaultClient instantiates a new client to communicate
// with Vault. It uses the DefaultConfig functionality, which
// sets the VAULT_ADDR == localhost:8200 unless the environment
// variable is set in the user's dev environment
func NewVaultClient(engine string) (*VaultClient, error) {
	config := api.DefaultConfig()

	client, err := api.NewClient(config)
	if err != nil {
		return &VaultClient{}, err
	}

	return &VaultClient{
		client:   client,
		kvEngine: engine,
	}, nil
}

// SearchPath loops over a given Vault path and searches for the secret
// inside this path.
func (v *VaultClient) SearchPath(path, secret string) error {
	return nil
}

// readSecret connects to the Vault cluster and reads a specific endpoint
func (v *VaultClient) readSecret(endpoint string) ([]string, error) {
	if v.kvEngine == "v2" {
		endpoint = mutateReadSecretAPIPath(endpoint)
	}

	secret, err := v.client.Logical().Read(endpoint)
	if err != nil {
		return []string{}, err
	}
	if secret == nil {
		return []string{}, ErrSecretNotFound
	}

	list := []string{}
	for _, v := range secret.Data[data].(map[string]interface{}) {
		list = append(list, v.(string))
	}

	return list, nil
}

// listSecret connects to the Vault cluster and lists all secrets at
// a designated path
func (v *VaultClient) listSecret(path string) ([]string, error) {
	if v.kvEngine == "v2" {
		path = mutateListSecretsAPIPath(path)
	}

	secret, err := v.client.Logical().List(path)
	if err != nil {
		return []string{}, err
	}

	list := []string{}
	for _, v := range secret.Data[keys].([]interface{}) {
		list = append(list, v.(string))
	}

	return list, nil
}

// mutateListSecretsApiPath mutates the api string due to Vault having multiple
// kv versions of their api. For kv v2, "/metadata" must be added to the path.
// "vault kv list secret" is the equivalent to LIST /secret/metadata
// https://www.vaultproject.io/api/secret/kv/kv-v2.html#list-secrets
func mutateListSecretsAPIPath(path string) string {
	fwdSlashIndex := strings.Index(path, "/")
	// forward slash as the first character is not required to read from vault
	// remove the first slash, if present
	if isPrepended(fwdSlashIndex) {
		path = path[1:]
		fwdSlashIndex = strings.Index(path, "/")
	}

	if fwdSlashIndex > -1 {
		return path[0:fwdSlashIndex] + "/metadata" + path[fwdSlashIndex:]
	}
	// no "/" appended to path. thus, the root path was already passed in
	return path + "/metadata"
}

// mutateReadSecretApiPath mutates the api string due to Vault having multiple
// versions of their api. For kv v2, "/data" must be added to the path.
// "vault kv get [path/to/secret]" is the equivalent to GET /path/data/to/secret
func mutateReadSecretAPIPath(endpoint string) string {
	fwdSlashIndex := strings.Index(endpoint, "/")
	// not error checking for an empty value, as cobra in root.go does not allow
	// for an empty secret
	return endpoint[0:fwdSlashIndex] + "/data" + endpoint[fwdSlashIndex:]
}

// isPrepended returns true if the index passed in is 0
func isPrepended(val int) bool {
	return val == 0
}

// fullPath combines two strings to form a full path for the listSecret
// search functionality
func _fullPath(p, q string) string {
	if strings.HasSuffix(p, "/") {
		return p + q
	}
	return p + "/" + q
}

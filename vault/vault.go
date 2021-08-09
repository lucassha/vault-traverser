package vault

import "github.com/hashicorp/vault/api"

// VaultClient is used to connect to the Vault API
// and stores a k/v engine version (either v1 or v2)
type VaultClient struct {
	client *api.Client
	engine string
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
		client: client,
		engine: engine,
	}, nil
}

func (v *VaultClient) ReadSecret(endpoint string) ([]string, error) {
	return []string{}, nil
}

func (v *VaultClient) ListSecret(path string) ([]string, error) {
	return []string{}, nil
}

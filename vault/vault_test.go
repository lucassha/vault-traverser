// Testing is currently only implemented for k/v v2. I haven't been able to get the testing
// working for v1

package vault

import (
	"reflect"
	"testing"
	"time"

	kv "github.com/hashicorp/vault-plugin-secrets-kv"
	"github.com/hashicorp/vault/api"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/logical"
	hashivault "github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/assert"
)

// TestMutateListSecretsAPIPath tests that /metadata path
// has been added into the path for vault kv v2
func TestMutateListSecretsAPIPath(t *testing.T) {
	testTable := []struct {
		name string
		path string
		want string
	}{
		{"forward slash at root", "secret/", "secret/metadata/"},
		{"forward slash in front", "/secret/", "secret/metadata/"},
		{"no forward slash at root", "secret", "secret/metadata"},
		{"multiple nested slashes", "secret/name/shannon", "secret/metadata/name/shannon"},
		{"forward slash with multiple nested slashes", "/secret/name/shannon", "secret/metadata/name/shannon"},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			got := mutateListSecretsAPIPath(tc.path)
			assert.Equal(t, tc.want, got, "these paths should be the same")
		})
	}
}

// TestMutateReadSecretAPIPath tests that /data path
// has been added into the path for vault kv v2
func TestMutateReadSecretAPIPath(t *testing.T) {
	testTable := []struct {
		name     string
		endpoint string
		want     string
	}{
		{"add off root val", "secret/endpoint", "secret/data/endpoint"},
		{"nested val", "secret/endpoint/abc/def/test", "secret/data/endpoint/abc/def/test"},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			got := mutateReadSecretAPIPath(tc.endpoint)
			assert.Equal(t, tc.want, got, "these paths should be the same")
		})
	}
}

// CreateTestVault spins up a Vault server and tests against
// an actual Vault instance. Currently, this is only set up for
// kv v2. Mostly copied from this github issue:
// https://github.com/hashicorp/vault/issues/8440
func createTestVault(t testing.TB) *hashivault.TestCluster {
	t.Helper()

	// CoreConfig parameterizes the Vault core config
	coreConfig := &hashivault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"kv": kv.Factory,
		},
	}

	cluster := hashivault.NewTestCluster(t, coreConfig, &hashivault.TestClusterOptions{
		// Handler returns an http.Handler for the API. This can be used on
		// its own to mount the Vault API within another web server.
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()

	// Create KV V2 mount on the path /test
	// It starts in cluster mode, so you just pick one of the three clients
	// In this case, Cores[0] is just always picking the first one
	if err := cluster.Cores[0].Client.Sys().Mount("test", &api.MountInput{
		Type: "kv",
		Options: map[string]string{
			"version": "2",
		},
	}); err != nil {
		t.Fatal(err)
	}

	return cluster
}

// TestReadandListSecrets creates a new test Vault cluster, writes secrets to paths, and
// then confirms the ReadSecret and ListSecret methods. Both tests are wrapped within this
// single function because tests running in parallel and not sequentially may fail
func TestReadandListSecrets(t *testing.T) {
	cluster := createTestVault(t)
	defer cluster.Cleanup()
	vaultClient := cluster.Cores[0].Client // only need a client from 1 of 3 clusters

	// instead of using NewVaultClient(engine string) constructor,
	// pass in the test client into the VaultClient struct
	vc := &VaultClient{
		client:   vaultClient,
		kvEngine: "v2",
	}

	// set up sample data to write into vault
	testData := []struct {
		path  string
		key   string
		value string
	}{
		{"test/data/test0", "test_0_key", "test_0_data"},
		{"test/data/test1", "test_1_key", "test_1_data"},
		{"test/data/test2", "test_2_key", "test_2_data"},
	}

	// write k/v data pairs into vault
	for _, v := range testData {
		_, err := vc.client.Logical().Write(v.path, map[string]interface{}{
			"data": map[string]interface{}{
				v.key: v.value,
			},
		})
		if err != nil {
			t.Fatal(err)
		}

		// time buffer required between each write
		// https://github.com/hashicorp/terraform-provider-vault/issues/677#issuecomment-609116328
		// Code 400: Errors: Upgrading from non-versioned to versioned data. This backend will be unavailable for a brief period and will resume service shortly.
		time.Sleep(2 * time.Second)
	}

	t.Run("Read Secret Testing", func(t *testing.T) {
		testTable := []struct {
			name       string
			endpoint   string
			key        string
			want       []string
			vaultError error
		}{
			{"find a k/v match", "test/test0", "test_0_key", []string{"test_0_data"}, nil},
			{"do not find a secret", "test/test123", "test_0_key", []string{"test_0_data"}, ErrSecretNotFound},
		}

		for _, tc := range testTable {
			t.Run(tc.name, func(t *testing.T) {
				secrets, err := vc.readSecret(tc.endpoint)
				if err != tc.vaultError {
					t.Fatal(err)
				}

				for i := 0; i < len(secrets); i++ {
					if secrets[i] != tc.want[i] {
						assert.Equal(t, tc.want[i], secrets[i])
					}
				}
			})
		}
	})

	t.Run("List Secrets Testing", func(t *testing.T) {
		testTable := []struct {
			name       string
			path       string
			want       []string
			vaultError error
		}{
			{"find three secrets in path", "test", []string{"test0", "test1", "test2"}, nil},
			{"do not find any secrets in path", "test/123", []string{}, ErrSecretNotFound},
		}

		for _, tc := range testTable {
			t.Run(tc.name, func(t *testing.T) {
				secrets, err := vc.listSecret(tc.path)
				if err != tc.vaultError {
					t.Fatal(err)
				}

				if !reflect.DeepEqual(secrets, tc.want) {
					t.Errorf("arrays not equal. expected %v but got %v", tc.want, secrets)
				}
			})
		}
	})
}

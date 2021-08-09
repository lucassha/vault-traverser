package vault

import (
	"testing"

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

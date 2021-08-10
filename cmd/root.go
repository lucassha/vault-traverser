package cmd

import (
	"errors"
	"fmt"

	"github.com/lucassha/vault-traverser/vault"
	"github.com/spf13/cobra"
)

type flagOptions struct {
	path   string
	secret string
	engine string
}

// var traverseFlag flagOptions
var (
	traverseFlag flagOptions
	ErrNoPath    = errors.New("no path provided")
	ErrNoSecret  = errors.New("no secret provided")
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "traverse",
	Short: "Traverse a Vault path to determine if a secret exists in the given path",
	Long: `Traverse allows you to search an entire Vault path(s) to search for hidden keys.

By default, traverse searches the 'secret' path.

Example:

# search the secret/ path for the secret "AKIA-123ASLDFD"
traverse --secret AKIA-123ASLDFD

# search the containers/teams path for the secret "AKIA-123ASLDFD"
traverse --path containers/teams --secret AKIA-123ASLDFD
`,
	RunE: traverse,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	// cobra.CheckErr(rootCmd.Execute())
	rootCmd.Execute()
}

func traverse(cmd *cobra.Command, args []string) error {
	if traverseFlag.path == "" {
		return ErrNoPath
	}
	if traverseFlag.secret == "" {
		return ErrNoSecret
	}

	vc, err := vault.NewVaultClient(traverseFlag.engine)
	if err != nil {
		return err
	}

	err = vc.SearchPath(traverseFlag.path, traverseFlag.secret)
	if err == vault.ErrSecretNotFound {
		fmt.Printf("secret [%s] not found\n", traverseFlag.secret)
		return nil
	}

	return err
}

func init() {
	rootCmd.AddCommand(versionCmd)

	rootCmd.PersistentFlags().StringVarP(&traverseFlag.path, "path", "p", "secret", "Vault path to search for a secret")
	rootCmd.PersistentFlags().StringVarP(&traverseFlag.engine, "engine", "e", "v2", "K/V secrets engine. Use KV v1 for < 0.10 Vault")
	rootCmd.PersistentFlags().StringVarP(&traverseFlag.secret, "secret", "s", "", "Secret key to search for")

	rootCmd.MarkFlagRequired("secret")
}

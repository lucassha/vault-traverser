package cmd

import (
	"github.com/spf13/cobra"
)

type flagOptions struct {
	path   string
	secret string
	engine string
}

var traverseFlag flagOptions

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
	cobra.CheckErr(rootCmd.Execute())
}

func traverse(cmd *cobra.Command, args []string) error {
	return nil
}

func init() {
	rootCmd.AddCommand(versionCmd)

	rootCmd.PersistentFlags().StringVarP(&traverseFlag.path, "path", "p", "/secret", "Vault path to search for a secret")
	rootCmd.PersistentFlags().StringVarP(&traverseFlag.path, "engine", "e", "v2", "K/V secrets engine. Use KV v1 for < 0.10 Vault")
	rootCmd.PersistentFlags().StringVarP(&traverseFlag.path, "secret", "s", "", "Secret key to search for")

	rootCmd.MarkFlagRequired("secret")
}

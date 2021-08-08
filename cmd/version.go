package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Output current versin of traverse",
	Run:   printVersion,
}

var version = "0.0.1"

func printVersion(cmd *cobra.Command, args []string) {
	fmt.Fprintf(os.Stdout, "%s\n", version)
}

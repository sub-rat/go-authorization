package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"go-authorization/cmd/apiserver"
	"go-authorization/cmd/migrate"
	"go-authorization/cmd/setup"
	"os"
)

func init() {
	rootCmd.AddCommand(apiserver.StartCmd)
	rootCmd.AddCommand(migrate.StartCmd)
	rootCmd.AddCommand(setup.StartCmd)
}

var rootCmd = &cobra.Command{
	Use:          "rbac-boilerplate",
	Short:        "rbac-boilerplate is casbin based rbac integrated boilerplate for kickstart your project ",
	SilenceUsage: true,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New(
				"requires at least one arg, " +
					"you can view the available parameters through `--help`")
		}
		return nil
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

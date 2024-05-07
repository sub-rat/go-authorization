package generate

import (
	"github.com/spf13/cobra"
	"go-authorization/generator"
	"go-authorization/lib"
)

var configFile string
var name string

func init() {
	pf := StartCmd.PersistentFlags()
	pf.StringVarP(&configFile, "config", "c",
		"config/config.yaml", "this parameter is used to start the service application")
	pf.StringVarP(&name, "name", "n",
		"", "this parameter is used to generate the controller, service and repository")
	cobra.MarkFlagRequired(pf, "config")
	cobra.MarkFlagRequired(pf, "name")
}

var StartCmd = &cobra.Command{
	Use:          "gen",
	Short:        "Set up data for the application",
	Example:      "{execfile} init -c config/settings.yml -n example",
	SilenceUsage: true,
	PreRun: func(cmd *cobra.Command, args []string) {
		lib.SetConfigPath(configFile)
	},
	Run: func(cmd *cobra.Command, args []string) {
		config := lib.NewConfig()
		logger := lib.NewLogger(config)
		generator.CreateFile(logger, name)
	},
}

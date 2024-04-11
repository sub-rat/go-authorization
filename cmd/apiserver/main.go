package apiserver

import (
	"github.com/spf13/cobra"
	"go-authorization/bootstrap"
	"go-authorization/lib"
)

var configFile string
var casbinModel string
var StartCmd = &cobra.Command{
	Use:          "api-server",
	Short:        "Start API Server",
	Example:      "{execfile} server -c config/config.yaml",
	SilenceUsage: true,
	PreRun: func(cmd *cobra.Command, args []string) {
		lib.SetConfigPath(configFile)
		lib.SetConfigCasbinModelPath(casbinModel)
	},
	Run: func(cmd *cobra.Command, args []string) {
		runApplication()
	},
}

func init() {
	pf := StartCmd.PersistentFlags()
	pf.StringVarP(&configFile,
		"config", "c",
		"config/config.yaml",
		"this parameter is used to start the service application")

	pf.StringVarP(&casbinModel,
		"casbin_model", "m",
		"config/casbin_model.conf", "this parameter is used for the running configuration of casbin")

	cobra.MarkFlagRequired(pf, "config")
	cobra.MarkFlagRequired(pf, "casbin_model")
}

func runApplication() {
	bootstrap.New()
}

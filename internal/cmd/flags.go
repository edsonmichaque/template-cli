package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func addGlobalFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().Bool(optSandbox, false, "Sandbox environment")
	cmd.PersistentFlags().String(configAccessToken, "", "Access token")
	cmd.PersistentFlags().String(configAccount, "", "Account")
	cmd.PersistentFlags().String(configBaseURL, "", "Base URL")
	cmd.PersistentFlags().StringVar(&profile, optProfile, defaultProfile, "Profile")
	cmd.PersistentFlags().StringVarP(&configFile, configConfigFile, "c", "", "Configuration file")

	cmd.MarkFlagsMutuallyExclusive(configBaseURL, optSandbox)

	viper.SetEnvPrefix(envPrefix)
	if err := viper.BindPFlags(cmd.PersistentFlags()); err != nil {
		panic(err)
	}
}

func applyOutputFlag(cmd *cobra.Command, f string) {
	cmd.Flags().StringP(optOutput, "o", f, "Output format")
}

func applyQueryFlag(cmd *cobra.Command) {
	cmd.Flags().StringP(optQuery, "q", "", "Query")
}

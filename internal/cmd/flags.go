package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func withFlagsGlobal() cmdOption {
	return func(cmd *cobra.Command) {
		cmd.PersistentFlags().Bool(optSandbox, false, "Sandbox environment")
		cmd.PersistentFlags().String(optAccessToken, "", "Access token")
		cmd.PersistentFlags().String(optAccount, "", "Account")
		cmd.PersistentFlags().String(optBaseURL, "", "Base URL")
		cmd.PersistentFlags().StringVar(&profile, optProfile, defaultProfile, "Profile")
		cmd.PersistentFlags().StringVarP(&configFile, optConfigFile, "c", "", "Configuration file")
		cmd.MarkFlagsMutuallyExclusive(optBaseURL, optSandbox)

		viper.SetEnvPrefix(envPrefix)
	}
}

// withFlagOutput adds output flag to command
func withFlagOutput(value string) cmdOption {
	return func(cmd *cobra.Command) {
		cmd.Flags().StringP(optOutput, "o", value, "Output format")
	}
}

// withFlagQuery adds query flag to command
func withFlagQuery() cmdOption {
	return func(cmd *cobra.Command) {
		cmd.Flags().StringP(optQuery, "q", "", "Query")
	}
}

func withOpts(opts *Opts) cmdOption {
	return func(cmd *cobra.Command) {
		cmd.SetOutput(opts.Stdout)
		cmd.SetIn(opts.Stdin)
		cmd.SetErr(opts.Stderr)
	}
}

// withCmd adds subcommand to command
func withCmd(children ...*Cmd) cmdOption {
	return func(cmd *cobra.Command) {
		for _, c := range children {
			cmd.AddCommand(c.Command)
		}
	}
}

// withFlagDomain adds domain flag to command
func withFlagDomain(value string, required bool) cmdOption {
	return func(cmd *cobra.Command) {
		cmd.Flags().StringP(optDomain, "d", value, "Domain")

		if required {
			cmd.MarkFlagRequired(optDomain)
		}
	}
}

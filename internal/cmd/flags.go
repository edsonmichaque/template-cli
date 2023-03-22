package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func withGlobalFlags() commandOption {
	return func(cmd *cobra.Command) {
		cmd.PersistentFlags().Bool(optSandbox, false, "Sandbox environment")
		cmd.PersistentFlags().String(optAccessToken, "", "Access token")
		cmd.PersistentFlags().String(optAccount, "", "Account")
		cmd.PersistentFlags().String(optBaseURL, "", "Base URL")
		cmd.PersistentFlags().StringVar(&profile, optProfile, defaultProfile, "Profile")
		cmd.PersistentFlags().StringVarP(&configFile, optConfigFile, "c", "", "Configuration file")

		cmd.MarkFlagsMutuallyExclusive(optBaseURL, optSandbox)

		viper.SetEnvPrefix(envPrefix)
		if err := viper.BindPFlags(cmd.PersistentFlags()); err != nil {
			panic(err)
		}
	}
}

// withOutputFlag adds output flag to command
func withOutputFlag(value string) commandOption {
	return func(cmd *cobra.Command) {
		cmd.Flags().StringP(optOutput, "o", value, "Output format")
	}
}

// withQueryFlag adds query flag to command
func withQueryFlag() commandOption {
	return func(cmd *cobra.Command) {
		cmd.Flags().StringP(optQuery, "q", "", "Query")
	}
}

func withOptions(opts *Options) commandOption {
	return func(cmd *cobra.Command) {
		cmd.SetOutput(opts.Stdout)
		cmd.SetIn(opts.Stdin)
		cmd.SetErr(opts.Stderr)
	}
}

// withSubcommand adds subcommand to command
func withSubcommand(children ...*Command) commandOption {
	return func(cmd *cobra.Command) {
		for _, c := range children {
			cmd.AddCommand(c.Command)
		}
	}
}

// withDomainFlag adds domain flag to command
func withDomainFlag(value string) commandOption {
	return func(cmd *cobra.Command) {
		cmd.Flags().StringP(optDomain, "d", value, "Domain")
	}
}

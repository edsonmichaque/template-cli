// Copyright 2023 Edson Michaque
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/edsonmichaque/template-cli/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	cfgFormatJSON = "json"
	cfgFormatYAML = "yaml"
	cfgFormatTOML = "toml"
	cfgFormatYML  = "yml"
)

var (
	prodBaseURL    = "https://api.dnsimple.com"
	sandboxBaseURL = "https://api.sandbox.dnsimple.com"

	configProps = map[string]struct{}{
		optAccount:     {},
		optBaseURL:     {},
		optAccessToken: {},
		optSandbox:     {},
	}

	cfgValidate = map[string]func(string) (interface{}, error){
		optSandbox: func(value string) (interface{}, error) {
			return strconv.ParseBool(value)
		},
		optAccount: func(value string) (interface{}, error) {
			return strconv.ParseInt(value, 10, 64)
		},
		optBaseURL: func(value string) (interface{}, error) {
			return value, nil
		},
		optAccessToken: func(value string) (interface{}, error) {
			return value, nil
		},
	}
)

func cmdConfig(opts *Opts) *Cmd {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage configurations",
	}

	return initCmd(
		cmd,
		withOpts(opts),
		withCmd(
			cmdConfigInit(opts),
			cmdConfigGet(opts),
			cmdConfigSet(opts),
		),
	)
}

func cmdConfigInit(opts *Opts) *Cmd {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Manage configurations",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.LoadWithValidation(false)
			if err != nil {
				return wrapError(exitFailure, err)
			}

			home, err := os.UserConfigDir()
			if err != nil {
				return wrapError(exitFailure, err)
			}

			profile := viper.GetString(optProfile)

			cmd.Println(fmt.Sprintf("Configuring profile '%s'", profile))

			cfg, ext, err := promptConfig(cfg)
			if err != nil {
				return wrapError(exitFailure, err)
			}

			cfgViper := viper.New()

			cfgViper.Set(optAccount, cfg.Account)
			cfgViper.Set(optAccessToken, cfg.AccessToken)

			if cfg.BaseURL != "" {
				cfgViper.Set(optBaseURL, cfg.BaseURL)
			}

			if cfg.Sandbox {
				cfgViper.Set(optSandbox, cfg.Sandbox)
			}

			cfgPath := filepath.Join(
				home,
				bin,
				fmt.Sprintf("%s.%s", profile, strings.ToLower(ext)),
			)
			if err := cfgViper.WriteConfigAs(cfgPath); err != nil {
				return wrapError(exitFailure, err)
			}

			return nil
		},
	}

	return initCmd(
		cmd,
		withOpts(opts),
	)
}

func cmdConfigGet(opts *Opts) *Cmd {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Manage configurations",
		Args: func(cmd *cobra.Command, args []string) error {
			return preRunE(
				func() error {
					return cobra.ExactArgs(1)(cmd, args)
				},
				func() error {
					if _, ok := configProps[args[0]]; !ok {
						return errors.New("not found")
					}

					return nil
				},
			)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Println(viper.GetString(args[0]))

			return nil
		},
	}

	return initCmd(
		cmd,
		withOpts(opts),
	)
}

func cmdConfigSet(opts *Opts) *Cmd {
	cmd := &cobra.Command{
		Use:   "set",
		Short: "Manage configurations",
		Args: func(cmd *cobra.Command, args []string) error {
			return preRunE(
				func() error {
					return cobra.ExactArgs(2)(cmd, args)
				},
				func() error {
					if _, ok := configProps[args[0]]; !ok {
						return errors.New("not found")
					}

					return nil
				},
			)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			validate := cfgValidate[args[0]]
			if validate == nil {
				return newError(exitFailure, "no validator found")
			}

			value, err := validate(args[1])
			if err != nil {
				return wrapError(exitFailure, err)
			}

			viper.Set(args[0], value)

			if err := viper.WriteConfig(); err != nil {
				return wrapError(exitFailure, err)
			}

			return nil
		},
	}

	return initCmd(
		cmd,
		withOpts(opts),
	)
}

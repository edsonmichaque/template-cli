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
	configFormatJSON = "json"
	configFormatYAML = "yaml"
	configFormatTOML = "toml"
	configFormatYML  = "yml"
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

	validateConfig = map[string]func(string) (interface{}, error){
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

func cmdConfig(opts *Options) *Cmd {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage configurations",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.NewWithValidation(false)
			if err != nil {
				return newError(1, err)
			}

			home, err := os.UserConfigDir()
			if err != nil {
				return newError(1, err)
			}

			profile := viper.GetString(optProfile)

			cmd.Println(fmt.Sprintf("Configuring profile '%s'", profile))
			cfg, ext, err := promptConfig(cfg)
			if err != nil {
				return newError(1, err)
			}

			v := viper.New()
			v.Set(optAccount, cfg.Account)
			v.Set(optAccessToken, cfg.AccessToken)
			if cfg.BaseURL != "" {
				v.Set(optBaseURL, cfg.BaseURL)
			}

			if cfg.Sandbox {
				v.Set(optSandbox, cfg.Sandbox)
			}

			cfgPath := filepath.Join(home, pathTemplate, fmt.Sprintf("%s.%s", profile, strings.ToLower(ext)))
			if err := v.WriteConfigAs(cfgPath); err != nil {
				return newError(1, err)
			}

			return nil
		},
	}

	return newCmd(
		cmd,
		withOptions(opts),
		withSubcommand(
			cmdConfigGet(opts),
			cmdConfigSet(opts),
		),
	)
}

func cmdConfigGet(opts *Options) *Cmd {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Manage configurations",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, ok := configProps[args[0]]; !ok {
				return newError(1, errors.New("not found"))
			}

			cmd.Println(viper.GetString(args[0]))

			return nil
		},
	}

	return newCmd(cmd, withOptions(opts))
}

func cmdConfigSet(opts *Options) *Cmd {
	cmd := &cobra.Command{
		Use:   "set",
		Short: "Manage configurations",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, ok := configProps[args[0]]; !ok {
				return newError(1, errors.New("not found"))
			}

			validate := validateConfig[args[0]]
			if validate == nil {
				return newError(1, errors.New("no validator found"))
			}

			value, err := validate(args[1])
			if err != nil {
				return newError(1, err)
			}

			viper.Set(args[0], value)

			if err := viper.WriteConfig(); err != nil {
				return newError(1, err)
			}

			return nil
		},
	}

	return newCmd(cmd, withOptions(opts))
}

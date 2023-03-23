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
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configFile string
	profile    string
)

const (
	binaryName              = "template"
	defaultConfigFileFormat = "yaml"
	defaultProfile          = "default"
	envDev                  = "DEV"
	envPrefix               = "TEMPLATE"
	envProd                 = "PROD"
	envSandbox              = "SANDBOX"
	envTemplateConfigFile   = "TEMPLATE_CONFIG_FILE"
	envTemplateProfile      = "TEMPLATE_PROFILE"
	envXDGConfigHome        = "XDG_CONFIG_HOME"
	formatJSON              = "json"
	formatTable             = "table"
	formatText              = "text"
	formatYAML              = "yaml"
	optAccessToken          = "access-token"
	optAccount              = "account"
	optBaseURL              = "base-url"
	optCollaboratorID       = "collaborator-id"
	optConfigFile           = "config-file"
	optConfirm              = "confirm"
	optDomain               = "domain"
	optOutput               = "output"
	optPage                 = "page"
	optPerPage              = "per-page"
	optProfile              = "profile"
	optQuery                = "query"
	optRecordID             = "record-id"
	optSandbox              = "sandbox"
	optionFromFile          = "from-file"
	pathConfigFile          = "/etc/template"
	pathTemplate            = "template"
)

func RunWithOptions(opts *Options) error {
	return cmdRoot(opts).Execute()
}

func Run() error {
	opts, err := NewOptions()
	if err != nil {
		return err
	}

	return RunWithOptions(opts)
}

func cmdRoot(opts *Options) *Cmd {
	cmd := &cobra.Command{
		Use: binaryName,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return viper.BindPFlags(cmd.PersistentFlags())
		},
		SilenceUsage: true,
	}

	cobra.OnInitialize(initConfig)

	return newCmd(
		cmd,
		withSubcommand(cmdFoo(opts)),
		withSubcommand(cmdBar(opts)),
		withSubcommand(cmdConfig(opts)),
		withSubcommand(cmdVersion(opts)),
		withGlobalFlags(),
	)
}

func initConfig() {
	var err error
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else if configFile := os.Getenv(envTemplateConfigFile); configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		configHome := os.Getenv(envXDGConfigHome)
		if configHome == "" {
			configHome, err = os.UserConfigDir()
			cobra.CheckErr(err)
		}

		viper.AddConfigPath(filepath.Join(configHome, pathTemplate))
		viper.AddConfigPath(pathConfigFile)
		viper.SetConfigType(defaultConfigFileFormat)

		viper.SetConfigName(defaultProfile)

		if profileEnv := os.Getenv(envTemplateProfile); profileEnv != "" {
			viper.SetConfigName(profileEnv)
		}

		if profile != "" {
			viper.SetConfigName(profile)
		}
	}

	bindEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Println("Found error: ", err.Error())
		}
	}
}

type Cmd struct {
	*cobra.Command
}

type cmdOption func(*cobra.Command)

func newCmd(c *cobra.Command, opts ...cmdOption) *Cmd {
	for _, opt := range opts {
		opt(c)
	}

	return &Cmd{
		Command: c,
	}
}

func bindEnv() {
	for _, v := range os.Environ() {
		parts := strings.Split(v, "=")
		if len(parts) == 0 {
			continue
		}

		if !strings.HasPrefix(parts[0], envPrefix+"_") {
			continue
		}

		viper.BindEnv(envToFlag(v), parts[0])
	}
}

func flagToEnv(env string) string {
	env = strings.ReplaceAll(env, "-", "_")
	env = strings.ToUpper(env)

	return fmt.Sprintf("%s_%s", envPrefix, env)
}

func envToFlag(env string) string {
	env = strings.TrimPrefix(env, envPrefix+"_")
	parts := strings.Split(env, "=")
	env = strings.ToLower(parts[0])
	env = strings.ReplaceAll(env, "_", "-")

	return env
}

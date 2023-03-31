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
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configFile string
	profile    string
)

const (
	binaryName        = "template"
	defaultProfile    = "default"
	envDev            = "DEV"
	envPrefix         = "TEMPLATE"
	envProd           = "PROD"
	envSandbox        = "SANDBOX"
	envCfgFile        = "TEMPLATE_CONFIG_FILE"
	envProfile        = "TEMPLATE_PROFILE"
	envCfgHome        = "XDG_CONFIG_HOME"
	formatJSON        = "json"
	formatTable       = "table"
	formatText        = "text"
	formatYAML        = "yaml"
	optAccessToken    = "access-token"
	optAccount        = "account"
	optBaseURL        = "base-url"
	optCollaboratorID = "collaborator-id"
	optConfigFile     = "config-file"
	optConfirm        = "confirm"
	optDomain         = "domain"
	optOutput         = "output"
	optPage           = "page"
	optPerPage        = "per-page"
	optProfile        = "profile"
	optQuery          = "query"
	optRecordID       = "record-id"
	optSandbox        = "sandbox"
	optionFromFile    = "from-file"
	pathConfigFile    = "/etc/template"
)

func init() {
	cobra.OnInitialize(initConfig)
	viperBindEnv()
}

func Run() error {
	opts, err := InitOpts()
	if err != nil {
		return err
	}

	return runWithOpts(opts)
}

func runWithOpts(opts *Opts) error {
	return cmdRoot(opts).Execute()
}

func cmdRoot(opts *Opts) *Cmd {

	cmd := &cobra.Command{
		Use: binaryName,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return viper.BindPFlags(cmd.PersistentFlags())
		},
		SilenceUsage: true,
	}

	return initCmd(
		cmd,
		withCmd(cmdFoo(opts)),
		withCmd(cmdBar(opts)),
		withCmd(cmdConfig(opts)),
		withCmd(cmdVersion(opts)),
		withFlagsGlobal(),
	)
}

func cfgCfg(env string) (string, error) {
	var (
		dir string
		err error
	)

	if dir = os.Getenv(envCfgHome); dir != "" {
		dir, err = os.UserConfigDir()
		if err != nil {
			return "", err
		}
	}

	cfgDir = filepath.Join(dir, binaryName)

	if env := os.Getenv(envProfile); env != "" {
		cfgName = env
	}
}

func initConfig() {
	var (
		cfgFile string
		cfgName string
		cfgDir  string
	)

	var err error
	if configFile != "" {
		cfgFile = configFile
	}

	if path := os.Getenv(envCfgFile); path != "" && configFile == "" {
		cfgFile = path
	}

	cfgName = defaultProfile

	if dir := os.Getenv(envCfgHome); dir != "" {
		dir, err = os.UserConfigDir()
		cobra.CheckErr(err)

		cfgDir = dir
	}

	if os.Getenv(envCfgHome) != "" {
		dir := os.Getenv(envCfgHome)
		if dir == "" {
			dir, err = os.UserConfigDir()
			cobra.CheckErr(err)
		}

		cfgDir = filepath.Join(dir, binaryName)

		if env := os.Getenv(envProfile); env != "" {
			cfgName = env
		}
	}

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	if cfgDir != "" && cfgName != "" {
		viper.AddConfigPath(cfgDir)
		viper.SetConfigName(cfgName)
	}

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

func initCmd(cmd *cobra.Command, applyFn ...cmdOption) *Cmd {
	for _, apply := range applyFn {
		apply(cmd)
	}

	return &Cmd{
		Command: cmd,
	}
}

func viperBindEnv() {
	for _, v := range os.Environ() {
		parts := strings.Split(v, "=")
		if len(parts) != 2 {
			continue
		}

		if !strings.HasPrefix(parts[0], envPrefix+"_") {
			continue
		}

		env, err := envToFlag(v)
		if err != nil {
			continue
		}

		viper.BindEnv(env, parts[0])
	}
}

func flagToEnv(env string) string {
	env = strings.ReplaceAll(env, "-", "_")
	env = strings.ToUpper(env)

	return fmt.Sprintf("%s_%s", envPrefix, env)
}

func envToFlag(env string) (string, error) {
	env = strings.TrimPrefix(env, envPrefix+"_")
	parts := strings.Split(env, "=")

	if len(parts) != 2 {
		return "", errors.New("Invalid env var")
	}

	env = strings.ToLower(parts[0])
	env = strings.ReplaceAll(env, "_", "-")

	return env, nil
}

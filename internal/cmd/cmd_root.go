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
	"io"
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
	bin               = "template"
	defaultProfile    = "default"
	envConfigFile     = "TEMPLATE_CONFIG_FILE"
	envConfigHome     = "XDG_CONFIG_HOME"
	envDev            = "DEV"
	envPrefix         = "TEMPLATE"
	envProd           = "PROD"
	envProfile        = "TEMPLATE_PROFILE"
	envSandbox        = "SANDBOX"
	optAccessToken    = "access-token"
	optAccount        = "account"
	optBaseURL        = "base-url"
	optCollaboratorID = "collaborator-id"
	optConfigFile     = "config-file"
	optConfirm        = "confirm"
	optDomain         = "domain"
	optFormat         = "format"
	optFromFile       = "from-file"
	optOutput         = "output"
	optPage           = "page"
	optPerPage        = "per-page"
	optProfile        = "profile"
	optQuery          = "query"
	optRecordID       = "record-id"
	optSandbox        = "sandbox"
	outputJSON        = "json"
	outputTable       = "table"
	outputText        = "text"
	outputYAML        = "yaml"
	pathConfigFile    = "/etc/template"
)

func init() {
	cobra.OnInitialize(initCfg)
	viperBindFlags()
}

func Run() error {
	return run()
}

func run() error {
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
		Use: bin,
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

func initCfg() {
	var (
		cfgFile string
		cfgName string
		cfgDir  string
	)

	var err error
	if configFile != "" {
		cfgFile = configFile
	}

	if path := os.Getenv(envConfigFile); path != "" && configFile == "" {
		cfgFile = path
	}

	cfgName = defaultProfile

	if dir := os.Getenv(envConfigHome); dir != "" {
		dir, err = os.UserConfigDir()
		cobra.CheckErr(err)

		cfgDir = dir
	} else {
		if os.Getenv(envConfigHome) != "" {
			dir := os.Getenv(envConfigHome)
			if dir == "" {
				dir, err = os.UserConfigDir()
				cobra.CheckErr(err)
			}

			cfgDir = filepath.Join(dir, bin)

			if env := os.Getenv(envProfile); env != "" {
				cfgName = env
			}
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

func initCmd(cmd *cobra.Command, opts ...cmdOption) *Cmd {
	for _, opt := range opts {
		opt(cmd)
	}

	return &Cmd{
		Command: cmd,
	}
}

func viperBindFlags() {
	for _, v := range os.Environ() {
		envParts := strings.Split(v, "=")
		if len(envParts) != 2 {
			continue
		}

		if !strings.HasPrefix(envParts[0], envPrefix+"_") {
			continue
		}

		env, err := envToFlag(v)
		if err != nil {
			continue
		}

		_ = viper.BindEnv(env, envParts[0])
	}
}

func flagToEnv(env string) string {
	env = strings.ReplaceAll(env, "-", "_")
	env = strings.ToUpper(env)

	return fmt.Sprintf("%s_%s", envPrefix, env)
}

func envToFlag(env string) (string, error) {
	env = strings.TrimPrefix(env, envPrefix+"_")
	envParts := strings.Split(env, "=")

	if len(envParts) != 2 {
		return "", errors.New("Invalid env var")
	}

	env = strings.ToLower(envParts[0])
	env = strings.ReplaceAll(env, "_", "-")

	return env, nil
}

func print(cmd *cobra.Command, r io.Reader) error {
	if _, err := io.Copy(cmd.OutOrStdout(), r); err != nil {
		return err
	}

	return nil
}

func flagContains(flag string, values []string) error {
	flagValue := viper.GetString(flag)

	for _, value := range values {
		if flagValue == value {
			return nil
		}
	}

	return fmt.Errorf(`flag "%s" has invalid value "%s"`, flag, flagValue)
}

func preRunE(fn ...func() error) error {
	for _, preRun := range fn {
		if err := preRun(); err != nil {
			return err
		}
	}

	return nil
}

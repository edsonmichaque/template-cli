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
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/edsonmichaque/template-cli/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func cmdBar(opts *Opts) *Cmd {
	cmd := &cobra.Command{
		Use:   "bar",
		Short: "List accounts",
		Example: heredoc.Doc(`
			template bar
			template bar --output=json
			template bar --output=yaml
			template bar --output=json --query="[].id"
		`),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return wrapError(exitFailure, err)
			}

			cmd.Printf("%#v", cfg)

			return nil
		},
	}

	return initCmd(
		cmd,
		withFlagOutput(formatTable),
		withFlagQuery(),
		withOpts(opts),
	)
}

const (
	exitFailure = 1
)

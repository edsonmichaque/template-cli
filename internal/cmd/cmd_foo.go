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
	"io"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/edsonmichaque/template-cli/internal/cmd/formatter"
	"github.com/edsonmichaque/template-cli/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func cmdWrite(cmd *cobra.Command, r io.Reader) error {
	if _, err := io.Copy(cmd.OutOrStdout(), r); err != nil {
		return err
	}

	return nil
}

func flagContainsValue(flag string, values []string) error {
	flagValue := viper.GetString(flag)

	for _, value := range values {
		if flagValue == value {
			return nil
		}
	}

	return fmt.Errorf(`flag "%s" has invalid value "%s"`, flag, flagValue)

}

func chainPreRunFunctions(fn ...func() error) error {
	for _, preRun := range fn {
		if err := preRun(); err != nil {
			return err
		}
	}

	return nil
}

func cmdFoo(opts *Opts) *Cmd {
	cmd := &cobra.Command{
		Use:   "foo",
		Short: "List accounts",
		Example: heredoc.Doc(`
			template foo
			template foo --output=json
			template foo --output=yaml
			template foo --output=json --query="[].id"
		`),
		Args: cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return chainPreRunFunctions(
				func() error {
					return viper.BindPFlags(cmd.Flags())
				},
				func() error {
					return flagContainsValue(
						optOutput,
						[]string{
							outputJSON,
							outputYAML,
							outputTable,
						},
					)
				},
			)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := config.Load()
			if err != nil {
				return wrapError(1, err)
			}

			list := formatter.FooList{
				formatter.Foo{
					ID:   1,
					Name: "First Name",
					Age:  "19",
				},
				formatter.Foo{
					ID:   2,
					Name: "First Name",
					Age:  "19",
				},
			}

			f, err := formatter.Format(list, &formatter.Opts{
				Output: formatter.Output(viper.GetString(optOutput)),
			})

			if err != nil {
				return wrapError(1, err)
			}

			if err := cmdWrite(cmd, f); err != nil {
				return wrapError(1, err)
			}

			return nil
		},
	}

	return initCmd(
		cmd,
		withFlagOutput(outputTable),
		withFlagQuery(),
		withOpts(opts),
	)
}

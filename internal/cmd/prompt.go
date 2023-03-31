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

	"github.com/AlecAivazis/survey/v2"
	"github.com/edsonmichaque/template-cli/internal/config"
)

func runConfigPrompt(c *config.Config) (*config.Config, string, error) {
	res, err := execPrompt(
		promptAccount(c.Account),
		promptAccessToken(c.AccessToken),
		promptBaseURL("https://example.con"),
		promptFileFmt(cfgFmtJSON),
		promptConfirm("Do you want to save?", true),
	)

	if err != nil {
		return nil, "", err
	}

	var cfg config.Config

	if value := res.GetString(optAccount); value != "" {
		cfg.Account = value
	}

	if value := res.GetString(optAccessToken); value != "" {
		cfg.AccessToken = value
	}

	if value := res.GetString(optBaseURL); value != "" {
		cfg.BaseURL = value
	}

	var fileFormat string
	if value := res.GetString(optFormat); value != "" {
		fileFormat = value
	}

	fmt.Printf("%#v\n", cfg)

	return &cfg, fileFormat, nil
}

func promptAccessToken(value string) runPromptFunc {
	return runPromptFunc(func() (*promptResult, error) {
		var token string

		if err := survey.AskOne(
			&survey.Input{
				Message: "Access token",
				Default: value,
			},
			&token,
		); err != nil {
			return nil, err
		}

		return &promptResult{
			Name:  optAccessToken,
			Value: token,
		}, nil
	})
}

func promptAccount(value string) runPromptFunc {
	return runPromptFunc(func() (*promptResult, error) {
		var account string

		if err := survey.AskOne(
			&survey.Input{
				Message: "Account",
				Default: value,
			},
			&account,
		); err != nil {
			return nil, err
		}

		return &promptResult{
			Name:  optAccount,
			Value: account,
		}, nil
	})
}

func promptEnvironment(value string) runPromptFunc {
	return runPromptFunc(func() (*promptResult, error) {
		var env string

		if err := survey.AskOne(
			&survey.Select{
				Message: "Environment",
				Options: []string{
					envProd,
					envSandbox,
					envDev,
				},
				Default: value,
			},
			&env,
		); err != nil {
			return nil, err
		}

		return &promptResult{
			Name:  optAccessToken,
			Value: env,
		}, nil
	})
}

func promptBaseURL(value string) runPromptFunc {
	return runPromptFunc(func() (*promptResult, error) {
		var url string

		if err := survey.AskOne(
			&survey.Input{
				Message: "Base URL",
				Default: value,
			},
			&url,
		); err != nil {
			return nil, err
		}

		return &promptResult{
			Name:  optBaseURL,
			Value: url,
		}, nil
	})
}

func promptFileFmt(value string) runPromptFunc {
	return runPromptFunc(func() (*promptResult, error) {
		var format string

		if err := survey.AskOne(
			&survey.Select{
				Message: "File format",
				Options: []string{
					cfgFmtJSON,
					cfgFmtYAML,
					cfgFmtTOML,
				},
				Default: value,
			},
			&format,
		); err != nil {
			return nil, err
		}

		return &promptResult{
			Name:  optFormat,
			Value: format,
		}, nil
	})
}

func promptConfirm(msg string, value bool) runPromptFunc {
	return runPromptFunc(func() (*promptResult, error) {
		var confirmation bool

		if err := survey.AskOne(
			&survey.Confirm{
				Message: msg,
				Default: false,
			},
			&confirmation,
		); err != nil {
			return nil, err
		}

		return &promptResult{
			Name:  "confirmation",
			Value: confirmation,
		}, nil
	})
}

type promptResult struct {
	Name  string
	Value interface{}
}

type promptRunner interface {
	runPrompt() (*promptResult, error)
}

type runPromptFunc func() (*promptResult, error)

func (r runPromptFunc) runPrompt() (*promptResult, error) {
	return r()
}

type runPromptResult map[string]interface{}

func (r runPromptResult) Get(key string) interface{} {
	if s, ok := r[key]; ok {
		return s
	}

	return nil
}

func (r runPromptResult) GetString(key string) string {
	value := r.Get(key)
	if value == nil {
		return ""
	}

	return value.(string)
}

func execPrompt(runners ...promptRunner) (runPromptResult, error) {
	result := make(map[string]interface{})

	for _, runner := range runners {
		res, err := runner.runPrompt()
		if err != nil {
			return nil, err
		}

		result[res.Name] = res.Value
	}

	return result, nil
}

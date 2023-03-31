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
	res, err := runPrompt(
		promptAccountID(c.Account),
		promptAccessToken(c.AccessToken),
		promptBaseURL("https://example.con"),
		promptFileFormat(cfgFormatJSON),
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

		prompt := &survey.Input{
			Message: "Access Token",
			Default: value,
		}

		var token string
		if err := survey.AskOne(prompt, &token); err != nil {
			return nil, err
		}

		return &promptResult{
			Name:  optAccessToken,
			Value: token,
		}, nil
	})
}

func promptAccountID(value string) runPromptFunc {
	return runPromptFunc(func() (*promptResult, error) {
		prompt := &survey.Input{
			Message: "Account ID",
			Default: value,
		}

		var accountID string
		if err := survey.AskOne(prompt, &accountID); err != nil {
			return nil, err
		}

		return &promptResult{
			Name:  optAccount,
			Value: accountID,
		}, nil
	})
}

func promptEnvironment(value string) runPromptFunc {
	return runPromptFunc(func() (*promptResult, error) {

		prompt := &survey.Select{
			Message: "Environment",
			Options: []string{envProd, envSandbox, envDev},
			Default: value,
		}

		var env string
		if err := survey.AskOne(prompt, &env); err != nil {
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
		prompt := &survey.Input{
			Message: "Base URL",
			Default: value,
		}

		var baseURL string
		if err := survey.AskOne(prompt, &baseURL); err != nil {
			return nil, err
		}

		return &promptResult{
			Name:  optBaseURL,
			Value: baseURL,
		}, nil
	})
}

func promptFileFormat(value string) runPromptFunc {
	return runPromptFunc(func() (*promptResult, error) {
		prompt := &survey.Select{
			Message: "File format",
			Options: []string{cfgFormatJSON, cfgFormatYAML, cfgFormatTOML},
			Default: value,
		}

		var fileFormat string
		if err := survey.AskOne(prompt, &fileFormat); err != nil {
			return nil, err
		}

		return &promptResult{
			Name:  optFormat,
			Value: fileFormat,
		}, nil
	})
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

func runPrompt(items ...promptRunner) (runPromptResult, error) {
	output := make(map[string]interface{})

	for _, item := range items {
		kv, err := item.RunPrompt()
		if err != nil {
			return nil, err
		}

		output[kv.Name] = kv.Value
	}

	return output, nil
}

type promptRunner interface {
	RunPrompt() (*promptResult, error)
}

type runPromptFunc func() (*promptResult, error)

func (p runPromptFunc) RunPrompt() (*promptResult, error) {
	return p()
}

func promptConfirm(msg string, value bool) runPromptFunc {
	return runPromptFunc(func() (*promptResult, error) {
		prompt := &survey.Confirm{
			Message: msg,
			Default: false,
		}

		var confirmation bool

		if err := survey.AskOne(prompt, &confirmation); err != nil {
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

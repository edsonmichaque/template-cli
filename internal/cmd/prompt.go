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

	"github.com/AlecAivazis/survey/v2"
	"github.com/edsonmichaque/template-cli/internal/config"
)

func runConfirm(domain string) bool {
	confirm := false
	prompt := &survey.Confirm{
		Message: fmt.Sprintf("Do you want to delete domain %s?", domain),
	}

	if err := survey.AskOne(prompt, &confirm); err != nil {
		return false
	}

	return confirm
}

func promptConfig(c *config.Config) (*config.Config, string, error) {
	/*baseURL := prodBaseURL

	accountID, err := promptAccountID(c.Account)
	if err != nil {
		return nil, "", err
	}

	accessToken, err := promptAccessToken(c.AccessToken)
	if err != nil {
		return nil, "", err
	}

	env, err := promptEnvironment(envProd)
	if err != nil {
		return nil, "", err
	}

	if env == envSandbox {
		baseURL = sandboxBaseURL
	}

	if env == envDev {
		baseURL, err = promptBaseURL(prodBaseURL)
		if err != nil {
			return nil, "", err
		}
	}

	fileFormat, err := promptFileFormat(configFormatJSON)
	if err != nil {
		return nil, "", err
	}

	confirmation, err := promptConfirmation("Do you want to save?", true)
	if err != nil {
		return nil, "", err
	}

	if !confirmation {
		return nil, "", errors.New("did not confirm")
	}

	cfg := config.Config{
		Account:     accountID,
		AccessToken: accessToken,
	}

	if env == envDev {
		cfg.BaseURL = baseURL
	}

	if env == envSandbox {
		cfg.Sandbox = true
	}

	return &cfg, fileFormat, nil*/
	return nil, "", errors.New("")
}

func promptAccessToken(value string) (string, error) {
	prompt := &survey.Input{
		Message: "Access Token",
		Default: value,
	}

	var token string
	if err := survey.AskOne(prompt, &token); err != nil {
		return "", err
	}

	return token, nil
}

func promptAccountID(value string) (string, error) {
	prompt := &survey.Input{
		Message: "Account ID",
		Default: value,
	}

	var accountID string
	if err := survey.AskOne(prompt, &accountID); err != nil {
		return "", err
	}

	return accountID, nil
}

func promptEnvironment(value string) (string, error) {
	prompt := &survey.Select{
		Message: "Environment",
		Options: []string{envProd, envSandbox, envDev},
		Default: value,
	}

	var env string
	if err := survey.AskOne(prompt, &env); err != nil {
		return "", err
	}

	return env, nil
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
			Name:  "format",
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
			Name:  "format",
			Value: fileFormat,
		}, nil
	})
}

func promptConfirmation(msg string, value bool) (bool, error) {
	var confirmation bool

	prompt := &survey.Confirm{
		Message: msg,
		Default: value,
	}

	if err := survey.AskOne(prompt, &confirmation); err != nil {
		return false, err
	}

	return confirmation, nil
}

type runPromptResult map[string]interface{}

func (r runPromptResult) Get(key string) interface{} {
	if s, ok := r[key]; ok {
		return s
	}

	return nil
}

func (r runPromptResult) GetSting(key string) string {
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

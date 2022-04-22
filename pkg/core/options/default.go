/*
 * Copyright 2022 The Furiko Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package options

import (
	"fmt"
	"strings"

	execution "github.com/furiko-io/furiko/apis/execution/v1alpha1"
)

// MakeDefaultOptions returns a substitution map containing default values for all options
// in the given OptionSpec.
func MakeDefaultOptions(cfg *execution.OptionSpec) (map[string]string, error) {
	opts := make(map[string]string)
	if cfg == nil {
		return opts, nil
	}

	for _, option := range cfg.Options {
		// Evaluate default option variables.
		val, err := EvaluateOptionDefault(option)
		if err != nil {
			return nil, err
		}
		opts[MakeOptionVariableName(option)] = val
	}

	return opts, nil
}

// EvaluateOptionDefault evaluates the default value of the option.
func EvaluateOptionDefault(option execution.Option) (string, error) {
	switch option.Type {
	case execution.OptionTypeBool:
		return EvaluateOptionDefaultBool(option.Bool)
	case execution.OptionTypeString:
		return EvaluateOptionDefaultString(option.String)
	case execution.OptionTypeSelect:
		return EvaluateOptionDefaultSelect(option.Select)
	case execution.OptionTypeMulti:
		return EvaluateOptionDefaultMulti(option.Multi)
	case execution.OptionTypeDate:
		return "", nil
	}
	return "", fmt.Errorf(`unsupported option type "%v"`, option.Type)
}

// EvaluateOptionDefaultBool evaluates the default value of the Bool option.
func EvaluateOptionDefaultBool(cfg *execution.BoolOptionConfig) (string, error) {
	if cfg == nil {
		cfg = &execution.BoolOptionConfig{}
	}
	return cfg.FormatValue(GetOptionDefaultBoolValue(cfg))
}

// EvaluateOptionDefaultString evaluates the default value of the String option.
func EvaluateOptionDefaultString(cfg *execution.StringOptionConfig) (string, error) {
	if cfg == nil {
		cfg = &execution.StringOptionConfig{}
	}
	value := GetOptionDefaultStringValue(cfg)
	if cfg.TrimSpaces {
		value = strings.TrimSpace(value)
	}
	return value, nil
}

// EvaluateOptionDefaultSelect evaluates the default value of the Select option.
func EvaluateOptionDefaultSelect(cfg *execution.SelectOptionConfig) (string, error) {
	if cfg == nil {
		cfg = &execution.SelectOptionConfig{}
	}
	return GetOptionDefaultSelectValue(cfg), nil
}

// EvaluateOptionDefaultMulti evaluates the default value of the Multi option.
func EvaluateOptionDefaultMulti(cfg *execution.MultiOptionConfig) (string, error) {
	if cfg == nil {
		cfg = &execution.MultiOptionConfig{}
	}
	return strings.Join(GetOptionDefaultMultiValue(cfg), cfg.Delimiter), nil
}

// GetOptionDefaultValue returns the default option value prior to evaluation.
func GetOptionDefaultValue(option execution.Option) (interface{}, error) {
	switch option.Type {
	case execution.OptionTypeBool:
		return GetOptionDefaultBoolValue(option.Bool), nil
	case execution.OptionTypeString:
		return GetOptionDefaultStringValue(option.String), nil
	case execution.OptionTypeSelect:
		return GetOptionDefaultSelectValue(option.Select), nil
	case execution.OptionTypeMulti:
		return GetOptionDefaultMultiValue(option.Multi), nil
	case execution.OptionTypeDate:
		return GetOptionDefaultDateValue(option.Date), nil
	}
	return "", fmt.Errorf(`unsupported option type "%v"`, option.Type)
}

// GetOptionDefaultBoolValue returns the default value of the Bool option.
func GetOptionDefaultBoolValue(cfg *execution.BoolOptionConfig) bool {
	if cfg == nil {
		cfg = &execution.BoolOptionConfig{}
	}
	return cfg.Default
}

// GetOptionDefaultStringValue returns the default value of the String option.
func GetOptionDefaultStringValue(cfg *execution.StringOptionConfig) string {
	if cfg == nil {
		cfg = &execution.StringOptionConfig{}
	}
	value := cfg.Default
	if cfg.TrimSpaces {
		value = strings.TrimSpace(value)
	}
	return value
}

// GetOptionDefaultSelectValue returns the default value of the Select option.
func GetOptionDefaultSelectValue(cfg *execution.SelectOptionConfig) string {
	if cfg == nil {
		cfg = &execution.SelectOptionConfig{}
	}
	return cfg.Default
}

// GetOptionDefaultMultiValue returns the default value of the Multi option.
func GetOptionDefaultMultiValue(cfg *execution.MultiOptionConfig) []string {
	if cfg == nil {
		cfg = &execution.MultiOptionConfig{}
	}
	return cfg.Default
}

// GetOptionDefaultDateValue returns the default value of the Date option.
func GetOptionDefaultDateValue(_ *execution.DateOptionConfig) string {
	return ""
}

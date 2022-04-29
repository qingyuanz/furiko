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

package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/furiko-io/furiko/pkg/cli/printer"
)

func NewGetCommand(streams *Streams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get one or more resources by name.",
	}

	// Add common flags
	cmd.PersistentFlags().StringP("output", "o", string(printer.OutputFormatPretty),
		fmt.Sprintf("Output format. One of: %v", strings.Join(printer.GetAllOutputFormatStrings(), "|")))

	cmd.AddCommand(NewGetJobCommand(streams))
	cmd.AddCommand(NewGetJobConfigCommand(streams))

	return cmd
}

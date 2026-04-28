// Copyright 2026 DataRobot, Inc. and its affiliates.
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

package get

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/datarobot/cli/internal/auth"
	"github.com/datarobot/cli/internal/pipelines"
	"github.com/datarobot/cli/tui"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	var outputFormat string

	cmd := &cobra.Command{
		Use:   "get <pipeline-id>",
		Short: "Display details of a pipeline.",
		Long: `Display full details of a pipeline including all versions.

By default, output is human-readable. Use --output json for machine-parseable output.

Example:
  dr pipelines get 8a8d6e5e-1234-5678-90ab-cdef01234567
  dr pipelines get 8a8d6e5e-1234-5678-90ab-cdef01234567 --output json`,
		Args:    cobra.ExactArgs(1),
		PreRunE: auth.EnsureAuthenticatedE,
		RunE: func(_ *cobra.Command, args []string) error {
			if outputFormat != "" && outputFormat != "json" {
				return fmt.Errorf("invalid output format: %s (supported: json)", outputFormat)
			}

			pipeline, err := pipelines.GetPipeline(args[0])
			if err != nil {
				return err
			}

			if outputFormat == "json" {
				return printGetJSON(*pipeline)
			}

			printGetHuman(*pipeline)

			return nil
		},
	}

	cmd.Flags().StringVar(&outputFormat, "output", "", "Output format (json)")

	return cmd
}

func printGetJSON(pipeline pipelines.Pipeline) error {
	data, err := json.MarshalIndent(pipeline, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(data))

	return nil
}

func printGetHuman(pipeline pipelines.Pipeline) {
	description := "\u2014"
	if pipeline.Description != "" {
		description = pipeline.Description
	}

	fmt.Println(tui.BaseTextStyle.Render("ID:          " + pipeline.PipelineID))
	fmt.Println(tui.BaseTextStyle.Render("Name:        " + pipeline.Name))
	fmt.Println(tui.BaseTextStyle.Render("Description: " + description))
	fmt.Println(tui.BaseTextStyle.Render("Mode:        " + pipeline.Mode))
	fmt.Println(tui.BaseTextStyle.Render(fmt.Sprintf("Active:      %t", pipeline.IsActive)))
	fmt.Println(tui.DimStyle.Render("Created:     " + pipeline.CreatedAt.UTC().Format(time.RFC3339)))
	fmt.Println(tui.DimStyle.Render("Updated:     " + pipeline.UpdatedAt.UTC().Format(time.RFC3339)))

	if len(pipeline.Versions) == 0 {
		return
	}

	fmt.Println()
	fmt.Println(tui.BaseTextStyle.Render(fmt.Sprintf("Versions (%d):", len(pipeline.Versions))))

	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	fmt.Fprintln(writer, "  VERSION\tSTATUS\tPYTHON\tCREATED\tELECTRONS")

	for _, version := range pipeline.Versions {
		electrons := "\u2014"
		if len(version.ElectronNames) > 0 {
			electrons = strings.Join(version.ElectronNames, ", ")
		}

		python := version.PythonVersion
		if python == "" {
			python = "\u2014"
		}

		fmt.Fprintf(writer, "  v%s\t%s\t%s\t%s\t%s\n",
			strconv.Itoa(version.Version),
			version.Status,
			python,
			version.CreatedAt.UTC().Format(time.RFC3339),
			electrons,
		)
	}

	_ = writer.Flush()

	for _, version := range pipeline.Versions {
		if version.ErrorDetail == "" {
			continue
		}

		fmt.Println(tui.DimStyle.Render(fmt.Sprintf("  v%d error: %s", version.Version, version.ErrorDetail)))
	}
}

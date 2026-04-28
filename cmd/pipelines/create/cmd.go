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

package create

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/datarobot/cli/internal/auth"
	"github.com/datarobot/cli/internal/pipelines"
	"github.com/datarobot/cli/tui"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	var (
		description  string
		mode         string
		outputFormat string
	)

	cmd := &cobra.Command{
		Use:   "create <file>",
		Short: "Upload a Python file to create a pipeline.",
		Long: `Upload a Python file containing a Covalent lattice to register a new pipeline.

The lattice name is extracted from the file and used as the pipeline name.
By default, output is human-readable. Use --output json for machine-parseable output.

Example:
  dr pipelines create ./my_pipeline.py
  dr pipelines create ./my_pipeline.py --description "First draft" --mode draft
  dr pipelines create ./my_pipeline.py --output json`,
		Args:    cobra.ExactArgs(1),
		PreRunE: auth.EnsureAuthenticatedE,
		RunE: func(_ *cobra.Command, args []string) error {
			if outputFormat != "" && outputFormat != "json" {
				return fmt.Errorf("invalid output format: %s (supported: json)", outputFormat)
			}

			if mode != "" && mode != pipelines.ModeDraft && mode != pipelines.ModeLocked {
				return fmt.Errorf("invalid mode: %s (supported: draft, locked)", mode)
			}

			result, err := pipelines.CreatePipeline(args[0], description, mode)
			if err != nil {
				return err
			}

			if outputFormat == "json" {
				return printCreateJSON(*result)
			}

			printCreateHuman(*result)

			return nil
		},
	}

	cmd.Flags().StringVar(&description, "description", "", "Optional description for the pipeline")
	cmd.Flags().StringVar(&mode, "mode", "", "Pipeline mode: draft (default) or locked")
	cmd.Flags().StringVar(&outputFormat, "output", "", "Output format (json)")

	return cmd
}

func printCreateJSON(result pipelines.CreateResponse) error {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(data))

	return nil
}

func printCreateHuman(result pipelines.CreateResponse) {
	electrons := "\u2014"
	if len(result.ElectronNames) > 0 {
		electrons = strings.Join(result.ElectronNames, ", ")
	}

	fmt.Println(tui.BaseTextStyle.Render("Pipeline:  " + result.PipelineID))
	fmt.Println(tui.BaseTextStyle.Render("Name:      " + result.Name))
	fmt.Println(tui.BaseTextStyle.Render("Version:   " + strconv.Itoa(result.Version)))
	fmt.Println(tui.BaseTextStyle.Render("Status:    " + result.Status))
	fmt.Println(tui.BaseTextStyle.Render("Mode:      " + result.Mode))
	fmt.Println(tui.BaseTextStyle.Render("Electrons: " + electrons))
	fmt.Println(tui.DimStyle.Render("Created:   " + result.CreatedAt.UTC().Format(time.RFC3339)))
}

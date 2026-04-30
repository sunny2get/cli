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
	"errors"
	"fmt"

	"github.com/datarobot/cli/cmd/pipelines/dispatch/dispatchutil"
	"github.com/datarobot/cli/cmd/pipelines/scopeflag"
	"github.com/datarobot/cli/internal/auth"
	"github.com/datarobot/cli/internal/pipelines"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	var (
		flags        scopeflag.Flags
		inputID      string
		outputFormat string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Trigger a pipeline dispatch",
		Long: `Trigger a new dispatch (single execution) of a pipeline.

The dispatch is created in PENDING state. Use ` + "`dr pipelines dispatch get`" + `
or ` + "`dr pipelines dispatch status`" + ` to follow its progress.

Example:
  dr pipelines dispatch create --pipeline <id> --input <input-id>
  dr pipelines dispatch create --pipeline <id> --version=2 --input <input-id> --output json`,
		Args:         cobra.NoArgs,
		PreRunE:      auth.EnsureAuthenticatedE,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if outputFormat != "" && outputFormat != "json" {
				return fmt.Errorf("invalid output format: %s (supported: json)", outputFormat)
			}

			if flags.PipelineID == "" {
				return errors.New("--pipeline is required")
			}

			if inputID == "" {
				return errors.New("--input is required")
			}

			scope, version, err := flags.Resolve(cmd)
			if err != nil {
				return err
			}

			result, err := pipelines.CreateDispatch(flags.PipelineID, scope, version, inputID)
			if err != nil {
				return err
			}

			if outputFormat == "json" {
				return dispatchutil.PrintDispatchJSON(*result)
			}

			dispatchutil.PrintDispatchHuman(*result)

			return nil
		},
	}

	flags.Bind(cmd)
	cmd.Flags().StringVar(&inputID, "input", "", "Input ID to dispatch")
	cmd.Flags().StringVar(&outputFormat, "output", "", "Output format (json)")

	return cmd
}

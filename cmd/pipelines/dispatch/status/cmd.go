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

package status

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/datarobot/cli/cmd/pipelines/dispatch/dispatchutil"
	"github.com/datarobot/cli/cmd/pipelines/scopeflag"
	"github.com/datarobot/cli/internal/auth"
	"github.com/datarobot/cli/internal/drapi"
	"github.com/datarobot/cli/internal/pipelines"
	"github.com/datarobot/cli/tui"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	var (
		flags        scopeflag.Flags
		outputFormat string
	)

	cmd := &cobra.Command{
		Use:   "status <dispatch-id>",
		Short: "Get the lightweight status of a pipeline dispatch",
		Long: `Poll a dispatch's current status without re-downloading the full record.

Example:
  dr pipelines dispatch status --pipeline <id> <dispatch-id>
  dr pipelines dispatch status --pipeline <id> --version=2 <dispatch-id> --output json`,
		Args:         cobra.ExactArgs(1),
		PreRunE:      auth.EnsureAuthenticatedE,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if outputFormat != "" && outputFormat != "json" {
				return fmt.Errorf("invalid output format: %s (supported: json)", outputFormat)
			}

			if flags.PipelineID == "" {
				return errors.New("--pipeline is required")
			}

			scope, version, err := flags.Resolve(cmd)
			if err != nil {
				return err
			}

			result, err := pipelines.GetDispatchStatus(flags.PipelineID, scope, version, args[0])
			if err != nil {
				return handleStatusError(err, args[0])
			}

			if outputFormat == "json" {
				return dispatchutil.PrintStatusJSON(*result)
			}

			dispatchutil.PrintStatusHuman(*result)

			return nil
		},
	}

	flags.Bind(cmd)
	cmd.Flags().StringVar(&outputFormat, "output", "", "Output format (json)")

	return cmd
}

func handleStatusError(err error, dispatchID string) error {
	var httpErr *drapi.HTTPError

	if errors.As(err, &httpErr) && httpErr.StatusCode == http.StatusNotFound {
		fmt.Println(tui.DimStyle.Render("No dispatch found with id: " + dispatchID))

		return nil
	}

	return err
}

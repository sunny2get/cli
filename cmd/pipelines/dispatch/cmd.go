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

package dispatch

import (
	"github.com/datarobot/cli/cmd/pipelines/dispatch/cancel"
	"github.com/datarobot/cli/cmd/pipelines/dispatch/create"
	"github.com/datarobot/cli/cmd/pipelines/dispatch/get"
	"github.com/datarobot/cli/cmd/pipelines/dispatch/list"
	"github.com/datarobot/cli/cmd/pipelines/dispatch/status"
	"github.com/spf13/cobra"
)

// Cmd returns the parent command for `dr pipelines dispatch`.
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dispatch",
		Short: "Manage pipeline dispatches",
		Long: `Trigger and inspect dispatches (single executions) of a pipeline.

Dispatches come in two scopes:
  - draft   : runs against the in-flight draft of a pipeline
  - locked  : runs against a specific frozen version

When --version is supplied, the locked scope is selected automatically.`,
	}

	cmd.AddCommand(
		create.Cmd(),
		list.Cmd(),
		get.Cmd(),
		status.Cmd(),
		cancel.Cmd(),
	)

	return cmd
}

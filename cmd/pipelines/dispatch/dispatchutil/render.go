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

// Package dispatchutil holds the rendering helpers shared by the
// `dr pipelines dispatch` verbs. Living in a sibling package keeps the
// parent `dispatch` package free of cycles.

package dispatchutil

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/datarobot/cli/internal/pipelines"
	"github.com/datarobot/cli/tui"
)

// PrintDispatchJSON marshals a dispatch as indented JSON.
func PrintDispatchJSON(d pipelines.Dispatch) error {
	data, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(data))

	return nil
}

// PrintDispatchHuman renders a single dispatch in a human-friendly form.
func PrintDispatchHuman(d pipelines.Dispatch) {
	scope := "draft"
	versionDisplay := "\u2014"

	if d.VersionID != nil {
		scope = "locked"
		versionDisplay = "v" + strconv.Itoa(*d.VersionID)
	}

	covalent := d.CovalentDispatchID
	if covalent == "" {
		covalent = "\u2014"
	}

	fmt.Println(tui.BaseTextStyle.Render("Dispatch ID:        " + d.DispatchID))
	fmt.Println(tui.BaseTextStyle.Render("Pipeline ID:        " + d.PipelineID))
	fmt.Println(tui.BaseTextStyle.Render("Scope:              " + scope))
	fmt.Println(tui.BaseTextStyle.Render("Version:            " + versionDisplay))
	fmt.Println(tui.BaseTextStyle.Render("Input ID:           " + d.InputID))
	fmt.Println(tui.BaseTextStyle.Render("Status:             " + d.Status))
	fmt.Println(tui.BaseTextStyle.Render("Triggered By:       " + d.TriggeredBy))
	fmt.Println(tui.BaseTextStyle.Render("Covalent Dispatch:  " + covalent))

	if d.ErrorDetail != "" {
		fmt.Println(tui.BaseTextStyle.Render("Error:              " + d.ErrorDetail))
	}

	fmt.Println(tui.DimStyle.Render("Created:            " + d.CreatedAt))
	fmt.Println(tui.DimStyle.Render("Updated:            " + d.UpdatedAt))
}

// PrintDispatchListJSON marshals a list of dispatches as indented JSON.
func PrintDispatchListJSON(items []pipelines.Dispatch) error {
	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(data))

	return nil
}

// PrintDispatchListHuman renders a tabular summary of dispatches.
func PrintDispatchListHuman(items []pipelines.Dispatch) {
	if len(items) == 0 {
		fmt.Println(tui.DimStyle.Render("No dispatches found"))

		return
	}

	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	fmt.Fprintln(writer, "DISPATCH_ID\tSCOPE\tVERSION\tSTATUS\tTRIGGER\tUPDATED")

	for _, d := range items {
		scope := "draft"
		ver := "\u2014"

		if d.VersionID != nil {
			scope = "locked"
			ver = "v" + strconv.Itoa(*d.VersionID)
		}

		fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s\t%s\n",
			d.DispatchID, scope, ver, d.Status, d.TriggeredBy, d.UpdatedAt,
		)
	}

	_ = writer.Flush()
}

// PrintStatusJSON marshals a lightweight status response as indented JSON.
func PrintStatusJSON(s pipelines.DispatchStatus) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(data))

	return nil
}

// PrintStatusHuman renders a lightweight status response.
func PrintStatusHuman(s pipelines.DispatchStatus) {
	covalent := s.CovalentDispatchID
	if covalent == "" {
		covalent = "\u2014"
	}

	fmt.Println(tui.BaseTextStyle.Render("Dispatch ID:        " + s.DispatchID))
	fmt.Println(tui.BaseTextStyle.Render("Status:             " + s.Status))
	fmt.Println(tui.BaseTextStyle.Render("Covalent Dispatch:  " + covalent))
}

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

package update

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"testing"
	"time"

	"github.com/datarobot/cli/internal/pipelines"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	fn()

	w.Close()

	os.Stdout = old

	var buf bytes.Buffer

	_, _ = io.Copy(&buf, r)

	return buf.String()
}

func sampleUpdateResponse() pipelines.CreateResponse {
	return pipelines.CreateResponse{
		PipelineID:    "6658f441-a8f5-4f21-b4d8-6cccf4c94c5b",
		Name:          "confluence_to_vdb",
		Version:       2,
		Status:        "READY",
		Mode:          "draft",
		ElectronNames: []string{"create_vector_database"},
		CreatedAt:     time.Date(2026, 4, 28, 12, 24, 54, 0, time.UTC),
	}
}

func TestPrintUpdateJSON(t *testing.T) {
	resp := sampleUpdateResponse()

	output := captureStdout(t, func() {
		err := printUpdateJSON(resp)
		require.NoError(t, err)
	})

	var parsed map[string]interface{}

	err := json.Unmarshal([]byte(output), &parsed)
	require.NoError(t, err)
	assert.Equal(t, resp.PipelineID, parsed["pipeline_id"])
	assert.EqualValues(t, 2, parsed["version"])
	assert.Equal(t, "READY", parsed["status"])
}

func TestPrintUpdateHuman_WithElectrons(t *testing.T) {
	resp := sampleUpdateResponse()

	output := captureStdout(t, func() {
		printUpdateHuman(resp)
	})

	assert.Contains(t, output, "Pipeline:  "+resp.PipelineID)
	assert.Contains(t, output, "Name:      confluence_to_vdb")
	assert.Contains(t, output, "Version:   2")
	assert.Contains(t, output, "Status:    READY")
	assert.Contains(t, output, "Mode:      draft")
	assert.Contains(t, output, "create_vector_database")
}

func TestPrintUpdateHuman_NoElectrons(t *testing.T) {
	resp := sampleUpdateResponse()
	resp.ElectronNames = nil

	output := captureStdout(t, func() {
		printUpdateHuman(resp)
	})

	assert.Contains(t, output, "Electrons: \u2014")
}

func TestCmd_RequiresTwoArgs(t *testing.T) {
	cmd := Cmd()
	cmd.SetArgs([]string{"only-one"})
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)

	err := cmd.Execute()
	require.Error(t, err)
}

func TestCmd_RejectsInvalidOutput(t *testing.T) {
	cmd := Cmd()
	cmd.SetArgs([]string{"some-id", "some-file.py", "--output", "yaml"})
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	cmd.PreRunE = nil

	err := cmd.Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid output format")
}

func TestCmd_HasOutputFlag(t *testing.T) {
	cmd := Cmd()
	assert.NotNil(t, cmd.Flags().Lookup("output"))
}

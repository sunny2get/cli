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

package graph

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/datarobot/cli/internal/drapi"
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

func sampleGraph() pipelines.Graph {
	return pipelines.Graph{
		Pipeline: pipelines.GraphPipeline{ID: "pipeline-0", Name: "wf"},
		Nodes: []pipelines.GraphNode{
			{ID: "pipeline-0", Type: "lattice", Name: "wf"},
			{ID: "task-0", Type: "electron", Name: "step1"},
		},
		Edges: []pipelines.GraphEdge{
			{Source: "pipeline-0", Target: "task-0"},
		},
	}
}

func TestPrintGraphJSON(t *testing.T) {
	output := captureStdout(t, func() {
		require.NoError(t, printGraphJSON(sampleGraph()))
	})

	var parsed map[string]any

	require.NoError(t, json.Unmarshal([]byte(output), &parsed))

	// JSON tag remains `lattice` until the API renames the wire field.
	pipeline, ok := parsed["lattice"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "pipeline-0", pipeline["id"])
	assert.Equal(t, "wf", pipeline["name"])
}

func TestPrintGraphHuman(t *testing.T) {
	output := captureStdout(t, func() {
		printGraphHuman(sampleGraph())
	})

	assert.Contains(t, output, "Pipeline: wf")
	assert.Contains(t, output, "ID:       pipeline-0")
	assert.Contains(t, output, "Nodes (2):")
	assert.Contains(t, output, "Edges (1):")
	assert.Contains(t, output, "task-0")
}

func TestPrintGraphHuman_EmptyGraph(t *testing.T) {
	output := captureStdout(t, func() {
		printGraphHuman(pipelines.Graph{Pipeline: pipelines.GraphPipeline{Name: "empty", ID: "x"}})
	})

	assert.Contains(t, output, "No nodes")
}

func TestHandleGraphError_404IsSuppressed(t *testing.T) {
	httpErr := &drapi.HTTPError{StatusCode: http.StatusNotFound, URL: "x"}

	err := handleGraphError(httpErr, "abc")
	assert.NoError(t, err)
}

func TestHandleGraphError_PropagatesOther(t *testing.T) {
	err := handleGraphError(errors.New("boom"), "abc")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "boom")
}

func TestCmd_RejectsMissingPipeline(t *testing.T) {
	cmd := Cmd()
	cmd.SetArgs([]string{})
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	cmd.PreRunE = nil

	err := cmd.Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "--pipeline")
}

func TestCmd_RejectsBadOutput(t *testing.T) {
	cmd := Cmd()
	cmd.SetArgs([]string{"--pipeline", "p", "--output", "yaml"})
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	cmd.PreRunE = nil

	err := cmd.Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid output format")
}

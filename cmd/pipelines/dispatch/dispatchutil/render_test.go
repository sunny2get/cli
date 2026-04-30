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

package dispatchutil

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"testing"

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

func intPtr(v int) *int { return &v }

func sampleDraftDispatch() pipelines.Dispatch {
	return pipelines.Dispatch{
		DispatchID:  "d-1",
		PipelineID:  "pl-1",
		InputID:     "in-1",
		TriggeredBy: "user@example.com",
		Status:      pipelines.DispatchStatusPending,
		CreatedAt:   "2026-04-29T10:00:00Z",
		UpdatedAt:   "2026-04-29T10:00:00Z",
	}
}

func sampleLockedDispatch() pipelines.Dispatch {
	d := sampleDraftDispatch()
	d.VersionID = intPtr(3)
	d.CovalentDispatchID = "cov-xyz"
	d.Status = pipelines.DispatchStatusRunning

	return d
}

func TestPrintDispatchJSON(t *testing.T) {
	output := captureStdout(t, func() {
		require.NoError(t, PrintDispatchJSON(sampleDraftDispatch()))
	})

	var parsed map[string]any

	require.NoError(t, json.Unmarshal([]byte(output), &parsed))
	assert.Equal(t, "d-1", parsed["dispatch_id"])
	assert.Equal(t, "PENDING", parsed["status"])
}

func TestPrintDispatchHuman_DraftMissingCovalent(t *testing.T) {
	output := captureStdout(t, func() { PrintDispatchHuman(sampleDraftDispatch()) })
	assert.Contains(t, output, "Dispatch ID:        d-1")
	assert.Contains(t, output, "Scope:              draft")
	assert.Contains(t, output, "Version:            \u2014")
	assert.Contains(t, output, "Covalent Dispatch:  \u2014")
	assert.Contains(t, output, "Status:             PENDING")
}

func TestPrintDispatchHuman_LockedShowsErrorWhenSet(t *testing.T) {
	d := sampleLockedDispatch()
	d.Status = pipelines.DispatchStatusFailed
	d.ErrorDetail = "boom"

	output := captureStdout(t, func() { PrintDispatchHuman(d) })
	assert.Contains(t, output, "Scope:              locked")
	assert.Contains(t, output, "Version:            v3")
	assert.Contains(t, output, "Covalent Dispatch:  cov-xyz")
	assert.Contains(t, output, "Error:              boom")
}

func TestPrintDispatchListJSON(t *testing.T) {
	output := captureStdout(t, func() {
		require.NoError(t, PrintDispatchListJSON([]pipelines.Dispatch{sampleDraftDispatch()}))
	})

	var parsed []map[string]any

	require.NoError(t, json.Unmarshal([]byte(output), &parsed))
	require.Len(t, parsed, 1)
	assert.Equal(t, "d-1", parsed[0]["dispatch_id"])
}

func TestPrintDispatchListHuman_Empty(t *testing.T) {
	output := captureStdout(t, func() { PrintDispatchListHuman(nil) })
	assert.Contains(t, output, "No dispatches found")
}

func TestPrintDispatchListHuman_RendersTable(t *testing.T) {
	output := captureStdout(t, func() {
		PrintDispatchListHuman([]pipelines.Dispatch{sampleDraftDispatch(), sampleLockedDispatch()})
	})

	assert.Contains(t, output, "DISPATCH_ID")
	assert.Contains(t, output, "STATUS")
	assert.Contains(t, output, "TRIGGER")
	assert.Contains(t, output, "draft")
	assert.Contains(t, output, "locked")
	assert.Contains(t, output, "v3")
	assert.Contains(t, output, "PENDING")
	assert.Contains(t, output, "RUNNING")
}

func TestPrintStatusJSON(t *testing.T) {
	status := pipelines.DispatchStatus{
		DispatchID:         "d-1",
		Status:             pipelines.DispatchStatusCompleted,
		CovalentDispatchID: "cov-xyz",
	}

	output := captureStdout(t, func() {
		require.NoError(t, PrintStatusJSON(status))
	})

	var parsed map[string]any

	require.NoError(t, json.Unmarshal([]byte(output), &parsed))
	assert.Equal(t, "COMPLETED", parsed["status"])
	assert.Equal(t, "cov-xyz", parsed["covalent_dispatch_id"])
}

func TestPrintStatusHuman_NoCovalentDispatchID(t *testing.T) {
	output := captureStdout(t, func() {
		PrintStatusHuman(pipelines.DispatchStatus{DispatchID: "d-1", Status: "PENDING"})
	})

	assert.Contains(t, output, "Dispatch ID:        d-1")
	assert.Contains(t, output, "Status:             PENDING")
	assert.Contains(t, output, "Covalent Dispatch:  \u2014")
}

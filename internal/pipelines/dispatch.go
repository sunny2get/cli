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

// dispatch.go wraps the pipeline dispatch endpoints described in
// pipelines-api/.../controllers/pipeline_dispatch.py. The CLI exposes the
// same draft/locked URL split as inputs via the shared Scope helpers.

package pipelines

import (
	"net/http"
	"net/url"
	"strconv"
)

// Dispatch lifecycle states (mirrors PipelineDispatchStatus).
const (
	DispatchStatusPending   = "PENDING"
	DispatchStatusRunning   = "RUNNING"
	DispatchStatusCompleted = "COMPLETED"
	DispatchStatusFailed    = "FAILED"
	DispatchStatusCancelled = "CANCELLED"
	DispatchStatusErrored   = "ERRORED"
)

// Dispatch mirrors PipelineDispatchResponse from the pipelines-api.
type Dispatch struct {
	DispatchID         string `json:"dispatch_id"`
	PipelineID         string `json:"pipeline_id"`
	VersionID          *int   `json:"version_id,omitempty"`
	InputID            string `json:"input_id"`
	CovalentDispatchID string `json:"covalent_dispatch_id,omitempty"`
	TriggeredBy        string `json:"triggered_by"`
	Status             string `json:"status"`
	ErrorDetail        string `json:"error_detail,omitempty"`
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_at"`
}

// DispatchStatus mirrors PipelineDispatchStatusResponse — the lightweight
// polling-friendly shape returned by GET .../status.
type DispatchStatus struct {
	DispatchID         string `json:"dispatch_id"`
	Status             string `json:"status"`
	CovalentDispatchID string `json:"covalent_dispatch_id,omitempty"`
}

// DispatchCreateRequest mirrors PipelineDispatchCreateRequest.
type DispatchCreateRequest struct {
	InputID string `json:"input_id"`
}

// CreateDispatch starts a new dispatch for the given input. Returns the
// freshly-created Dispatch (status PENDING).
func CreateDispatch(pipelineID string, scope Scope, version *int, inputID string) (*Dispatch, error) {
	endpoint, err := EndpointFor(pipelineID, scope, version, "dispatches")
	if err != nil {
		return nil, err
	}

	body := DispatchCreateRequest{InputID: inputID}

	var result Dispatch

	err = doJSON(http.MethodPost, endpoint, body, "create dispatch", &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// ListDispatches returns a paginated slice of dispatches for the given scope.
func ListDispatches(pipelineID string, scope Scope, version *int, offset, limit int) ([]Dispatch, error) {
	endpoint, err := EndpointFor(pipelineID, scope, version, "dispatches")
	if err != nil {
		return nil, err
	}

	query := url.Values{}
	if offset > 0 {
		query.Set("offset", strconv.Itoa(offset))
	}

	if limit > 0 {
		query.Set("limit", strconv.Itoa(limit))
	}

	if encoded := query.Encode(); encoded != "" {
		endpoint = endpoint + "?" + encoded
	}

	var dispatches []Dispatch

	err = doJSON(http.MethodGet, endpoint, nil, "dispatches", &dispatches)
	if err != nil {
		return nil, err
	}

	return dispatches, nil
}

// GetDispatch fetches a single dispatch by id within the given scope.
func GetDispatch(pipelineID string, scope Scope, version *int, dispatchID string) (*Dispatch, error) {
	endpoint, err := EndpointFor(pipelineID, scope, version, "dispatches/"+dispatchID)
	if err != nil {
		return nil, err
	}

	var dispatch Dispatch

	err = doJSON(http.MethodGet, endpoint, nil, "dispatch", &dispatch)
	if err != nil {
		return nil, err
	}

	return &dispatch, nil
}

// GetDispatchStatus calls the lightweight GET .../status endpoint useful
// for polling without re-downloading the full dispatch record.
func GetDispatchStatus(pipelineID string, scope Scope, version *int, dispatchID string) (*DispatchStatus, error) {
	endpoint, err := EndpointFor(pipelineID, scope, version, "dispatches/"+dispatchID+"/status")
	if err != nil {
		return nil, err
	}

	var status DispatchStatus

	err = doJSON(http.MethodGet, endpoint, nil, "dispatch status", &status)
	if err != nil {
		return nil, err
	}

	return &status, nil
}

// CancelDispatch issues a DELETE on a dispatch, transitioning it to
// CANCELLED if it is still in a non-terminal state.
func CancelDispatch(pipelineID string, scope Scope, version *int, dispatchID string) error {
	endpoint, err := EndpointFor(pipelineID, scope, version, "dispatches/"+dispatchID)
	if err != nil {
		return err
	}

	return doDelete(endpoint, "cancel dispatch")
}

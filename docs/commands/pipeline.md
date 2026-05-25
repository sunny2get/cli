# `dr pipeline` - Pipelines API management

Manage AI/ML pipelines orchestrated by Covalent through the DataRobot
pipelines service. The `dr pipeline` group is a thin CLI wrapper over
the pipelines REST API: every subcommand maps directly to a single
endpoint.

## Quick start

```bash
# List pipelines registered with the pipelines service
dr pipeline list

# Register a new draft pipeline by uploading a DataRobot pipeline source file
dr pipeline create ./my_pipeline.py --description "First draft"

# Append a new version after editing the file
dr pipeline update <pipeline-id> ./my_pipeline.py

# Promote the draft to locked when you are happy with it
dr pipeline lock <pipeline-id>

# Cancel a stuck run
dr pipeline run cancel --pipeline <pipeline-id> <run-id>
```

> [!NOTE]
> The `pipeline` command is currently behind a feature gate. Enable it
> by exporting `DATAROBOT_CLI_FEATURE_PIPELINE=true` before running any
> `dr pipeline` subcommand. See
> [Feature gates](../development/feature-gates.md) for details.

> [!NOTE]
> **First time?** If you're new to the CLI, start with the
> [Quick start](../../README.md#quick-start) for step-by-step setup
> instructions.

## Synopsis

```bash
dr pipeline <command> [subcommand] [flags]
```

## Description

A **pipeline** is a versioned bundle of Python source defining a DataRobot pipeline (one or more tasks). Each
top-level `dr pipeline` subcommand operates on one of four resources:

- the **pipeline** itself (create, list, get, update, delete, lock),
- pipeline **versions** (list, get, graph),
- pipeline **inputs** — JSON payloads supplied to a run,
- pipeline **runs** — concrete executions on Covalent,
- pipeline **schedules** — recurring runs on a cron expression,
- pipeline **environments** — named, immutable-versioned bags of pip
  packages that pipelines can be built against.

Versions are created automatically:

- The first `create` call registers the source as **v1** in `draft`
  mode.
- `update` re-uploads the same file (or an edited copy) and appends
  **v2**, **v3**, etc., as long as the pipeline name still matches and
  the pipeline is still in `draft` mode.
- `lock` promotes a draft to **locked** mode. Locked pipelines are
  immutable; their inputs and schedules become valid.

Inputs, runs, and the graph endpoint exist in two scopes —
**draft** (mutable, no version pinned) and **locked** (immutable, tied
to a frozen version) — selected via the shared `--scope` and
`--version` flags. Schedules are locked-only.

> For an exhaustive table mapping every CLI command to its API endpoint,
> see [pipeline-reference.md](pipeline-reference.md).

## Command groups

| Group                        | Endpoint(s)                                      | Purpose                                                |
|------------------------------|--------------------------------------------------|--------------------------------------------------------|
| `dr pipeline create`        | `POST   /api/v2/pipelines`                       | Upload a Python file to register a new pipeline.       |
| `dr pipeline list`          | `GET    /api/v2/pipelines`                       | Paginated list with mode filtering.                    |
| `dr pipeline get`           | `GET    /api/v2/pipelines/{id}`                  | Pipeline detail including all versions.                |
| `dr pipeline update`        | `PATCH  /api/v2/pipelines/{id}`                  | Re-upload a file to append a new version.              |
| `dr pipeline delete`        | `DELETE /api/v2/pipelines/{id}`                  | Remove a pipeline and all of its versions.             |
| `dr pipeline lock`          | `PATCH  /api/v2/pipelines/{id}/mode`             | Promote a draft to locked mode.                        |
| `dr pipeline version …`     | `…/versions[/{ver}]`                             | Inspect pipeline versions.                             |
| `dr pipeline graph`         | `…/graph` (draft) or `…/versions/{ver}/graph`    | Render the pipeline/task DAG.                          |
| `dr pipeline input …`       | `…/inputs` and `…/inputs/{input_id}`             | Manage JSON payloads for runs.                         |
| `dr pipeline run …`         | `…/dispatches` and `…/dispatches/{dispatch_id}`  | Trigger, inspect, and cancel runs.                     |
| `dr pipeline schedule …`    | `…/versions/{ver}/schedules`                     | Manage recurring (cron) runs on locked versions.       |
| `dr pipeline environment …` | `/api/v2/pipelines/environments[/...]`           | Manage named, versioned pip-package execution environments. |

## Subcommands

### `create`

Upload a Python file defining a DataRobot pipeline (one or more tasks) and
register a new pipeline. The pipeline name is extracted from the file and used as the
pipeline name.

```bash
dr pipeline create <file> [flags]
dr pipeline create --from-file=<file> [flags]
```

**Arguments:**

- `<file>` — path to a `.py` file containing a single DataRobot pipeline.
  Mutually exclusive with `--from-file`.

**Flags:**

- `--from-file <path>` — alternative to the positional file argument.
- `--description <text>` — optional human-readable description stored on
  the pipeline.
- `--mode <draft|locked>` — pipeline lifecycle mode. Defaults to `draft`.
- `--output <json>` — emit machine-parseable JSON instead of the
  human-readable summary.

**Example:**

```bash
$ dr pipeline create ./confluence_to_vdb.py --description "test"
Pipeline ID:  683c2a1b4f8e1a2b3c4d5e6f
Name:         confluence_to_vdb
Version:      1
Status:       READY
Mode:         draft
Tasks:        create_vector_database, ingest_confluence_files, setup_credential_and_datastore
Created:      2026-04-28T11:42:28Z
```

### `list`

List pipelines registered with the pipelines service, with optional
mode filtering and pagination.

```bash
dr pipeline list [flags]
```

**Flags:**

- `--mode <draft|locked>` — filter by pipeline mode.
- `--offset <N>` — pagination offset. Default `0`.
- `--limit <N>` — pagination limit (1-200). Default `50`.
- `--output <json>` — emit machine-parseable JSON instead of a table.

**Example:**

```bash
$ dr pipeline list
Showing 1 of 1 (offset=0 limit=50)

ID                                    NAME               MODE   ACTIVE  VERSION  UPDATED
683c2a1b4f8e1a2b3c4d5e6f  confluence_to_vdb  draft  true    v3       2026-04-28T12:25:11Z
```

| Column    | Meaning                                                         |
|-----------|-----------------------------------------------------------------|
| `ID`      | Pipeline ObjectId, used as the argument to `get` / `update` / etc.  |
| `NAME`    | Pipeline name extracted from the originally uploaded file.     |
| `MODE`    | `draft` (mutable) or `locked` (immutable).                      |
| `ACTIVE`  | `true` while the pipeline has not been soft-deleted.            |
| `VERSION` | Latest version number, or `—` when no versions exist yet.       |
| `UPDATED` | Last modification time in UTC (RFC 3339).                       |

### `get`

Display full details of a single pipeline including all versions.

```bash
dr pipeline get <pipeline-id> [flags]
```

**Arguments:**

- `<pipeline-id>` — the ObjectId returned by `create` / shown in `pipeline list`.

**Flags:**

- `--output <json>` — emit machine-parseable JSON.

**Example:**

```bash
$ dr pipeline get 683c2a1b4f8e1a2b3c4d5e6f
ID:          683c2a1b4f8e1a2b3c4d5e6f
Name:        confluence_to_vdb
Description: test
Mode:        draft
Active:      true
Created:     2026-04-28T11:42:28Z
Updated:     2026-04-28T12:25:11Z

Versions (3):
  VERSION  STATUS  PYTHON  CREATED               TASKS
  v1       READY   3.12    2026-04-28T11:42:28Z  create_vector_database, ingest_confluence_files, setup_credential_and_datastore
  v2       READY   3.12    2026-04-28T12:24:54Z  create_vector_database, ingest_confluence_files, setup_credential_and_datastore
  v3       READY   3.12    2026-04-28T12:25:11Z  create_vector_database, ingest_confluence_files, setup_credential_and_datastore
```

If a version failed to register, its `error_detail` is rendered as a
dim line underneath the table.

If the pipeline doesn't exist, `get` prints
`No pipeline found with id: <id>` and exits 0 instead of dumping an
HTTP error.

### `update`

Re-upload a Python file to update a draft pipeline. A new version is
appended.

```bash
dr pipeline update <pipeline-id> <file> [flags]
dr pipeline update <pipeline-id> --from-file=<file> [flags]
```

**Arguments:**

- `<pipeline-id>` — the ObjectId of the pipeline to update.
- `<file>` — path to the updated `.py` file. Mutually exclusive with
  `--from-file`.

**Flags:**

- `--from-file <path>` — alternative to the positional file argument.
- `--output <json>` — emit machine-parseable JSON.

**Constraints:**

- The pipeline name encoded in the uploaded file **must match** the pipeline's
  existing name. To register a different pipeline, use `create` instead.
- Locked pipelines cannot be updated. The API responds with
  `409 Conflict`.

### `delete`

Delete a pipeline and all of its versions.

```bash
dr pipeline delete <pipeline-id>
```

**Arguments:**

- `<pipeline-id>` — the ObjectId of the pipeline to delete.

**Example:**

```bash
$ dr pipeline delete 683c2a1b4f8e1a2b3c4d5e6f
Deleted pipeline: 683c2a1b4f8e1a2b3c4d5e6f
```

If the pipeline doesn't exist, `delete` prints
`No pipeline found with id: <id>` and exits 0.

### `lock`

Promote a draft pipeline to locked mode. Once locked, the pipeline can
no longer be updated and locked runs/inputs/schedules become
valid.

```bash
dr pipeline lock <pipeline-id> [flags]
```

**Arguments:**

- `<pipeline-id>` — the ObjectId of the pipeline to lock.

**Flags:**

- `--output <json>` — emit machine-parseable JSON.

**Example:**

```bash
$ dr pipeline lock 683c2a1b4f8e1a2b3c4d5e6f
Pipeline ID:  683c2a1b4f8e1a2b3c4d5e6f
Name:         confluence_to_vdb
Mode:         locked
Version:      v3
Status:       READY
Tasks:        create_vector_database, ingest_confluence_files, setup_credential_and_datastore
Locked:       2026-04-28T12:30:00Z
```

### `version`

Read-only access to pipeline versions.

```bash
dr pipeline version list --pipeline <id> [--offset N] [--limit N] [--output json]
dr pipeline version get  --pipeline <id> <version-id>     [--output json]
```

`version list` returns the same data that's shown inline by
`pipeline get`, but in a paginated stand-alone view. `version get`
shows a single version in detail (pipeline, tasks, Python version,
creation timestamp, and any error detail).

### `graph`

Display the pipeline/task DAG for a pipeline as either a JSON
payload (for visualisation tooling) or a human-readable summary.

```bash
dr pipeline graph --pipeline <id>                       # draft graph
dr pipeline graph --pipeline <id> --version=N           # locked-version graph
dr pipeline graph --pipeline <id> --scope=draft         # explicit draft
dr pipeline graph --pipeline <id> --scope=locked --version=N
dr pipeline graph --pipeline <id> --output json         # JSON payload
```

The human view prints the pipeline header followed by `Nodes (N):` and
`Edges (M):` tables. Pass `--output json` to get the structured `Graph`
object (`lattice`, `nodes[]`, `edges[]` — JSON keys preserved while the API
wire format is unchanged) suitable for piping to
visualisation tooling.

See [Shared `--scope` / `--version` semantics](#scope--version-flags)
below for the full flag truth table.

### `input`

Manage JSON payloads that drive a run.

```bash
dr pipeline input create --pipeline <id> <payload-file>            # draft scope
dr pipeline input create --pipeline <id> --version=N <payload-file>  # locked scope
dr pipeline input list   --pipeline <id> [--scope|--version] [--offset N] [--limit N]
dr pipeline input get    --pipeline <id> <input-id>      [--scope|--version]
dr pipeline input update --pipeline <id> <input-id> <payload-file>   # draft only
dr pipeline input delete --pipeline <id> <input-id>      [--scope|--version]
```

- The payload file must contain a JSON object. The CLI wraps it in
  `{"payload": …}` before sending.
- All verbs accept `--scope` / `--version` (see below). Inputs in the
  `locked` scope are immutable, so `input update` is draft-only.

### `run`

Trigger, inspect, and cancel pipeline executions.

```bash
dr pipeline run create --pipeline <id> --input <input-id>          # draft
dr pipeline run create --pipeline <id> --version=N --input <input-id>  # locked
dr pipeline run list   --pipeline <id> [--scope|--version]
dr pipeline run get    --pipeline <id> <run-id> [--scope|--version]
dr pipeline run status --pipeline <id> <run-id> [--scope|--version]
dr pipeline run cancel --pipeline <id> <run-id> [--scope|--version]
```

`run status` is a lighter-weight call than `run get` —
intended for polling — and returns just the run ID, status, and
the corresponding Covalent dispatch ID.

`run cancel` returns `409 Conflict` if the run is already in
a terminal state (COMPLETED / FAILED / CANCELLED).

### `schedule`

Manage recurring (cron) runs on locked versions only. Both
`--pipeline` and `--version` are required for every verb.

```bash
dr pipeline schedule create --pipeline <id> --version=N \
    --cron "0 * * * *" --input <input-id> [--timezone UTC]
dr pipeline schedule list   --pipeline <id> --version=N [--offset N] [--limit N]
dr pipeline schedule get    --pipeline <id> --version=N <schedule-id>
dr pipeline schedule update --pipeline <id> --version=N <schedule-id> [--cron "*/15 * * * *"] [--timezone Europe/Berlin]
dr pipeline schedule delete --pipeline <id> --version=N <schedule-id>
```

`schedule update` requires at least one of `--cron` or `--timezone`.

### `environment`

Manage pipeline execution environments — named, immutable-versioned
bags of pip packages that pipelines can be built against. Environments
live at the top of the pipelines namespace (not nested under a specific
pipeline) and have their own lifecycle.

```bash
dr pipeline environment create --name <name> --package <spec> [--package <spec>] ...
dr pipeline environment list   [--offset N] [--limit N] [--output json]
dr pipeline environment update <environment-id> --package <spec> [...]
dr pipeline environment delete <environment-id>
dr pipeline environment version delete --environment <id> <version>
```

`create` registers a new environment with an initial v1 containing the
supplied pip packages; the returned record reports the build status of
that first version. `update` adds packages to an existing environment
by creating a new immutable version (older versions are unchanged).
`delete` soft-deletes the most recent active version (and cascades the
parent if no versions remain). `version delete` targets a specific
older version without touching the parent.

`--package` is repeatable and also accepts comma-separated values:

```bash
dr pipeline environment create --name ml-base \
    --package numpy --package pandas==2.0
dr pipeline environment create --name ml-base \
    --package "numpy,pandas==2.0,scikit-learn"
```

> [!NOTE]
> The pipelines-api currently does not surface `GET` endpoints for a
> single environment or for the version list. The full version history
> is only returned in the `create` and `update` responses.

## Shared flags

### `--scope` / `--version` flags

Inputs, runs, and `graph` mirror the API's two URL shapes —
`/pipelines/{id}/…` for the mutable draft and
`/pipelines/{id}/versions/{ver}/…` for a locked version — through a
pair of optional flags:

| Flags supplied                       | Resolved scope    | URL used                            |
|--------------------------------------|-------------------|-------------------------------------|
| _(none)_                             | `draft`           | `/pipelines/{id}/…`                 |
| `--version=N`                        | `locked` (auto)   | `/pipelines/{id}/versions/N/…`      |
| `--scope=draft`                      | `draft`           | `/pipelines/{id}/…`                 |
| `--scope=locked --version=N`         | `locked`          | `/pipelines/{id}/versions/N/…`      |
| `--scope=draft --version=N`          | **error**         | `--scope=draft cannot be combined with --version` |
| `--scope=locked` (no `--version`)    | **error**         | `--scope=locked requires --version=<n>` |
| `--scope=garbage`                    | **error**         | `invalid --scope: "garbage" (supported: draft, locked)` |

Schedules do not accept `--scope`; they are locked-only and require
`--version` on every verb.

### `--from-file` / positional file

`pipeline create`, `pipeline update`, `pipeline input create`, and
`pipeline input update` all accept the input file in two equivalent
ways:

```bash
dr pipeline create ./my_pipeline.py
dr pipeline create --from-file=./my_pipeline.py
```

Exactly one of the two must be supplied; passing both yields
`specify the file either as a positional argument or via --from-file, not both`,
and supplying neither yields `a file path is required …` (or
`a JSON payload file is required …` for input verbs).

### `--output`

Every read/write verb that produces a payload accepts `--output json`
to emit the underlying response struct as indented JSON. Any other
value is rejected with
`invalid output format: <value> (supported: json)`.

```bash
dr pipeline list --output json | jq '.items[].pipeline_id'
```

### Global options

All [global flags](README.md#global-flags) are available, notably
`--debug` for protocol-level tracing and `--skip-auth` for advanced
scenarios.

## Local development

While iterating against a locally running pipelines-api (default port
`8100`), point the CLI at `http://localhost:8100` and bypass token
verification using the prefixed environment variables:

```bash
export DATAROBOT_CLI_FEATURE_PIPELINE=true
export DATAROBOT_CLI_ENDPOINT=http://localhost:8100/api/v2
export DATAROBOT_CLI_TOKEN=local
export DATAROBOT_CLI_SKIP_AUTH=true

./dist/dr pipeline list
```

Why each variable matters:

| Variable                          | Purpose                                                                                                  |
|-----------------------------------|----------------------------------------------------------------------------------------------------------|
| `DATAROBOT_CLI_FEATURE_PIPELINE` | Reveals the feature-gated `pipeline` command group.                                                     |
| `DATAROBOT_CLI_ENDPOINT`          | Auto-bound to viper's `endpoint` key via `SetEnvPrefix("DATAROBOT_CLI") + AutomaticEnv()` in the CLI.    |
| `DATAROBOT_CLI_TOKEN`             | Same prefix story — bound to viper's `token` key for outbound `Authorization: Bearer` headers.           |
| `DATAROBOT_CLI_SKIP_AUTH`         | Skips token verification against `<endpoint>/version/`, which the local stub does not implement.         |

The unprefixed `DATAROBOT_ENDPOINT` / `DATAROBOT_API_TOKEN` variables
are **only** bound to viper after a successful token verification and
therefore do not work alongside `--skip-auth` or
`DATAROBOT_CLI_SKIP_AUTH=true`. Always prefer the `DATAROBOT_CLI_*`
names during local development.

> [!TIP]
> Use `--debug` to see the full request URL, headers, and body the CLI
> sends. Logs are written to `.dr-tui-debug.log`.

For a step-by-step walkthrough of how `dr pipeline list` was wired up,
see [Adding a command](../development/adding-a-command.md).

## Examples

### Pipeline lifecycle

```bash
# Register a draft, append a version, lock it, then delete it
dr pipeline create ./my_pipeline.py --description "Initial draft"
dr pipeline update <pipeline-id> ./my_pipeline.py
dr pipeline lock   <pipeline-id>
dr pipeline delete <pipeline-id>
```

### Inspect versions and graph

```bash
dr pipeline version list --pipeline <pipeline-id>
dr pipeline version get  --pipeline <pipeline-id> 2
dr pipeline graph        --pipeline <pipeline-id> --version=2 --output json
```

### Trigger a run

```bash
# 1. Register a JSON input on the draft scope
dr pipeline input create --pipeline <pipeline-id> ./input.json

# 2. Trigger a run with that input
dr pipeline run create --pipeline <pipeline-id> --input <input-id>

# 3. Poll status until it reaches a terminal state
dr pipeline run status --pipeline <pipeline-id> <run-id>
```

### Schedule a recurring run on a locked version

```bash
dr pipeline schedule create \
    --pipeline <pipeline-id> --version=2 \
    --cron "0 */6 * * *" --input <input-id> --timezone America/Los_Angeles
```

### Scripting friendliness

```bash
# All UUIDs of locked pipelines, one per line
dr pipeline list --mode locked --output json | jq -r '.items[].pipeline_id'
```

## Error handling

The CLI surfaces backend errors verbatim and exits non-zero. The most
common status codes you will see:

| Status | Cause                                                                              |
|--------|------------------------------------------------------------------------------------|
| `400`  | Invalid Python file, mismatched pipeline name, or malformed JSON payload.          |
| `404`  | The provided `<pipeline-id>` / version / input / run / schedule does not exist.    |
| `409`  | Tried to update a `locked` pipeline, or cancel an already-terminal run.            |

For most `get` / `delete` verbs the CLI translates a 404 into a
friendly informational line (e.g. `No pipeline found with id: …`) and
exits 0, so a no-op delete won't dump usage at the user.

## See also

- [Pipelines reference](pipeline-reference.md) — exhaustive table
  mapping every CLI command to its API endpoint, usage variants, and
  inputs.
- [Authentication](auth.md) — how `dr auth login` and `--skip-auth`
  interact.
- [Configuration](../user-guide/configuration.md) — config file and
  environment-variable precedence.
- [Adding a command](../development/adding-a-command.md) — how the
  pipelines verbs were built.
- [Feature gates](../development/feature-gates.md) — flipping
  `DATAROBOT_CLI_FEATURE_PIPELINE` on and off.

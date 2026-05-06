# `dr pipelines` - Pipelines API management

Manage AI/ML pipelines orchestrated by Covalent through the DataRobot
pipelines service. The `dr pipelines` group is a thin CLI wrapper over
the pipelines REST API: every subcommand maps directly to a single
endpoint.

## Quick start

```bash
# List pipelines registered with the pipelines service
dr pipelines list

# Register a new draft pipeline by uploading a DataRobot pipeline source file
dr pipelines create ./my_pipeline.py --description "First draft"

# Append a new version after editing the file
dr pipelines update <pipeline-id> ./my_pipeline.py

# Promote the draft to locked when you are happy with it
dr pipelines lock <pipeline-id>

# Cancel a stuck dispatch
dr pipelines dispatch cancel --pipeline <pipeline-id> <dispatch-id>
```

> [!NOTE]
> The `pipelines` command is currently behind a feature gate. Enable it
> by exporting `DATAROBOT_CLI_FEATURE_PIPELINES=true` before running any
> `dr pipelines` subcommand. See
> [Feature gates](../development/feature-gates.md) for details.

> [!NOTE]
> **First time?** If you're new to the CLI, start with the
> [Quick start](../../README.md#quick-start) for step-by-step setup
> instructions.

## Synopsis

```bash
dr pipelines <command> [subcommand] [flags]
```

## Description

A **pipeline** is a versioned bundle of Python source defining a DataRobot pipeline (one or more tasks). Each
top-level `dr pipelines` subcommand operates on one of four resources:

- the **pipeline** itself (create, list, get, update, delete, lock),
- pipeline **versions** (list, get, graph),
- pipeline **inputs** — JSON payloads supplied to a dispatch,
- pipeline **dispatches** — concrete executions on Covalent,
- pipeline **schedules** — recurring dispatches on a cron expression,
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

Inputs, dispatches, and the graph endpoint exist in two scopes —
**draft** (mutable, no version pinned) and **locked** (immutable, tied
to a frozen version) — selected via the shared `--scope` and
`--version` flags. Schedules are locked-only.

> For an exhaustive table mapping every CLI command to its API endpoint,
> see [pipelines-reference.md](pipelines-reference.md).

## Command groups

| Group                        | Endpoint(s)                                      | Purpose                                                |
|------------------------------|--------------------------------------------------|--------------------------------------------------------|
| `dr pipelines create`        | `POST   /api/v2/pipelines`                       | Upload a Python file to register a new pipeline.       |
| `dr pipelines list`          | `GET    /api/v2/pipelines`                       | Paginated list with mode filtering.                    |
| `dr pipelines get`           | `GET    /api/v2/pipelines/{id}`                  | Pipeline detail including all versions.                |
| `dr pipelines update`        | `PATCH  /api/v2/pipelines/{id}`                  | Re-upload a file to append a new version.              |
| `dr pipelines delete`        | `DELETE /api/v2/pipelines/{id}`                  | Remove a pipeline and all of its versions.             |
| `dr pipelines lock`          | `PATCH  /api/v2/pipelines/{id}/mode`             | Promote a draft to locked mode.                        |
| `dr pipelines version …`     | `…/versions[/{ver}]`                             | Inspect pipeline versions.                             |
| `dr pipelines graph`         | `…/graph` (draft) or `…/versions/{ver}/graph`    | Render the pipeline/task DAG.                          |
| `dr pipelines input …`       | `…/inputs` and `…/inputs/{input_id}`             | Manage JSON payloads for dispatches.                   |
| `dr pipelines dispatch …`    | `…/dispatches` and `…/dispatches/{dispatch_id}`  | Trigger, inspect, and cancel runs.                     |
| `dr pipelines schedule …`    | `…/versions/{ver}/schedules`                     | Manage recurring (cron) dispatches on locked versions. |
| `dr pipelines environment …` | `/api/v2/pipelines/environments[/...]`           | Manage named, versioned pip-package execution environments. |

## Subcommands

### `create`

Upload a Python file defining a DataRobot pipeline (one or more tasks) and
register a new pipeline. The pipeline name is extracted from the file and used as the
pipeline name.

```bash
dr pipelines create <file> [flags]
dr pipelines create --from-file=<file> [flags]
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
$ dr pipelines create ./confluence_to_vdb.py --description "test"
Pipeline ID:  6658f441-a8f5-4f21-b4d8-6cccf4c94c5b
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
dr pipelines list [flags]
```

**Flags:**

- `--mode <draft|locked>` — filter by pipeline mode.
- `--offset <N>` — pagination offset. Default `0`.
- `--limit <N>` — pagination limit (1-200). Default `50`.
- `--output <json>` — emit machine-parseable JSON instead of a table.

**Example:**

```bash
$ dr pipelines list
Showing 1 of 1 (offset=0 limit=50)

ID                                    NAME               MODE   ACTIVE  VERSION  UPDATED
6658f441-a8f5-4f21-b4d8-6cccf4c94c5b  confluence_to_vdb  draft  true    v3       2026-04-28T12:25:11Z
```

| Column    | Meaning                                                         |
|-----------|-----------------------------------------------------------------|
| `ID`      | Pipeline UUID, used as the argument to `get` / `update` / etc.  |
| `NAME`    | Pipeline name extracted from the originally uploaded file.     |
| `MODE`    | `draft` (mutable) or `locked` (immutable).                      |
| `ACTIVE`  | `true` while the pipeline has not been soft-deleted.            |
| `VERSION` | Latest version number, or `—` when no versions exist yet.       |
| `UPDATED` | Last modification time in UTC (RFC 3339).                       |

### `get`

Display full details of a single pipeline including all versions.

```bash
dr pipelines get <pipeline-id> [flags]
```

**Arguments:**

- `<pipeline-id>` — the UUID returned by `create` / shown in `pipelines list`.

**Flags:**

- `--output <json>` — emit machine-parseable JSON.

**Example:**

```bash
$ dr pipelines get 6658f441-a8f5-4f21-b4d8-6cccf4c94c5b
ID:          6658f441-a8f5-4f21-b4d8-6cccf4c94c5b
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
dr pipelines update <pipeline-id> <file> [flags]
dr pipelines update <pipeline-id> --from-file=<file> [flags]
```

**Arguments:**

- `<pipeline-id>` — the UUID of the pipeline to update.
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
dr pipelines delete <pipeline-id>
```

**Arguments:**

- `<pipeline-id>` — the UUID of the pipeline to delete.

**Example:**

```bash
$ dr pipelines delete 6658f441-a8f5-4f21-b4d8-6cccf4c94c5b
Deleted pipeline: 6658f441-a8f5-4f21-b4d8-6cccf4c94c5b
```

If the pipeline doesn't exist, `delete` prints
`No pipeline found with id: <id>` and exits 0.

### `lock`

Promote a draft pipeline to locked mode. Once locked, the pipeline can
no longer be updated and locked dispatches/inputs/schedules become
valid.

```bash
dr pipelines lock <pipeline-id> [flags]
```

**Arguments:**

- `<pipeline-id>` — the UUID of the pipeline to lock.

**Flags:**

- `--output <json>` — emit machine-parseable JSON.

**Example:**

```bash
$ dr pipelines lock 6658f441-a8f5-4f21-b4d8-6cccf4c94c5b
Pipeline ID:  6658f441-a8f5-4f21-b4d8-6cccf4c94c5b
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
dr pipelines version list --pipeline <id> [--offset N] [--limit N] [--output json]
dr pipelines version get  --pipeline <id> <version-id>     [--output json]
```

`version list` returns the same data that's shown inline by
`pipelines get`, but in a paginated stand-alone view. `version get`
shows a single version in detail (pipeline, tasks, Python version,
creation timestamp, and any error detail).

### `graph`

Display the pipeline/task DAG for a pipeline as either a JSON
payload (for visualisation tooling) or a human-readable summary.

```bash
dr pipelines graph --pipeline <id>                       # draft graph
dr pipelines graph --pipeline <id> --version=N           # locked-version graph
dr pipelines graph --pipeline <id> --scope=draft         # explicit draft
dr pipelines graph --pipeline <id> --scope=locked --version=N
dr pipelines graph --pipeline <id> --output json         # JSON payload
```

The human view prints the pipeline header followed by `Nodes (N):` and
`Edges (M):` tables. Pass `--output json` to get the structured `Graph`
object (`lattice`, `nodes[]`, `edges[]` — JSON keys preserved while the API
wire format is unchanged) suitable for piping to
visualisation tooling.

See [Shared `--scope` / `--version` semantics](#scope--version-flags)
below for the full flag truth table.

### `input`

Manage JSON payloads that drive a dispatch.

```bash
dr pipelines input create --pipeline <id> <payload-file>            # draft scope
dr pipelines input create --pipeline <id> --version=N <payload-file>  # locked scope
dr pipelines input list   --pipeline <id> [--scope|--version] [--offset N] [--limit N]
dr pipelines input get    --pipeline <id> <input-id>      [--scope|--version]
dr pipelines input update --pipeline <id> <input-id> <payload-file>   # draft only
dr pipelines input delete --pipeline <id> <input-id>      [--scope|--version]
```

- The payload file must contain a JSON object. The CLI wraps it in
  `{"payload": …}` before sending.
- All verbs accept `--scope` / `--version` (see below). Inputs in the
  `locked` scope are immutable, so `input update` is draft-only.

### `dispatch`

Trigger, inspect, and cancel pipeline executions.

```bash
dr pipelines dispatch create --pipeline <id> --input <input-id>          # draft
dr pipelines dispatch create --pipeline <id> --version=N --input <input-id>  # locked
dr pipelines dispatch list   --pipeline <id> [--scope|--version]
dr pipelines dispatch get    --pipeline <id> <dispatch-id> [--scope|--version]
dr pipelines dispatch status --pipeline <id> <dispatch-id> [--scope|--version]
dr pipelines dispatch cancel --pipeline <id> <dispatch-id> [--scope|--version]
```

`dispatch status` is a lighter-weight call than `dispatch get` —
intended for polling — and returns just the dispatch ID, status, and
the corresponding Covalent dispatch ID.

`dispatch cancel` returns `409 Conflict` if the dispatch is already in
a terminal state (COMPLETED / FAILED / CANCELLED).

### `schedule`

Manage recurring (cron) dispatches on locked versions only. Both
`--pipeline` and `--version` are required for every verb.

```bash
dr pipelines schedule create --pipeline <id> --version=N \
    --cron "0 * * * *" --input <input-id> [--timezone UTC]
dr pipelines schedule list   --pipeline <id> --version=N [--offset N] [--limit N]
dr pipelines schedule get    --pipeline <id> --version=N <schedule-id>
dr pipelines schedule update --pipeline <id> --version=N <schedule-id> [--cron "*/15 * * * *"] [--timezone Europe/Berlin]
dr pipelines schedule delete --pipeline <id> --version=N <schedule-id>
```

`schedule update` requires at least one of `--cron` or `--timezone`.

### `environment`

Manage pipeline execution environments — named, immutable-versioned
bags of pip packages that pipelines can be built against. Environments
live at the top of the pipelines namespace (not nested under a specific
pipeline) and have their own lifecycle.

```bash
dr pipelines environment create --name <name> --package <spec> [--package <spec>] ...
dr pipelines environment list   [--offset N] [--limit N] [--output json]
dr pipelines environment update <environment-id> --package <spec> [...]
dr pipelines environment delete <environment-id>
dr pipelines environment version delete --environment <id> <version>
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
dr pipelines environment create --name ml-base \
    --package numpy --package pandas==2.0
dr pipelines environment create --name ml-base \
    --package "numpy,pandas==2.0,scikit-learn"
```

> [!NOTE]
> The pipelines-api currently does not surface `GET` endpoints for a
> single environment or for the version list. The full version history
> is only returned in the `create` and `update` responses.

## Shared flags

### `--scope` / `--version` flags

Inputs, dispatches, and `graph` mirror the API's two URL shapes —
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

`pipelines create`, `pipelines update`, `pipelines input create`, and
`pipelines input update` all accept the input file in two equivalent
ways:

```bash
dr pipelines create ./my_pipeline.py
dr pipelines create --from-file=./my_pipeline.py
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
dr pipelines list --output json | jq '.items[].pipeline_id'
```

### Global options

All [global flags](README.md#global-flags) are available, notably
`--debug` for protocol-level tracing and `--skip-auth` for advanced
scenarios.

## Local development

While iterating against a locally running pipelines-api (default port
`8000`), point the CLI at `http://localhost:8000` and bypass token
verification using the prefixed environment variables:

```bash
export DATAROBOT_CLI_FEATURE_PIPELINES=true
export DATAROBOT_CLI_ENDPOINT=http://localhost:8000/api/v2
export DATAROBOT_CLI_TOKEN=local
export DATAROBOT_CLI_SKIP_AUTH=true

./dist/dr pipelines list
```

Why each variable matters:

| Variable                          | Purpose                                                                                                  |
|-----------------------------------|----------------------------------------------------------------------------------------------------------|
| `DATAROBOT_CLI_FEATURE_PIPELINES` | Reveals the feature-gated `pipelines` command group.                                                     |
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

For a step-by-step walkthrough of how `dr pipelines list` was wired up,
see [Adding a command](../development/adding-a-command.md).

## Examples

### Pipeline lifecycle

```bash
# Register a draft, append a version, lock it, then delete it
dr pipelines create ./my_pipeline.py --description "Initial draft"
dr pipelines update <pipeline-id> ./my_pipeline.py
dr pipelines lock   <pipeline-id>
dr pipelines delete <pipeline-id>
```

### Inspect versions and graph

```bash
dr pipelines version list --pipeline <pipeline-id>
dr pipelines version get  --pipeline <pipeline-id> 2
dr pipelines graph        --pipeline <pipeline-id> --version=2 --output json
```

### Run a dispatch

```bash
# 1. Register a JSON input on the draft scope
dr pipelines input create --pipeline <pipeline-id> ./input.json

# 2. Trigger a dispatch with that input
dr pipelines dispatch create --pipeline <pipeline-id> --input <input-id>

# 3. Poll status until it reaches a terminal state
dr pipelines dispatch status --pipeline <pipeline-id> <dispatch-id>
```

### Schedule a recurring run on a locked version

```bash
dr pipelines schedule create \
    --pipeline <pipeline-id> --version=2 \
    --cron "0 */6 * * *" --input <input-id> --timezone America/Los_Angeles
```

### Scripting friendliness

```bash
# All UUIDs of locked pipelines, one per line
dr pipelines list --mode locked --output json | jq -r '.items[].pipeline_id'
```

## Error handling

The CLI surfaces backend errors verbatim and exits non-zero. The most
common status codes you will see:

| Status | Cause                                                                              |
|--------|------------------------------------------------------------------------------------|
| `400`  | Invalid Python file, mismatched pipeline name, or malformed JSON payload.          |
| `404`  | The provided `<pipeline-id>` / version / input / dispatch / schedule does not exist. |
| `409`  | Tried to update a `locked` pipeline, or cancel an already-terminal dispatch.       |

For most `get` / `delete` verbs the CLI translates a 404 into a
friendly informational line (e.g. `No pipeline found with id: …`) and
exits 0, so a no-op delete won't dump usage at the user.

## See also

- [Pipelines reference](pipelines-reference.md) — exhaustive table
  mapping every CLI command to its API endpoint, usage variants, and
  inputs.
- [Authentication](auth.md) — how `dr auth login` and `--skip-auth`
  interact.
- [Configuration](../user-guide/configuration.md) — config file and
  environment-variable precedence.
- [Adding a command](../development/adding-a-command.md) — how the
  pipelines verbs were built.
- [Feature gates](../development/feature-gates.md) — flipping
  `DATAROBOT_CLI_FEATURE_PIPELINES` on and off.

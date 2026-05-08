<!-- TOOD: remove this file, just to keep track internally while development is happening -->
# `dr pipelines` command reference

Complete cross-reference of every `dr pipelines …` subcommand, the
`pipelines-api` endpoint each one calls, sample invocations, and the
inputs (positional args, flags, request body fields) each command
accepts.

> All commands below assume the `pipelines` feature is enabled
> (`DATAROBOT_CLI_FEATURE_PIPELINES=true`).

## How to read this document

- **Method + path** is relative to `/api/v2`. The CLI prefixes the host
  from `DATAROBOT_CLI_ENDPOINT` (or `DATAROBOT_ENDPOINT`).
- **Usage** lists the canonical invocation plus common variants.
- **Inputs** names every positional argument and flag the command
  accepts. Flags shared by many commands (`--output`, `--scope`,
  `--version`, `--from-file`, `--skip-auth`) are described once at the
  bottom under "Shared flag semantics".

---

## Pipeline lifecycle

| Command | API endpoint | Usage | Inputs |
|---|---|---|---|
| `dr pipelines create` | `POST /pipelines` | `dr pipelines create ./my_pipeline.py` <br> `dr pipelines create --from-file=./my_pipeline.py` <br> `dr pipelines create ./my_pipeline.py --description "First draft" --mode draft` <br> `dr pipelines create --from-file=./my_pipeline.py --output json` | **Positional:** `<file>` (Python file defining a DataRobot pipeline; mutually exclusive with `--from-file`). <br> **Flags:** `--from-file=<path>`, `--description <text>`, `--mode draft\|locked`, `--output json`. |
| `dr pipelines list` | `GET /pipelines` | `dr pipelines list` <br> `dr pipelines list --mode draft` <br> `dr pipelines list --offset 50 --limit 10 --output json` | **Flags:** `--mode draft\|locked`, `--offset <n>`, `--limit <n>`, `--output json`. |
| `dr pipelines get` | `GET /pipelines/{pipeline_id}` | `dr pipelines get <pipeline-id>` <br> `dr pipelines get <pipeline-id> --output json` | **Positional:** `<pipeline-id>` (required). <br> **Flags:** `--output json`. |
| `dr pipelines update` | `PATCH /pipelines/{pipeline_id}` | `dr pipelines update <pipeline-id> ./my_pipeline.py` <br> `dr pipelines update <pipeline-id> --from-file=./my_pipeline.py` <br> `dr pipelines update <pipeline-id> --from-file=./my_pipeline.py --output json` | **Positional:** `<pipeline-id>` (required), `<file>` (mutually exclusive with `--from-file`). <br> **Flags:** `--from-file=<path>`, `--output json`. |
| `dr pipelines delete` | `DELETE /pipelines/{pipeline_id}` | `dr pipelines delete <pipeline-id>` | **Positional:** `<pipeline-id>` (required). |
| `dr pipelines lock` | `PATCH /pipelines/{pipeline_id}/mode` | `dr pipelines lock <pipeline-id>` <br> `dr pipelines lock <pipeline-id> --output json` | **Positional:** `<pipeline-id>` (required). <br> **Flags:** `--output json`. <br> **Body:** none (the API uses absence-of-body as the promote signal). |

---

## Versions

| Command | API endpoint | Usage | Inputs |
|---|---|---|---|
| `dr pipelines version list` | `GET /pipelines/{pipeline_id}/versions` | `dr pipelines version list --pipeline <id>` <br> `dr pipelines version list --pipeline <id> --offset 10 --limit 5 --output json` | **Flags:** `--pipeline <id>` (required), `--offset <n>`, `--limit <n>`, `--output json`. |
| `dr pipelines version get` | `GET /pipelines/{pipeline_id}/versions/{version_id}` | `dr pipelines version get --pipeline <id> 2` <br> `dr pipelines version get --pipeline <id> 2 --output json` | **Positional:** `<version-id>` (positive integer, required). <br> **Flags:** `--pipeline <id>` (required), `--output json`. |
| `dr pipelines graph` | `GET /pipelines/{pipeline_id}/graph` (draft) <br> `GET /pipelines/{pipeline_id}/versions/{version_id}/graph` (locked) | `dr pipelines graph --pipeline <id>` (draft) <br> `dr pipelines graph --pipeline <id> --version=2` (locked) <br> `dr pipelines graph --pipeline <id> --version=2 --output json` | **Flags:** `--pipeline <id>` (required), `--scope draft\|locked`, `--version <n>`, `--output json`. |

---

## Inputs (`dr pipelines input …`)

Inputs come in two scopes — **draft** (mutable, no version pinned) and
**locked** (immutable, tied to a frozen version). Scope selection rules
are documented under "Shared flag semantics" below.

| Command | API endpoint | Usage | Inputs |
|---|---|---|---|
| `dr pipelines input create` | `POST /pipelines/{id}/inputs` (draft) <br> `POST /pipelines/{id}/versions/{ver}/inputs` (locked) | `dr pipelines input create --pipeline <id> ./payload.json` <br> `dr pipelines input create --pipeline <id> --from-file=./payload.json` <br> `dr pipelines input create --pipeline <id> --version=2 ./payload.json --output json` | **Positional:** `<payload-file>` (JSON object; mutually exclusive with `--from-file`). <br> **Flags:** `--pipeline <id>` (required), `--scope`, `--version`, `--from-file=<path>`, `--output json`. <br> **Body sent to API:** `{"payload": <object from file>}`. |
| `dr pipelines input list` | `GET /pipelines/{id}/inputs` (draft) <br> `GET /pipelines/{id}/versions/{ver}/inputs` (locked) | `dr pipelines input list --pipeline <id>` (draft) <br> `dr pipelines input list --pipeline <id> --version=2` (locked) <br> `dr pipelines input list --pipeline <id> --offset 10 --limit 5 --output json` | **Flags:** `--pipeline <id>` (required), `--scope`, `--version`, `--offset <n>`, `--limit <n>`, `--output json`. |
| `dr pipelines input get` | `GET /pipelines/{id}/inputs/{input_id}` (draft) <br> `GET /pipelines/{id}/versions/{ver}/inputs/{input_id}` (locked) | `dr pipelines input get --pipeline <id> <input-id>` <br> `dr pipelines input get --pipeline <id> --version=2 <input-id> --output json` | **Positional:** `<input-id>` (required). <br> **Flags:** `--pipeline <id>` (required), `--scope`, `--version`, `--output json`. |
| `dr pipelines input update` | `PATCH /pipelines/{id}/inputs/{input_id}` (draft only) | `dr pipelines input update --pipeline <id> <input-id> ./payload.json` <br> `dr pipelines input update --pipeline <id> <input-id> --from-file=./payload.json --output json` | **Positional:** `<input-id>` (required), `<payload-file>` (JSON object; mutually exclusive with `--from-file`). <br> **Flags:** `--pipeline <id>` (required), `--from-file=<path>`, `--output json`. <br> **Body sent to API:** `{"payload": <object from file>}`. |
| `dr pipelines input delete` | `DELETE /pipelines/{id}/inputs/{input_id}` (draft) <br> `DELETE /pipelines/{id}/versions/{ver}/inputs/{input_id}` (locked) | `dr pipelines input delete --pipeline <id> <input-id>` <br> `dr pipelines input delete --pipeline <id> --version=2 <input-id>` | **Positional:** `<input-id>` (required). <br> **Flags:** `--pipeline <id>` (required), `--scope`, `--version`. |

---

## Runs (`dr pipelines run …`)

Same draft/locked scope rules as inputs. The wire-level URLs still use
the legacy term `dispatches` / `dispatch_id`, but the CLI's `--output
json` remaps these to `run_id` / `covalent_run_id` so the JSON output
matches the rest of the CLI vocabulary.

| Command | API endpoint | Usage | Inputs |
|---|---|---|---|
| `dr pipelines run create` | `POST /pipelines/{id}/dispatches` (draft) <br> `POST /pipelines/{id}/versions/{ver}/dispatches` (locked) | `dr pipelines run create --pipeline <id> --input <input-id>` <br> `dr pipelines run create --pipeline <id> --version=2 --input <input-id> --output json` | **Flags:** `--pipeline <id>` (required), `--input <input-id>` (required), `--scope`, `--version`, `--output json`. <br> **Body sent to API:** `{"input_id": "<input-id>"}`. |
| `dr pipelines run list` | `GET /pipelines/{id}/dispatches` (draft) <br> `GET /pipelines/{id}/versions/{ver}/dispatches` (locked) | `dr pipelines run list --pipeline <id>` <br> `dr pipelines run list --pipeline <id> --version=2 --output json` | **Flags:** `--pipeline <id>` (required), `--scope`, `--version`, `--offset <n>`, `--limit <n>`, `--output json`. |
| `dr pipelines run get` | `GET /pipelines/{id}/dispatches/{dispatch_id}` (draft) <br> `GET /pipelines/{id}/versions/{ver}/dispatches/{dispatch_id}` (locked) | `dr pipelines run get --pipeline <id> <run-id>` <br> `dr pipelines run get --pipeline <id> --version=2 <run-id> --output json` | **Positional:** `<run-id>` (required). <br> **Flags:** `--pipeline <id>` (required), `--scope`, `--version`, `--output json`. |
| `dr pipelines run status` | `GET /pipelines/{id}/dispatches/{dispatch_id}/status` (draft) <br> `GET /pipelines/{id}/versions/{ver}/dispatches/{dispatch_id}/status` (locked) | `dr pipelines run status --pipeline <id> <run-id>` <br> `dr pipelines run status --pipeline <id> --version=2 <run-id> --output json` | **Positional:** `<run-id>` (required). <br> **Flags:** `--pipeline <id>` (required), `--scope`, `--version`, `--output json`. |
| `dr pipelines run cancel` | `DELETE /pipelines/{id}/dispatches/{dispatch_id}` (draft) <br> `DELETE /pipelines/{id}/versions/{ver}/dispatches/{dispatch_id}` (locked) | `dr pipelines run cancel --pipeline <id> <run-id>` <br> `dr pipelines run cancel --pipeline <id> --version=2 <run-id>` | **Positional:** `<run-id>` (required). <br> **Flags:** `--pipeline <id>` (required), `--scope`, `--version`. |

---

## Schedules (`dr pipelines schedule …`)

Schedules are **locked-only** — every verb requires both `--pipeline` and
`--version`. There is no draft scope or `--scope` flag.

| Command | API endpoint | Usage | Inputs |
|---|---|---|---|
| `dr pipelines schedule create` | `POST /pipelines/{id}/versions/{ver}/schedules` | `dr pipelines schedule create --pipeline <id> --version=2 --cron "0 * * * *" --input <input-id>` <br> `dr pipelines schedule create --pipeline <id> --version=2 --cron "0 9 * * *" --input <input-id> --timezone America/Los_Angeles` <br> `… --output json` | **Flags:** `--pipeline <id>` (required), `--version <n>` (required, > 0), `--cron "<expr>"` (required), `--input <input-id>` (required), `--timezone <iana>` (default `UTC`), `--output json`. <br> **Body sent to API:** `{"cron_expression": "...", "pipeline_input_id": "...", "timezone": "..."}`. |
| `dr pipelines schedule list` | `GET /pipelines/{id}/versions/{ver}/schedules` | `dr pipelines schedule list --pipeline <id> --version=2` <br> `dr pipelines schedule list --pipeline <id> --version=2 --offset 10 --limit 5 --output json` | **Flags:** `--pipeline <id>` (required), `--version <n>` (required, > 0), `--offset <n>`, `--limit <n>`, `--output json`. |
| `dr pipelines schedule get` | `GET /pipelines/{id}/versions/{ver}/schedules/{schedule_id}` | `dr pipelines schedule get --pipeline <id> --version=2 <schedule-id>` <br> `… --output json` | **Positional:** `<schedule-id>` (required). <br> **Flags:** `--pipeline <id>` (required), `--version <n>` (required, > 0), `--output json`. |
| `dr pipelines schedule update` | `PATCH /pipelines/{id}/versions/{ver}/schedules/{schedule_id}` | `dr pipelines schedule update --pipeline <id> --version=2 <schedule-id> --cron "*/15 * * * *"` <br> `dr pipelines schedule update --pipeline <id> --version=2 <schedule-id> --timezone Europe/Berlin` <br> `… --cron "0 0 * * *" --timezone UTC --output json` | **Positional:** `<schedule-id>` (required). <br> **Flags:** `--pipeline <id>` (required), `--version <n>` (required, > 0), `--cron "<expr>"`, `--timezone <iana>`, `--output json`. At least one of `--cron`/`--timezone` must be supplied. <br> **Body sent to API:** `{"cron_expression"?: "...", "timezone"?: "..."}` (only fields you changed). |
| `dr pipelines schedule delete` | `DELETE /pipelines/{id}/versions/{ver}/schedules/{schedule_id}` | `dr pipelines schedule delete --pipeline <id> --version=2 <schedule-id>` | **Positional:** `<schedule-id>` (required). <br> **Flags:** `--pipeline <id>` (required), `--version <n>` (required, > 0). |

---

## Execution environments (`dr pipelines environment …`)

Pipeline execution environments are named, immutable-versioned bags of
pip packages. They live at the top of the pipelines namespace (not
nested under a specific pipeline) and are created/updated independently.
Each `update` appends a new version; older versions can be deleted
individually.

| Command | API endpoint | Usage | Inputs |
|---|---|---|---|
| `dr pipelines environment create` | `POST /pipelines/environments` | `dr pipelines environment create --name ml-base --package numpy --package pandas==2.0` <br> `dr pipelines environment create --name ml-base --package "numpy,pandas==2.0" --description "training base" --output json` | **Flags:** `--name <name>` (required), `--description <text>`, `--package <spec>` (required, repeatable, also accepts comma-separated values), `--output json`. <br> **Body sent to API:** `{"name": "...", "description"?: "...", "packages": ["..."]}`. |
| `dr pipelines environment list` | `GET /pipelines/environments` | `dr pipelines environment list` <br> `dr pipelines environment list --offset 50 --limit 10 --output json` | **Flags:** `--offset <n>`, `--limit <n>`, `--output json`. |
| `dr pipelines environment update` | `PATCH /pipelines/environments/{environment_id}` | `dr pipelines environment update <environment-id> --package scikit-learn` <br> `dr pipelines environment update <environment-id> --package "scikit-learn,torch" --output json` | **Positional:** `<environment-id>` (required). <br> **Flags:** `--package <spec>` (required, repeatable, also accepts comma-separated values), `--output json`. <br> **Body sent to API:** `{"packages": ["..."]}`. |
| `dr pipelines environment delete` | `DELETE /pipelines/environments/{environment_id}` | `dr pipelines environment delete <environment-id>` | **Positional:** `<environment-id>` (required). |
| `dr pipelines environment version delete` | `DELETE /pipelines/environments/{environment_id}/versions/{version_id}` | `dr pipelines environment version delete --environment <id> 2` | **Positional:** `<version>` (positive integer, required). <br> **Flags:** `--environment <id>` (required). |

> [!NOTE]
> The pipelines-api currently does not expose `GET` endpoints for a
> single environment or for individual versions, so the CLI does not
> ship `environment get` or `environment version get`. The full version
> history is only returned in the `create` and `update` responses.

---

## Shared flag semantics

### `--scope` / `--version` (inputs, runs, graph)

The CLI mirrors the API's two URL shapes — `/pipelines/{id}/…` for the
mutable draft and `/pipelines/{id}/versions/{ver}/…` for a locked
version — through a pair of optional flags:

| Flags supplied | Resolved scope | URL used |
|---|---|---|
| _(none)_ | `draft` | `/pipelines/{id}/…` |
| `--version=N` | `locked` (auto) | `/pipelines/{id}/versions/N/…` |
| `--scope=draft` | `draft` | `/pipelines/{id}/…` |
| `--scope=locked --version=N` | `locked` | `/pipelines/{id}/versions/N/…` |
| `--scope=draft --version=N` | **error** | `--scope=draft cannot be combined with --version` |
| `--scope=locked` (no `--version`) | **error** | `--scope=locked requires --version=<n>` |
| `--scope=garbage` | **error** | `invalid --scope: "garbage" (supported: draft, locked)` |

### `--from-file` / positional file (create + update verbs)

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

Every read/write verb that produces a payload accepts `--output json` to
emit the underlying response struct as indented JSON. Any other value
(e.g. `--output yaml`, `--output csv`) is rejected with
`invalid output format: <value> (supported: json)`.

### `auth` / `--skip-auth`

All verbs run `auth.EnsureAuthenticatedE` as their `PreRunE`. Pass the
global `--skip-auth` flag (or set `DATAROBOT_CLI_SKIP_AUTH=true`) when
exercising a local API stub that doesn't implement `/version/`.

---

## Quick endpoint lookup

| API endpoint | CLI command |
|---|---|
| `POST /pipelines` | `dr pipelines create` |
| `GET /pipelines` | `dr pipelines list` |
| `GET /pipelines/{id}` | `dr pipelines get` |
| `PATCH /pipelines/{id}` | `dr pipelines update` |
| `DELETE /pipelines/{id}` | `dr pipelines delete` |
| `PATCH /pipelines/{id}/mode` | `dr pipelines lock` |
| `GET /pipelines/{id}/versions` | `dr pipelines version list` |
| `GET /pipelines/{id}/versions/{ver}` | `dr pipelines version get` |
| `GET /pipelines/{id}/graph` | `dr pipelines graph` (draft) |
| `GET /pipelines/{id}/versions/{ver}/graph` | `dr pipelines graph --version=N` |
| `POST /pipelines/{id}/inputs` | `dr pipelines input create` (draft) |
| `POST /pipelines/{id}/versions/{ver}/inputs` | `dr pipelines input create --version=N` |
| `GET /pipelines/{id}/inputs` | `dr pipelines input list` (draft) |
| `GET /pipelines/{id}/versions/{ver}/inputs` | `dr pipelines input list --version=N` |
| `GET /pipelines/{id}/inputs/{input_id}` | `dr pipelines input get` (draft) |
| `GET /pipelines/{id}/versions/{ver}/inputs/{input_id}` | `dr pipelines input get --version=N` |
| `PATCH /pipelines/{id}/inputs/{input_id}` | `dr pipelines input update` |
| `DELETE /pipelines/{id}/inputs/{input_id}` | `dr pipelines input delete` (draft) |
| `DELETE /pipelines/{id}/versions/{ver}/inputs/{input_id}` | `dr pipelines input delete --version=N` |
| `POST /pipelines/{id}/dispatches` | `dr pipelines run create` (draft) |
| `POST /pipelines/{id}/versions/{ver}/dispatches` | `dr pipelines run create --version=N` |
| `GET /pipelines/{id}/dispatches` | `dr pipelines run list` (draft) |
| `GET /pipelines/{id}/versions/{ver}/dispatches` | `dr pipelines run list --version=N` |
| `GET /pipelines/{id}/dispatches/{dispatch_id}` | `dr pipelines run get` (draft) |
| `GET /pipelines/{id}/versions/{ver}/dispatches/{dispatch_id}` | `dr pipelines run get --version=N` |
| `GET /pipelines/{id}/dispatches/{dispatch_id}/status` | `dr pipelines run status` (draft) |
| `GET /pipelines/{id}/versions/{ver}/dispatches/{dispatch_id}/status` | `dr pipelines run status --version=N` |
| `DELETE /pipelines/{id}/dispatches/{dispatch_id}` | `dr pipelines run cancel` (draft) |
| `DELETE /pipelines/{id}/versions/{ver}/dispatches/{dispatch_id}` | `dr pipelines run cancel --version=N` |
| `POST /pipelines/{id}/versions/{ver}/schedules` | `dr pipelines schedule create` |
| `GET /pipelines/{id}/versions/{ver}/schedules` | `dr pipelines schedule list` |
| `GET /pipelines/{id}/versions/{ver}/schedules/{schedule_id}` | `dr pipelines schedule get` |
| `PATCH /pipelines/{id}/versions/{ver}/schedules/{schedule_id}` | `dr pipelines schedule update` |
| `DELETE /pipelines/{id}/versions/{ver}/schedules/{schedule_id}` | `dr pipelines schedule delete` |
| `POST /pipelines/environments` | `dr pipelines environment create` |
| `GET /pipelines/environments` | `dr pipelines environment list` |
| `PATCH /pipelines/environments/{environment_id}` | `dr pipelines environment update` |
| `DELETE /pipelines/environments/{environment_id}` | `dr pipelines environment delete` |
| `DELETE /pipelines/environments/{environment_id}/versions/{version_id}` | `dr pipelines environment version delete` |

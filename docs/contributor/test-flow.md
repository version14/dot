# test-flow: End-to-End Flow Testing

`test-flow` is a command-line tool that runs the full DOT scaffold pipeline against scripted JSON fixtures — no terminal interaction required. It is the primary way to verify that a flow, its generators, and their commands all work together correctly.

---

## Table of Contents

- [When to use test-flow](#when-to-use-test-flow)
- [How it works](#how-it-works)
- [Running test-flow](#running-test-flow)
- [Flags](#flags)
- [Writing a fixture](#writing-a-fixture)
- [Fixture schema reference](#fixture-schema-reference)
- [Loop fixtures](#loop-fixtures)
- [Plugin injection fixtures](#plugin-injection-fixtures)
- [expected_visited](#expected_visited)
- [Exit codes and output](#exit-codes-and-output)
- [Existing fixtures](#existing-fixtures)
- [When to add a new fixture](#when-to-add-a-new-fixture)

---

## When to use test-flow

Use `test-flow` whenever you:

- Add or modify a flow.
- Add or modify a generator.
- Add or modify a plugin injection.
- Change a `PostGenerationCommand` or `TestCommand` in a manifest.
- Want to verify that the full pipeline (flow → generators → validators → commands) produces a working project.

`test-flow` is not a replacement for unit tests. Unit tests cover individual functions in isolation; `test-flow` exercises the entire scaffold pipeline end-to-end.

---

## How it works

For each fixture, `test-flow`:

1. Looks up the flow by `flow_id` in the flows registry.
2. Replays the fixture's `answers` through a `scriptedAdapter` (no TUI).
3. Runs the full generator pipeline into a fresh temp directory.
4. Runs `Validators` against the generated files.
5. Runs `PostGenerationCommands` (unless `skip_post_commands: true`).
6. Runs `TestCommands` including background dev servers (unless `skip_test_commands: true`).
7. Reports pass / fail per step with timing.

Plugin injections are active. The `scriptedAdapter` must provide answers for any questions injected by active plugins (e.g. `biome_extras.strict_mode`).

---

## Running test-flow

```bash
# Run all fixtures
go run ./tools/test-flow

# Run a specific fixture by name
go run ./tools/test-flow -only turborepo_ts_react

# Run multiple fixtures
go run ./tools/test-flow -only "turborepo_ts_react,single_go"

# Skip post-gen commands globally (faster, offline)
go run ./tools/test-flow -skip-post

# Skip test commands globally (much faster)
go run ./tools/test-flow -skip-test

# Keep scratch directories for inspection after run
go run ./tools/test-flow -keep

# Use a custom testdata directory
go run ./tools/test-flow -dir ./my-testdata

# Use a custom temp directory for scratch dirs
go run ./tools/test-flow -tmp /tmp/dot-test-runs
```

The Makefile shortcut:

```bash
make test-flow            # equivalent to go run ./tools/test-flow -skip-test
```

---

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-dir DIR` | `tools/test-flow/testdata` | Directory containing `*.json` fixture files. |
| `-tmp DIR` | system temp | Parent directory for per-case scratch directories. |
| `-skip-post` | `false` | Skip all `PostGenerationCommands` globally. Overrides the fixture's `skip_post_commands`. |
| `-skip-test` | `false` | Skip all `TestCommands` globally. Overrides the fixture's `skip_test_commands`. |
| `-only NAMES` | (all) | Comma-separated list of fixture `name` values to run. |
| `-keep` | `false` | Do not delete scratch directories after the run. Lets you inspect generated files. |

---

## Writing a fixture

Create a `.json` file in `tools/test-flow/testdata/`. The filename is used as the fixture name if `name` is not set.

Minimum fixture:

```json
{
  "name": "my_flow_basic",
  "flow_id": "my-flow",
  "answers": {
    "project_name": "test-project",
    "confirm_generate": true
  }
}
```

The fixture must provide answers for **every question the engine visits** — including questions injected by active plugins. If an answer is missing, the scripted adapter returns an error and the case fails with:

```
✗ scaffold  : test-flow: no scripted answer for question "plugin_id.question_name"
```

### Finding which questions are visited

Run with `expected_visited` empty first, then inspect the output:

```
✗ verify visited  : mismatch
      expected: []
      actual:   [project_name monorepo_type stack use_react use_biome biome_extras.strict_mode confirm_generate]
```

Copy the `actual` list into `expected_visited`.

---

## Fixture schema reference

```json
{
  "name": "fixture_name",
  "flow_id": "flow-id",
  "answers": {
    "question_id": "answer_value"
  },
  "expected_visited": ["question_id", "..."],
  "skip_post_commands": false,
  "skip_test_commands": false
}
```

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `name` | string | No | Identifier for `-only` and reports. Defaults to filename. |
| `flow_id` | string | Yes | Must match a registered flow ID. |
| `answers` | object | Yes | Map of question ID → answer. |
| `expected_visited` | string[] | No | If set, the engine must visit exactly these IDs in this order. |
| `skip_post_commands` | bool | No | Skip `PostGenerationCommands` for this fixture only. |
| `skip_test_commands` | bool | No | Skip `TestCommands` for this fixture only. |

### Answer types

| Question type | JSON type | Example |
|--------------|-----------|---------|
| `TextQuestion` | string | `"my-project"` |
| `ConfirmQuestion` | boolean | `true` or `false` |
| `OptionQuestion` (single) | string | `"turborepo"` |
| `OptionQuestion` (multi) | string[] | `["eslint", "prettier"]` |
| `LoopQuestion` | array of objects | See [Loop fixtures](#loop-fixtures) |

---

## Loop fixtures

Loop questions expect an array of objects — one per iteration:

```json
{
  "name": "microservices_three",
  "flow_id": "microservices",
  "answers": {
    "project_name": "platform",
    "services": [
      {"service_name": "auth",    "service_port": "3001"},
      {"service_name": "users",   "service_port": "3002"},
      {"service_name": "billing", "service_port": "3003"}
    ],
    "confirm_generate": true
  }
}
```

Each object in the array provides answers for one iteration of the loop body. The scripted adapter runs the loop body once per element, collecting answers from the corresponding object, then falls back to the top-level `answers` for any key not found in the iteration object.

The number of iterations equals the number of objects in the array.

---

## Plugin injection fixtures

If active plugins inject questions into the flow, the fixture must include answers for those questions. The injected question IDs are prefixed with the plugin's ID:

```json
{
  "name": "turborepo_ts_react",
  "flow_id": "init",
  "answers": {
    "project_name": "demo-app",
    "monorepo_type": "turborepo",
    "stack": "typescript",
    "use_react": true,
    "use_biome": true,
    "biome_extras.strict_mode": false,
    "confirm_generate": true
  }
}
```

`biome_extras.strict_mode` is injected by the `biome_extras` plugin (an `InsertAfter` on `use_biome`). The fixture provides it because the plugin is active when `test-flow` runs.

---

## expected_visited

`expected_visited` is an optional ordered list of question IDs the engine must visit. If the actual visited sequence does not match exactly, the case fails:

```
✗ verify visited
      expected: [project_name monorepo_type stack confirm_generate]
      actual:   [project_name monorepo_type stack use_react confirm_generate]
```

Use it to:

- Verify that branching logic is correct (certain branches are taken / skipped).
- Catch regressions when a flow's question graph changes.
- Document the intended happy-path question sequence.

When `expected_visited` is empty, the check is skipped.

---

## Exit codes and output

`test-flow` prints a structured report per case:

```
[1/3] turborepo_ts_react (flow=monorepo)
  ✓ flow                        — 7 nodes visited
  ✓ verify visited              — matches expected
  ✓ resolved generators         — base_project, typescript_base, react_app, biome_config
  ✓ scaffolded files            — → /tmp/dot-test-monorepo-xyz/demo-app
  ✓ validators                  — 12 passed
  → post-gen commands (1)
    ✓  pnpm install              [8.4s]
  → test commands (3)
    ✓  pnpm exec tsc --noEmit    [2.1s]
    ✓  pnpm exec vite build      [6.8s]
    ✓  pnpm exec vite            [background, ready+stop 4.0s]
  PASS

✓ All 3 cases passed
```

| Exit code | Meaning |
|-----------|---------|
| `0` | All cases passed |
| `1` | One or more cases failed |
| `2` | Usage / I/O error (no fixtures found, unknown flag) |

---

## Existing fixtures

| Fixture | Flow | What it covers |
|---------|------|----------------|
| `single_go.json` | `monorepo` | Single-package Go project, no React, no Biome |
| `turborepo_ts_react.json` | `monorepo` | Turborepo + TypeScript + React + Biome (non-strict) |
| `biome_extras_strict.json` | `monorepo` | Biome strict mode plugin injection |
| `fullstack_react.json` | `fullstack` | Full-stack with React + Biome |
| `fullstack_no_ui.json` | `fullstack` | Full-stack without UI |
| `microservices_three.json` | `microservices` | 3-service loop — auth, users, billing |
| `plugin_template_full.json` | `plugin-template` | Full plugin scaffold with injection + generator |

---

## When to add a new fixture

| Scenario | Action |
|----------|--------|
| New flow added | Add at least one fixture covering the happy path |
| New branch in an existing flow | Add or extend a fixture that exercises the branch |
| New plugin injection | Add or update a fixture that sets the injected question's answer |
| Loop body question added | Update loop fixtures to include the new key in each iteration object |
| `PostGenerationCommand` added | Update the fixture's `skip_post_commands` to `false` and verify it passes |
| `TestCommand` added | Update `skip_test_commands` to `false` and add the command's expected output |

Every fixture is a contract. If the flow changes and the fixture breaks, it is a signal to update the fixture and document what changed.

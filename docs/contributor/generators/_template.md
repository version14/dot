<!-- ─────────────────────────────────────────────────────────────────────────
  GENERATOR DOC TEMPLATE
  ─────────────────────────────────────────────────────────────────────────
  How to use:
    1. Copy this file to docs/generators/<name>.md
    2. Replace every <!-- TODO --> comment and every _placeholder_ value
    3. Delete this header block before committing

  Naming: the filename must match Manifest.Name exactly.
  Add the new file to the table in docs/README.md (generators section).
  ───────────────────────────────────────────────────────────────────────── -->

# Generator: `_name_`

<!-- TODO: one sentence — what does this generator produce and when is it used? -->

---

## Identity

| Field | Value |
|-------|-------|
| Name | `_name_` |
| Version | `_0.1.0_` |
| Package | `generators/_name_` |

---

## Dependencies

<!-- TODO: list every generator in Manifest.DependsOn, with a "Why" column.
     If there are none, write: None. -->

| Generator | Why |
|-----------|-----|
| `_dep_` | _reason_ |

---

## Answers consumed

<!-- TODO: list every ctx.Answers key the generator reads, its expected type,
     whether it is required or optional, and any fallback behaviour.
     If the generator reads no answers, write: None. -->

| Key | Type | Required | Notes |
|-----|------|----------|-------|
| `_key_` | _string_ | _Yes / No_ | _description, default if optional_ |

---

## Files written

<!-- TODO: list every path passed to ctx.State.WriteFile / UpdateJSON / etc.
     For dynamic paths (loop generators), describe the pattern instead. -->

| Path | Description |
|------|-------------|
| `_path_` | _what it contains_ |

<!-- TODO: if the generator merges into files created by other generators,
     add a second table like this:

Also merges into:

| Path | Keys added / updated |
|------|---------------------|
| `_path_` | `_key_` — _description_ |
-->

---

## Validators

<!-- TODO: copy from Manifest.Validators. If there are none, write: None. -->

| Check | Type | Passes when |
|-------|------|-------------|
| `_path_` | `file_exists` | File is present after generation |
| `_key_` in `_path_` | `json_key_exists` | JSON key exists |

---

## Post-generation commands

<!-- TODO: copy from Manifest.PostGenerationCommands.
     If there are none, write: No PostGenerationCommands. -->

| Command | WorkDir | Notes |
|---------|---------|-------|
| `_cmd_` | `_workdir or "project root"_` | _notes_ |

## Test commands

<!-- TODO: copy from Manifest.TestCommands.
     If there are none, write: No TestCommands. -->

| Command | Background | Ready delay | Notes |
|---------|-----------|-------------|-------|
| `_cmd_` | No | — | _notes_ |

---

## Conflicts

<!-- TODO: list every generator in Manifest.ConflictsWith.
     If there are none, write: None. -->

None.

---

<!-- TODO: add any "see also" links if this generator is extended by a plugin
     or has a related flow doc. Delete this section if unused.

## See also

- [docs/plugins/_plugin_.md](../plugins/_plugin_.md)
- [docs/flows/_flow_.md](../flows/_flow_.md)
-->

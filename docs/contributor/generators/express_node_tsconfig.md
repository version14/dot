# Generator: `express_node_tsconfig`

Overrides tsconfig.json compiler options for a Node.js/CommonJS Express backend (module: CommonJS, moduleResolution: Node, rootDir/outDir). Also sets `"type": "commonjs"` in package.json to prevent ESM conflicts.

---

## Identity

| Field | Value |
|-------|-------|
| Name | `express_node_tsconfig` |
| Version | `0.1.0` |
| Package | `generators/express_node_tsconfig` |

---

## Dependencies

| Generator | Why |
|-----------|-----|
| `typescript_base` | `tsconfig.json` must exist before it can be overridden |

---

## Answers consumed

None.

---

## Files written

Also merges into:

| Path | Keys added / updated |
|------|---------------------|
| `tsconfig.json` | `compilerOptions.module = "CommonJS"`, `compilerOptions.moduleResolution = "Node"`, `compilerOptions.rootDir`, `compilerOptions.resolveJsonModule`, `include`, `exclude` |
| `package.json` | `type = "commonjs"` (overrides `typescript_base`'s `"module"`) |

---

## Validators

| Check | Type | Passes when |
|-------|------|-------------|
| `compilerOptions.module` in `tsconfig.json` | `json_key_exists` | — |

---

## Post-generation commands

No PostGenerationCommands.

## Test commands

No TestCommands.

---

## Conflicts

None.

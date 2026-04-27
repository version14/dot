# Getting Started with DOT

DOT is a generative project scaffolding CLI. You answer a few questions, DOT builds a production-ready project and records the choices so you can update or extend it later.

---

## Table of Contents

- [Install](#install)
- [Your first scaffold](#your-first-scaffold)
- [Explore what was built](#explore-what-was-built)
- [Update a project](#update-a-project)
- [Check project health](#check-project-health)
- [Manage plugins](#manage-plugins)
- [Uninstall](#uninstall)

---

## Install

### macOS / Linux — curl (no Go required)

```bash
curl -fsSL https://raw.githubusercontent.com/version14/dot/main/install.sh | sh
```

Installs to `/usr/local/bin/dot`. To choose a different location:

```bash
INSTALL_DIR=~/bin sh -c "$(curl -fsSL https://raw.githubusercontent.com/version14/dot/main/install.sh)"
```

### macOS — Homebrew

```bash
brew install version14/tap/dot
```

### go install

```bash
go install github.com/version14/dot/cmd/dot@latest
```

Requires Go 1.21+. The binary lands in `$GOPATH/bin` (usually already on `$PATH`).

### From source

```bash
git clone https://github.com/version14/dot.git
cd dot
make build        # produces bin/dot
export PATH="$PWD/bin:$PATH"
dot version
```

### Verify the installation

```bash
dot version
# dot 0.1.0
```

---

## Your first scaffold

### 1. List available flows

```bash
dot flows
```

Output:

```
Flows
  monorepo      Turborepo monorepo (apps + shared packages)
  fullstack     Full-stack app (Next.js + Go API)
  microservices Go microservices (multiple services, Docker Compose)
  plugin-template  Plugin Repository Template
```

### 2. Scaffold a project

```bash
dot scaffold
```

If there is more than one flow, DOT lists them and asks you to pick one:

```
Multiple flows available — re-run with one of:
  dot scaffold monorepo
  dot scaffold fullstack
  dot scaffold microservices
  dot scaffold plugin-template
```

Or pass the flow ID directly:

```bash
dot scaffold monorepo
```

### 3. Answer the interactive questions

DOT opens a terminal form. Use arrow keys to navigate, `Enter` to confirm, `Ctrl+C` to abort. Pressing `Escape` in any field takes you back to the previous question.

### 4. Wait for generators and post-gen commands

Once you confirm, DOT:

1. Runs the generator pipeline (outputs a line per generator).
2. Runs any `PostGenerationCommands` (e.g. `pnpm install`) with a live spinner.

```
✓ base_project
✓ typescript_base
✓ react_app

post-gen commands (1)
  ✓  pnpm install  [12.3s]

✓ scaffolded my-app in ./my-app
```

To skip post-generation commands (useful for CI or offline runs):

```bash
dot scaffold monorepo --skip-post
```

To write the project to a specific directory:

```bash
dot scaffold monorepo -out /tmp
# Creates /tmp/my-project/
```

---

## Explore what was built

After scaffolding, the project directory contains a `.dot/` folder:

```
my-project/
├── .dot/
│   ├── spec.json        ← your answers + the flow that was used
│   └── manifest.json    ← which generators ran, at what version, when
├── README.md
├── ...
```

`spec.json` is the machine-readable record of every decision. DOT reads it on `dot update` and `dot doctor`. You can inspect it:

```bash
cat my-project/.dot/spec.json
```

---

## Update a project

After a DOT upgrade or after changing a flow, re-run the generators against an existing project without re-answering the questions:

```bash
dot update ./my-project
```

DOT reads `.dot/spec.json`, replays the same answers through the current generators, and writes any changed files.

---

## Check project health

`dot doctor` compares the `.dot/spec.json` against the currently installed generators and flags drift:

```bash
dot doctor ./my-project
```

Output example:

```
✓ flow: monorepo
✓ generators: all found
⚠ version drift: base_project spec=0.1.0 installed=0.2.0
✓ validators: 12 passed
```

Exit code 0 when no issues are found, 1 when issues are detected.

---

## Manage plugins

Plugins add generators and inject new questions into existing flows. They are installed from git repositories.

### List plugins

```bash
dot plugin list
```

Shows built-in plugins (compiled in) and installed plugins:

```
Built-in plugins
  biome_extras              1 generators · 2 injections

Installed plugins (~/.dot/plugins)
  (none)
```

### Install a plugin from GitHub

```bash
dot plugin install github.com/version14/dot-plugin-example
```

Accepts any of these source forms:

| Form | Meaning |
|------|---------|
| `github.com/owner/repo` | Clone `https://github.com/owner/repo.git` at default branch |
| `github.com/owner/repo@v1.2.0` | Clone and check out tag `v1.2.0` |
| `https://example.com/path/repo.git` | Clone any HTTPS URL |
| `git@github.com:owner/repo.git` | SSH clone |

Optionally pin a ref with `-ref`:

```bash
dot plugin install github.com/me/my-plugin -ref v0.2.0
```

### Install a local plugin (development)

```bash
dot plugin install -from ./my-local-plugin
```

### Uninstall a plugin

```bash
dot plugin uninstall my-plugin
```

Plugins are stored in `~/.dot/plugins/<id>/`. Uninstall removes that directory.

> **Note:** DOT is a compiled binary. After installing a plugin, you need to rebuild `dot` with the plugin imported for its `init()` function to register hooks. This constraint is removed when dynamic loading is introduced.

---

## Uninstall

### macOS — Homebrew

```bash
brew uninstall dot
```

### curl / go install / from source

```bash
curl -fsSL https://raw.githubusercontent.com/version14/dot/main/uninstall.sh | sh
```

Project `.dot/` directories are **not** removed. Delete them manually if needed:

```bash
rm -rf my-project/.dot
```

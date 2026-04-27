package state

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

// GoMod models the minimal subset of a go.mod file the engine needs to
// manipulate: the module path, Go version, and require directives.
//
// It intentionally avoids modfile parsing to keep generator-time edits cheap
// and predictable; the file is regenerated deterministically on Marshal.
type GoMod struct {
	Module    string
	GoVersion string
	Requires  []GoModRequire
}

type GoModRequire struct {
	Path    string
	Version string
}

func NewGoMod() *GoMod { return &GoMod{} }

func (m *GoMod) Load(data []byte) error {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	inBlock := false
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}
		switch {
		case strings.HasPrefix(line, "module "):
			m.Module = strings.TrimSpace(strings.TrimPrefix(line, "module "))
		case strings.HasPrefix(line, "go "):
			m.GoVersion = strings.TrimSpace(strings.TrimPrefix(line, "go "))
		case strings.HasPrefix(line, "require ("):
			inBlock = true
		case inBlock && line == ")":
			inBlock = false
		case strings.HasPrefix(line, "require "):
			parts := strings.Fields(strings.TrimPrefix(line, "require "))
			if len(parts) >= 2 {
				m.Requires = append(m.Requires, GoModRequire{Path: parts[0], Version: parts[1]})
			}
		case inBlock:
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				m.Requires = append(m.Requires, GoModRequire{Path: parts[0], Version: parts[1]})
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("gomod: scan: %w", err)
	}
	return nil
}

func (m *GoMod) AddModule(path string) { m.Module = path }
func (m *GoMod) SetGoVersion(v string) { m.GoVersion = v }

func (m *GoMod) AddRequire(path, version string) {
	for i, r := range m.Requires {
		if r.Path == path {
			m.Requires[i].Version = version
			return
		}
	}
	m.Requires = append(m.Requires, GoModRequire{Path: path, Version: version})
}

func (m *GoMod) RemoveRequire(path string) {
	out := m.Requires[:0]
	for _, r := range m.Requires {
		if r.Path != path {
			out = append(out, r)
		}
	}
	m.Requires = out
}

func (m *GoMod) Marshal() []byte {
	var b bytes.Buffer
	if m.Module != "" {
		fmt.Fprintf(&b, "module %s\n\n", m.Module)
	}
	if m.GoVersion != "" {
		fmt.Fprintf(&b, "go %s\n", m.GoVersion)
	}
	if len(m.Requires) > 0 {
		b.WriteString("\nrequire (\n")
		for _, r := range m.Requires {
			fmt.Fprintf(&b, "\t%s %s\n", r.Path, r.Version)
		}
		b.WriteString(")\n")
	}
	return b.Bytes()
}

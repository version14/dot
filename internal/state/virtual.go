package state

import (
	"time"

	"github.com/version14/dot/internal/spec"
)

// VirtualProjectState — in-memory filesystem, all generators write here before disk

type VirtualProjectState struct {
	Files    map[string]*FileNode
	Metadata spec.ProjectMetadata
}

// internal/state/file.go
type FileNode struct {
	Path            string
	Content         []byte
	ContentType     ContentType // Raw, JSON, YAML, GoMod
	CreatedBy       string      // generator name
	Transformations []string    // audit trail of edits
	ModifiedAt      time.Time
}

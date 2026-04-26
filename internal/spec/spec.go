package spec

import (
	"time"
)

// ProjectSpec definition — flow ID, metadata, recursive answer tree, visited nodes, constraints

type ProjectSpec struct {
	FlowID               string                `json:"flow_id"`
	CreatedAt            time.Time             `json:"created_at"`
	Metadata             ProjectMetadata       `json:"metadata"`
	Answers              map[string]AnswerNode `json:"answers"`       // recursive tree
	VisitedNodes         []string              `json:"visited_nodes"` // traversal audit trail
	LoadedPlugins        []string              `json:"loaded_plugins"`
	GeneratorConstraints map[string]string     `json:"generator_constraints"`
}

type ProjectMetadata struct {
	ProjectName string `json:"project_name"`
	ToolVersion string `json:"tool_version"`
}

package spec

import "time"

// AnswerNode mirrors flow.AnswerNode but lives here to keep the spec package
// importable without pulling in the flow package (avoids cycles when generators
// only need to read the spec).
//
// Concrete types: string | bool | int | []string | map[string]AnswerNode | []map[string]AnswerNode
type AnswerNode = interface{}

// ProjectSpec is the persisted shape of a scaffolded project. It is written
// to .dot/spec.json after generation and re-read on `dot scaffold` re-runs.
type ProjectSpec struct {
	FlowID               string                `json:"flow_id"`
	CreatedAt            time.Time             `json:"created_at"`
	Metadata             ProjectMetadata       `json:"metadata"`
	Answers              map[string]AnswerNode `json:"answers"`
	VisitedNodes         []string              `json:"visited_nodes"`
	LoadedPlugins        []string              `json:"loaded_plugins"`
	GeneratorConstraints map[string]string     `json:"generator_constraints"`
}

type ProjectMetadata struct {
	ProjectName string `json:"project_name"`
	ToolVersion string `json:"tool_version"`
}

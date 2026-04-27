package servicewriter

import "github.com/version14/dot/pkg/dotapi"

// Manifest declares service_writer — a per-iteration generator that scaffolds
// one service folder. The microservices flow runs it once per service in a
// LoopQuestion, with its loop scope providing the service's name + port.
var Manifest = dotapi.Manifest{
	Name:        "service_writer",
	Version:     "0.1.0",
	Description: "Writes one service skeleton per loop iteration (services/<name>/)",
	DependsOn:   []string{"base_project"},
	// Outputs are dynamic (based on loop scope); leave empty.
	Outputs: nil,
}

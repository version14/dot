package postgresdockercompose

import "github.com/version14/dot/pkg/dotapi"

var Manifest = dotapi.Manifest{
	Name:        "postgres_docker_compose",
	Version:     "0.1.0",
	Description: "Docker Compose service for PostgreSQL development environment",
	DependsOn:   []string{},
	Outputs: []string{
		"docker-compose.yml",
	},
	Validators: []dotapi.Validator{
		{
			Name: "postgres-docker-compose",
			Checks: []dotapi.Check{
				{Type: dotapi.CheckFileExists, Path: "docker-compose.yml"},
			},
		},
	},
}

// Package portutil derives stable, per-project host ports for generated
// docker-compose services, so two scaffolded projects don't collide on the
// same fixed port when run side by side.
package portutil

import "hash/fnv"

// postgresBase and postgresSpan bound the host ports handed out to
// scaffolded Postgres containers: [40000, 50000).
const (
	postgresBase = 40000
	postgresSpan = 10000
)

// PostgresHostPort deterministically maps a project name to a host port for
// its Postgres container. The same name always yields the same port, so the
// docker-compose.yml and .env.example generators can each call this
// independently and stay in agreement.
func PostgresHostPort(projectName string) int {
	h := fnv.New32a()
	_, _ = h.Write([]byte(projectName))
	return postgresBase + int(h.Sum32()%postgresSpan)
}

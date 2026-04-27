package cli_test

import (
	"github.com/version14/dot/internal/cli"
	"github.com/version14/dot/internal/flow"
	"testing"
)

func TestHuhFormRunnerImplementsFlowRunner(t *testing.T) {
	var _ flow.FlowRunner = cli.NewHuhFormRunner()
}

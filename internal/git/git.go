package git

import (
	"os"
	"os/exec"
)

type Repository struct{}

func (r *Repository) Close() {
	// No resources to clean up for now, but this allows us to add some in the future without changing the API.
}

func Init(path, defaultBranch string) (*Repository, error) {
	cmd := exec.Command("git", "init", "-b", defaultBranch)
	cmd.Dir = path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	return &Repository{}, nil
}

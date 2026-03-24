package internal

import (
	"errors"
	"os/exec"
	"strings"
)

func ListContent(repo string) ([]string, error) {
	if repo == "" {
		return nil, errors.New("repo folder doesn't exist")
	}

	cmd := exec.Command(
		"git", "-C", repo,
		"ls-tree", "-r", "HEAD", "--name-only",
	)

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	files := strings.Split(strings.TrimSpace(string(output)), "\n")

	return files, nil
}
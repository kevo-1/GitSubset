package internal

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func Clone(URL string) (GithubLink, error) {
	link, err := ParseURL(URL)
	
	if err != nil {
		return GithubLink{}, fmt.Errorf("Invalid repo link\nerror: %w", err)
	}

	ok, err := checkRepoExist(URL)

	if err != nil {
		return GithubLink{}, err
	}
	if !ok {
		return GithubLink{}, errors.New("repository not found")
	}

	dest, err := filepath.Abs(link.Repo)
	
	if err != nil {
		return GithubLink{}, err
	}

	if _, statErr := os.Stat(dest); os.IsNotExist(statErr) {
		args := []string{"git", "clone", "--filter=blob:none", "--no-checkout", URL}
		cmd := exec.Command(args[0], args[1:]...)

		err = cmd.Run()
		if err != nil {
			return GithubLink{}, err
		}
	}

	link.Path = dest

	return link, nil
}

func checkRepoExist(URL string) (bool, error) {
	cmd := exec.Command("git", "ls-remote", "--exit-code", "-h", URL)

	err := cmd.Run()
	if err != nil {
		return false, fmt.Errorf("Repo not found\nerror: %w", err)
	}

	return true, nil
}
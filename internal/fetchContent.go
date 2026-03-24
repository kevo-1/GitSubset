package internal

import (
	"os/exec"
	"strings"
)

func FetchContent(repo string, files []string) error {
	existing := getExistingPatterns(repo)

	merged := make(map[string]bool)
	for _, f := range existing {
		merged[f] = true
	}
	for _, f := range files {
		merged[f] = true
	}

	all := make([]string, 0, len(merged))
	for f := range merged {
		all = append(all, f)
	}

	args := []string{"-C", repo, "sparse-checkout", "set", "--no-cone"}
	args = append(args, all...)

	cmd := exec.Command("git", args...)
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("git", "-C", repo, "checkout")
	return cmd.Run()
}

func getExistingPatterns(repo string) []string {
	cmd := exec.Command("git", "-C", repo, "sparse-checkout", "list")
	output, err := cmd.Output()
	if err != nil {
		return nil
	}

	raw := strings.TrimSpace(string(output))
	if raw == "" {
		return nil
	}

	return strings.Split(raw, "\n")
}
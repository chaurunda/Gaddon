package git

import (
	"log"
	"os/exec"
	"strings"
)

func Check_last_local_commit_id(path string) string {
	cmd := exec.Command("git", "-C", path, "rev-parse", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(out), "\n")

	return lines[0]
}

func Check_last_distant_commit_id(path string) string {
	cmd := exec.Command("git", "-C", path, "rev-parse", "origin/HEAD")
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(out), "\n")

	return lines[0]
}

func Install_git_repo(git_url string, path string) {
	cmd := exec.Command("git", "clone", git_url, path)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func Update_git_repo(path string) {
	cmd := exec.Command("git", "-C", path, "pull")
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

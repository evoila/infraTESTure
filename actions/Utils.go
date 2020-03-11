package actions

import (
	"os/exec"
	"strings"
)

func gitClone(url string, repoPath string, tag string) error {
	var tagClone string

	if tag != "" {
		tagClone = "--branch " + tag + " --single-branch"
	}

	cmd := exec.Command("bash", "-c", "git clone "+url+" "+repoPath+" "+tagClone)
	err := cmd.Run()

	if err != nil {
		return err
	}

	return nil
}

func appendSlash(dir string) string {
	if !strings.HasSuffix(dir, "/") {
		return dir + "/"
	}
	return dir
}

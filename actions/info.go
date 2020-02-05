package actions

import (
	"github.com/evoila/infraTESTure/parser"
	"github.com/fatih/color"
	"github.com/urfave/cli"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"sort"
	"strings"
)

func Info(c *cli.Context) error {
	url := c.String("repository")
	tag := c.String("tag")
	tmpDir := appendSlash(os.TempDir()) + "infra-tmp"

	err := gitClone(url, tmpDir, tag)

	if err != nil {
		logError(err, "Could not clone repository")
		return err
	}

	log.Printf("The following services and tests were found in %v:\n\n", url)

	dirs, err := ioutil.ReadDir(tmpDir)
	if err != nil {
		logError(err, "No repository found")
		return err
	}

	// Iterate through all directories and files inside these directories to get a list
	// of all annotations, and therefore a list of all offered tests
	for _, dir := range dirs {
		if dir.IsDir() {
			goFiles, err := ioutil.ReadDir(appendSlash(tmpDir)+dir.Name())
			if err != nil {
				logError(err, "Failed to acquire offered tests")
				return err
			}

			if !strings.HasPrefix(dir.Name(), ".") {
				log.Printf("├── %v", color.GreenString(dir.Name()))
			}

			var testNames []string

			for _, goFile := range goFiles {
				if strings.HasSuffix(goFile.Name(), ".go") {
					tmpTestNames := parser.GetAnnotations(appendSlash(tmpDir)+dir.Name()+"/"+goFile.Name())

					for i := range tmpTestNames {
						testNames = append(testNames, tmpTestNames[i])
					}
				}
			}

			sort.Strings(testNames)

			for i := range testNames {
				log.Printf("│ \t├── %v", testNames[i])
			}
		}
	}

	cmd := exec.Command("bash", "-c", "rm -rf " + tmpDir)
	err = cmd.Run()

	if err != nil {
		logError(err, "Could not delete directory")
		return err
	}

	return nil
}

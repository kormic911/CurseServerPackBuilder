package packtester

import (
	"os/exec"
	"os"
	"strings"
	"fmt"
)

var serverDir string
var forge string

func InstallForge(modServerDir string, forgeInstaller string) error {
	serverDir = modServerDir
	forge = forgeInstaller
	os.Chdir(modServerDir)
	installer := exec.Command("java", "-jar", forgeInstaller, "--installServer")
	err := installer.Start()
	if err != nil {
		return err
	}

	installer.Wait()
	return nil
}

func Start() error {
	makeEula()

	forgeUniversal := strings.Replace(forge, "installer", "universal", -1)
	fmt.Println(forgeUniversal)
	return nil
}

func makeEula() error {
	return nil
}
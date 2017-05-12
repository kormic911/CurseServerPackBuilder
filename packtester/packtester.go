package packtester

import (
	"os/exec"
	"os"
	"strings"
	"fmt"
	"bufio"
	"io/ioutil"
	"path/filepath"
	"io"
	"regexp"
)

var serverDir string
var forge string

func InstallForge(modServerDir string, forgeInstaller string) error {
	modServerDir, err := filepath.Abs(modServerDir)
	serverDir = modServerDir
	forge = forgeInstaller
	os.Chdir(modServerDir)
	installer := exec.Command("java", "-jar", forgeInstaller, "--installServer")
	stdout, err := installer.StdoutPipe()

	err = installer.Start()
	if err != nil {
		return err
	}

	r := bufio.NewReader(stdout)
	for {
		_, _, err := r.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		//fmt.Println("** ", string(line[:]))
	}

	installer.Wait()
	return nil
}

func Start() error {
	makeEula()

	forgeUniversal := strings.Replace(forge, "installer", "universal", -1)
	os.Chdir(serverDir)
	server := exec.Command("java", "-jar", forgeUniversal)
	stdout, err := server.StdoutPipe()
	if err != nil {
		return err
	}

	err = server.Start()
	if err != nil {
		return err
	}

	r := bufio.NewReader(stdout)
	for {
		line, _, err := r.ReadLine()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		matched, err := regexp.Match("Done \\(", line)
		if err != nil {
			return err
		}

		fmt.Println(string(line[:]))
		if matched {
			fmt.Println("Server Finished")
			server.Process.Kill()
		}
	}

	return nil
}

func makeEula() error {
	eulaText := []byte("eula=true\n")
	err := ioutil.WriteFile(filepath.FromSlash(serverDir+"/eula.txt"), eulaText, 0644)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
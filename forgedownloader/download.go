package forgedownloader

import (
	"github.com/kormic911/CurseServerPackBuilder/manifestreader"
	"fmt"
	"strings"
	"net/http"
	"os"
	"path/filepath"
	"io"
	"encoding/json"
)

var forgeVersions string = "http://files.minecraftforge.net/maven/net/minecraftforge/forge/json"
var forgeDownloadURL string = "http://files.minecraftforge.net/maven/net/minecraftforge/forge"

func DownloadForge(minecraft manifestreader.MinecraftType, modServerDir string) (forgeInstaller string, err error) {
	var modloaderVersion string
	for _, modLoader := range minecraft.Modloaders {
		if modLoader.Primary {
			modloaderVersion = strings.Replace(modLoader.Id, "forge-", "", -1)
		}
	}


	forgeVersion, err := forgeDownloadPath(modloaderVersion)

	forgeInstaller = fmt.Sprintf("forge-%s-installer.jar", forgeVersion )

	fmt.Printf("Downloading %s for %s\n", forgeInstaller, minecraft.Version)
	response, err := http.Get(fmt.Sprintf("%s/%s/%s", forgeDownloadURL, forgeVersion, forgeInstaller ))
	if err != nil {
		return
	}

	output, err := os.Create(filepath.FromSlash(fmt.Sprintf("%s/%s", modServerDir, forgeInstaller)))
	if err != nil {
		return
	}

	defer output.Close()
	defer response.Body.Close()

	_, err = io.Copy(output, response.Body)
	if err != nil {
		return
	}

	return
}

func forgeDownloadPath(forgeVersion string) (string, error) {
	forgeVersionParts := strings.Split(forgeVersion, ".")
	baseVersion := forgeVersionParts[len(forgeVersionParts)-1]

	response, err := http.Get(forgeVersions)
	if err != nil {
		fmt.Println("error: ",err)
		return "", err
	}

	defer response.Body.Close()

	var f map[string]interface{}
	json.NewDecoder(response.Body).Decode(&f)
	m := f["number"].(map[string]interface{})
	build := m[baseVersion].(map[string]interface{})

	if build["branch"] != nil {
		return fmt.Sprintf("%s-%s-%s", build["mcversion"], build["version"], build["branch"]), nil
	} else {
		return fmt.Sprintf("%s-%s", build["mcversion"], build["version"]), nil
	}
	//forge-<mcversion>-<version>[-<branch>]
}
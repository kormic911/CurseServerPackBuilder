package forgedownloader

import (
	"github.com/kormic911/CurseServerPackBuilder/manifestreader"
	"fmt"
	"strings"
	"net/http"
	"os"
	"path/filepath"
	"io"
)

var forgeVersions string = "http://files.minecraftforge.net/maven/net/minecraftforge/forge"

func DownloadForge(minecraft manifestreader.MinecraftType, modServerDir string) (forgeInstaller string, err error) {
	var modloaderVersion string
	for _, modLoader := range minecraft.Modloaders {
		if modLoader.Primary {
			modloaderVersion = strings.Replace(modLoader.Id, "forge-", "", -1)
		}
	}

	forgeInstaller = fmt.Sprintf("forge-%s-%s-installer.jar", minecraft.Version, modloaderVersion )

	fmt.Printf("Downloading forge[%s] for %s\n", modloaderVersion, minecraft.Version)
	response, err := http.Get(fmt.Sprintf("%s/%s-%s/%s", forgeVersions, minecraft.Version, modloaderVersion, forgeInstaller ))
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
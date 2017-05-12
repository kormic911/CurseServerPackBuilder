package main
import (
	"fmt"
	"github.com/kormic911/CurseServerPackBuilder/zipper"
	"github.com/kormic911/CurseServerPackBuilder/manifestreader"
	"github.com/kormic911/CurseServerPackBuilder/cursedownloader"
	"github.com/kormic911/CurseServerPackBuilder/forgedownloader"
	"github.com/kormic911/CurseServerPackBuilder/packtester"
	"flag"
	"path/filepath"
	"os"
	"strings"
	"io"
)

func main() {
	modArchive := flag.String("mod", "", "The modpack archive to turn into a server package")
	modBaseDir := flag.String("dir", "", "The directory to do all the work in")
	modFlushCacheDir := flag.Bool("flushcache", false, "Flush the cache of downloaded mods")

	flag.Parse()

	if *modArchive == "" || *modBaseDir == "" {
		flag.PrintDefaults()
		return
	}

	modPackDir := filepath.FromSlash(*modBaseDir + "/modpack")
	modServerDir := filepath.FromSlash(*modBaseDir + "/server")
	modServerTestDir := filepath.FromSlash(*modBaseDir + "/servertest")
	modCacheDir := filepath.FromSlash(*modBaseDir + "/cache")

	err := os.RemoveAll(modPackDir)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = os.RemoveAll(modServerDir)
	if err != nil {
		fmt.Println(err)
		return
	}

	if *modFlushCacheDir {
		err = os.RemoveAll(modCacheDir)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	err = os.RemoveAll(modServerTestDir)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Extracting mod archive")
	err = zipper.Unzip(*modArchive, modPackDir)
	if err != nil {
		fmt.Println(err)
		return
	}

	manifestreader.Read(filepath.FromSlash(modPackDir + "/manifest.json"))
	fmt.Println("Cloning mod overrides")
	err = cloneDir(filepath.FromSlash(modPackDir + "/" + manifestreader.OverridesDir()), modServerDir)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = cursedownloader.DownloadMods(manifestreader.Modlist(), filepath.FromSlash(modServerDir + "/mods"), modCacheDir)
	if err != nil {
		fmt.Println(err)
		return
	}

	forgeInstaller, err := forgedownloader.DownloadForge(manifestreader.Minecraft(), modServerDir)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Creating test environment")
	err = cloneDir(modServerDir, modServerTestDir)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = packtester.InstallForge(modServerTestDir, forgeInstaller)
	if err != nil {
		fmt.Println("ERROR: ",err)
		return
	}

	packtester.Start()
}

func cloneDir(modPackOverridesDir, modServerDir string) error {
	err := filepath.Walk(modPackOverridesDir, copyDir(modServerDir, modPackOverridesDir))
	if err != nil {
		return err
	}
	return nil
}

func copyDir(modServerDir, modPackOverridesDir string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return nil
		}

		newpath := strings.Replace(path, modPackOverridesDir, modServerDir, -1)

		if info.IsDir() {
			err := os.MkdirAll(newpath, 0755)
			if err != nil {
				fmt.Println(err)
				return err
			}
		} else {
			src, err := os.Open(path)
			if err != nil {
				fmt.Println(err)
				return err
			}

			defer src.Close()

			dest, err := os.Create(newpath)
			if err != nil {
				fmt.Println(err)
				return err
			}

			defer dest.Close()

			_, err = io.Copy(dest, src)
			if err != nil {
				fmt.Println(err)
				return err
			}
		}

		return nil
	}
}

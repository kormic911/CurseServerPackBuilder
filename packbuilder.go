package main
import (
	"fmt"
	"zipper"
	"manifestreader"
	"cursedownloader"
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

	modPackDir := *modBaseDir + "/modpack"
	modServerDir := *modBaseDir + "/server"
	modCacheDir := *modBaseDir + "/cache"

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

	err = zipper.Unzip("/home/jwilson/sample/FTBBeyond-1.6.0-1.10.2.zip", modPackDir)
	if err != nil {
		fmt.Println(err)
		return
	}

	manifestreader.Read(modPackDir + "/manifest.json")
	err = cloneOverrides(modPackDir + "/" + manifestreader.OverridesDir(), modServerDir)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = cursedownloader.DownloadMods(manifestreader.Modlist(), modServerDir + "/mods", modCacheDir)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func cloneOverrides(modPackOverridesDir, modServerDir string) error {
	err := filepath.Walk(modPackOverridesDir, copyOverride(modServerDir, modPackOverridesDir))
	if err != nil {
		return err
	}
	return nil
}

func copyOverride(modServerDir, modPackOverridesDir string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		newpath := strings.Replace(path, modPackOverridesDir, modServerDir, -1)

		if info.IsDir() {
			err := os.MkdirAll(newpath, 0755)
			if err != nil {
				return err
			}
		} else {
			src, err := os.Open(path)
			if err != nil {
				return err
			}

			defer src.Close()

			dest, err := os.Create(newpath)
			if err != nil {
				return err
			}

			defer dest.Close()

			_, err = io.Copy(dest, src)
			if err != nil {
				return err
			}
		}

		return nil
	}
}

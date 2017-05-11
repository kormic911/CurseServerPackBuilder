package cursedownloader

import (
	"fmt"
	"manifestreader"
	"net/http"
	"net/url"
	"strings"
	"os"
	"io"
	"github.com/gosuri/uilive"
	"time"
)

func DownloadMods(mods []manifestreader.FileType, destPath, cachePath string) error {
	writer := uilive.New()
	writer.Start()
	for index, mod := range mods {
		fmt.Fprintf(writer, "Downloading (%d/%d) mods\n", (index + 1), len(mods))
		err := downloadFile(mod, destPath, cachePath)
		if err != nil {
			return err
		}
	}
	time.Sleep(time.Millisecond * 10)
	fmt.Fprintf(writer, "Finished: Downloaded %d mods\n", len(mods))
	writer.Stop()
	return nil
}

func downloadFile(mod manifestreader.FileType, destPath, cachePath string) error {
	cachePathDir := fmt.Sprintf("%s/%d/%d", cachePath, mod.ProjectID, mod.FileID)
	if _, err := os.Stat(cachePathDir); os.IsNotExist(err) {
		modUrl := fmt.Sprintf("https://minecraft.curseforge.com/projects/%d/files/%d/download", mod.ProjectID, mod.FileID)
		response, err := http.Get(modUrl)
		if err != nil { 
			return err
		}

		downloadPath, err := url.PathUnescape(response.Request.URL.String())
		if err != nil {
			return err
		}

		tokens := strings.Split(downloadPath, "/")
		fileName := tokens[len(tokens)-1]

		if _, err := os.Stat(cachePathDir); os.IsNotExist(err) {
			err := os.MkdirAll(cachePathDir, 0755)
			if err != nil {
				return err
			}
		}

		output, err := os.Create(fmt.Sprintf("%s/%s", cachePathDir, fileName))
		if err != nil {
			return err
		}

		defer output.Close()
		defer response.Body.Close()

		_, err = io.Copy(output, response.Body)
		if err != nil {
			return err
		}
		copyFromCache(destPath, cachePathDir)
	} else {
		copyFromCache(destPath, cachePathDir)
		time.Sleep(time.Millisecond * 5)
	}

	return nil
}

func copyFromCache(destPath, cachePath string) error {
	return nil
}
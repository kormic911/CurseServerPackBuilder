package cursedownloader

import (
	"fmt"
	"github.com/kormic911/CurseServerPackBuilder/manifestreader"
	"net/http"
	"net/url"
	"net/http/cookiejar"
	"strings"
	"os"
	"io"
	"github.com/gosuri/uilive"
	"time"
	"path/filepath"
	"errors"
)

var cookieJar *cookiejar.Jar

func DownloadMods(mods []manifestreader.FileType, destPath, cachePath string) error {
	cookieJar, _ = cookiejar.New(nil)
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
		httpClient := &http.Client{
			Jar: cookieJar,
		}
		response, err := httpClient.Get(modUrl)
		if err != nil { 
			return err
		}

		downloadPath, err := url.PathUnescape(response.Request.URL.String())
		if err != nil {
			return err
		}

		tokens := strings.Split(downloadPath, "/")
		fileName := tokens[len(tokens)-1]

		if fileName == "download" {
			return errors.New("Mod missing/not available for download.")
		}

		if _, err := os.Stat(cachePathDir); os.IsNotExist(err) {
			err := os.MkdirAll(cachePathDir, 0755)
			if err != nil {
				return err
			}
		}

		output, err := os.Create(filepath.FromSlash(fmt.Sprintf("%s/%s", cachePathDir, fileName)))
		if err != nil {
			return err
		}

		defer output.Close()
		defer response.Body.Close()

		_, err = io.Copy(output, response.Body)
		if err != nil {
			return err
		}
		copyFromCache(filepath.FromSlash(destPath), filepath.FromSlash(cachePathDir))
	} else {
		copyFromCache(filepath.FromSlash(destPath), filepath.FromSlash(cachePathDir))
		time.Sleep(time.Millisecond * 5)
	}

	return nil
}

func copyFromCache(destPath, cachePath string) error {
	err := filepath.Walk(cachePath, copyPath(destPath, cachePath))
	if err != nil {
		return err
	}
	return nil
}

func copyPath(destPath, cachePath string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return nil
		}

		newpath := strings.Replace(path, cachePath, destPath, -1)

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

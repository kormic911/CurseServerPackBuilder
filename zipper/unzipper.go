package zipper

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

func Unzip(archive, target string) error {
	reader, err := zip.OpenReader(archive)

	if err != nil {
		return err
	}

	if err := os.MkdirAll(target, 0755); err != nil {
		return err
	}

	for _, file := range reader.File {
		path := filepath.Join(target, file.Name)
		basePath := filepath.Dir(path)
		// Sometimes zips do not reutrn folders as part of the list of 'files' so we just get the base path of every file and make that directory
		if err := os.MkdirAll(basePath, 0755); err != nil {
			return err
		}

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(path, file.Mode()); err != nil {
				return err
			}
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			if fileReader != nil {
				fileReader.Close()
			}
			return err
		}

		fileWriter, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			fileReader.Close()
			if fileWriter != nil {
				fileWriter.Close()
			}

			return err
		}

		if _, err := io.Copy(fileWriter, fileReader); err != nil {
			fileReader.Close()
			fileWriter.Close()

			return err
		}

		fileReader.Close()
		fileWriter.Close()
	}

	return nil
} 

package archive

import (
	"archive/zip"
	"github.com/pkg/errors"
	"io"
	"os"
	"path/filepath"
)

func ZipFiles(filename string, files []string) error {
	newZipFile, err := os.Create(filename)
	if err != nil {
		return errors.WithStack(err)
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	for _, file := range files {
		if err = addFileToZip(zipWriter, file); err != nil {
			return err
		}
	}
	return nil
}

func addFileToZip(zipWriter *zip.Writer, filename string) error {
	fileToZip, err := os.Open(filename)
	if err != nil {
		return errors.WithStack(err)
	}
	defer fileToZip.Close()

	info, err := fileToZip.Stat()
	if err != nil {
		return errors.WithStack(err)
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return errors.WithStack(err)
	}

	header.Name = filepath.Base(filename)

	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return errors.WithStack(err)
	}

	if _, err := io.Copy(writer, fileToZip); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

package csv

import (
	"encoding/csv"
	"github.com/pkg/errors"
	"golang.org/x/text/encoding/charmap"
	"os"
	"path/filepath"
)

func Writer(data [][]string, outPath string, winEncoding bool) error {
	if err := checkDirExists(outPath); err != nil {
		return err
	}

	file, err := os.Create(outPath)
	if err != nil {
		return errors.WithStack(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = ';'

	if winEncoding {
		var dataWin [][]string

		encoder := charmap.Windows1251.NewEncoder()

		for _, row := range data {
			var rowWin []string

			for _, str := range row {
				strWin, err := encoder.String(str)
				if err != nil {
					return errors.WithStack(err)
				}

				rowWin = append(rowWin, strWin)
			}

			dataWin = append(dataWin, rowWin)
		}

		if err := writer.WriteAll(dataWin); err != nil {
			return errors.WithStack(err)
		}
	} else {
		if err := writer.WriteAll(data); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func checkDirExists(outPath string) error {
	baseDir := filepath.Dir(outPath)

	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		if err := os.MkdirAll(baseDir, 0644); err != nil {
			return errors.WithStack(err)
		}
	} else {
		return errors.WithStack(err)
	}

	return nil
}

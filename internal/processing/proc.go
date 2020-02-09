package processing

import (
	"database/sql"
	"github.com/pkg/errors"
	"log"
	"loyaltyReport/internal/config"
	"loyaltyReport/pkg/archive"
	"loyaltyReport/pkg/csv"
	"loyaltyReport/pkg/mail"
	"os"
	"path/filepath"
	"strings"
)

func GenerateReport(conf config.Config, db *sql.DB, name string, outPath string, log *log.Logger) error {
	log.Printf("Find report \"%v\" in config", name)

	report, ok := conf.Reports[name]
	if !ok {
		return errors.Errorf("Report with name %v not found in configuration file", name)
	}

	log.Println("Fetch data from DB")
	queryResult, err := getFromDB(db, report.Query)
	if err != nil {
		return err
	}

	csvPath := getCSVPath(outPath)

	log.Printf("Create temp CSV file: %v\n", csvPath)
	if err := csv.Writer(queryResult, csvPath, report.WinEncoding); err != nil {
		return err
	}

	log.Printf("Create archive: %v\n", outPath)
	if err := archive.ZipFiles(outPath, []string{csvPath}); err != nil {
		return err
	}

	if err := os.Remove(csvPath); err != nil {
		return errors.WithStack(err)
	}

	mailServer := mail.Server{
		Host:     conf.Mail.Host,
		Port:     conf.Mail.Port,
		Username: conf.Mail.Username,
		Password: conf.Mail.Password,
	}

	log.Printf("Send report \"%v\" to %v\n", name, report.SendTo)

	if err := mail.Send(mailServer, conf.Mail.SendFrom, report.SendTo, report.Subject, report.Text, []string{outPath}); err != nil {
		return err
	}

	return nil
}

func nullToStr(nullStr []sql.NullString) []string {
	var result []string
	for _, v := range nullStr {
		result = append(result, v.String)
	}

	return result
}

func getFromDB(db *sql.DB, query string) ([][]string, error) {
	var queryResult [][]string

	rows, err := db.Query(query)
	if err != nil {
		return queryResult, errors.WithStack(err)
	}

	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return queryResult, errors.WithStack(err)
	}

	queryResult = append(queryResult, columns)

	rowRead := make([]interface{}, len(columns))
	rowWrite := make([]sql.NullString, len(columns))

	for i := 0; i < len(rowRead); i++ {
		rowRead[i] = &rowWrite[i]
	}

	for rows.Next() {
		if err := rows.Scan(rowRead...); err != nil {
			return queryResult, errors.WithStack(err)
		}
		queryResult = append(queryResult, nullToStr(rowWrite))
	}

	return queryResult, nil
}

func getCSVPath(zipPath string) string {
	dir, file := filepath.Split(zipPath)
	ext := filepath.Ext(file)

	fileName := strings.TrimSuffix(file, ext) + ".csv"

	pathCSV := filepath.Join(dir, fileName)

	return pathCSV
}

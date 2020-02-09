package main

import (
	"flag"
	"fmt"
	"github.com/pkg/errors"
	"loyaltyReport/internal/config"
	"loyaltyReport/internal/processing"
	"loyaltyReport/pkg/database"
	"loyaltyReport/pkg/logger"
	"os"
	"path/filepath"
	"time"
)

func main() {
	rootDir := filepath.Dir(os.Args[0])
	outDir := filepath.Join(rootDir, "output")

	reportName := flag.String("n", "", "Report name")
	flag.Parse()

	if *reportName == "" {
		fmt.Println(`Missing required launch flags, see "-h"`)
		os.Exit(1)
	}

	reportFile := *reportName + "_" + time.Now().Format("2006-01-02") + ".zip"

	log, err := logger.GetLogger(filepath.Join(rootDir, "app.log"))
	if err != nil {
		fmt.Printf("%+v\n", err)
	}

	conf, err := config.GetConfig(filepath.Join(rootDir, "env.yml"), log)
	if err != nil {
		log.Fatalf("ERROR - %+v\n", err)
	}

	db, err := database.Open(conf.DataSource.URL, conf.DataSource.User, conf.DataSource.Password, 1, log)
	if err != nil {
		log.Fatalf("ERROR - %+v\n", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("%+v\n", errors.WithStack(err))
		}
		log.Println("Close database connection")
	}()

	if err := processing.GenerateReport(conf, db, *reportName, filepath.Join(outDir, reportFile), log); err != nil {
		log.Fatalf("ERROR - %+v\n", err)
	}
}

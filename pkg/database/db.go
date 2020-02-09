package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"log"
	"regexp"
)

func Open(url string, user string, password string, maxConn int, log *log.Logger) (*sql.DB, error) {
	connStr, err := getConnStr(url, user, password)
	if err != nil {
		return nil, err
	}

	log.Println("Open database connection")

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := db.Ping(); err != nil {
		return nil, errors.WithStack(err)
	}

	db.SetMaxIdleConns(maxConn)
	db.SetMaxOpenConns(maxConn)

	return db, nil
}

func getConnStr(url string, user string, password string) (string, error) {
	var connStr string

	regex, err := regexp.Compile(`postgresql://([^:]+):(\d+)/([^?\n]+)`)
	if err != nil {
		return connStr, errors.WithStack(err)
	}

	group := regex.FindStringSubmatch(url)

	if len(group) != 4 {
		return connStr, errors.New("Error parsing connection string")
	}

	connStr = fmt.Sprintf("host=%v port=%v dbname=%v user=%v password=%v sslmode=disable",
		group[1], group[2], group[3], user, password)

	return connStr, nil
}

package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Connect connects to a MySQL database
func Connect(user string, pass string, dbname string, address string) (*sql.DB, error) {
	// build the dsn string
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true",
		user, pass, address, dbname)

	// connect to database
	db, err := sql.Open("mysql", dsn)
	if err != nil || db == nil {
		return nil, fmt.Errorf("connect dsn[%s]: %v", dsn, err)
	}

	// test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("connect dsn[%s]: %v", dsn, err)
	}

	log.Printf("Connected to database at [%s]/[%s]", address, dbname)

	return db, nil
}

// Migrate sets up database tables
func Migrate(db *sql.DB) error {
	const (
		migrationDir = "migration"
	)

	if db == nil {
		return fmt.Errorf("Migrate: null pointer error")
	}

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	migrationPath := fmt.Sprintf("file:///%s", filepath.Join(exPath, migrationDir))

	driver, _ := mysql.WithInstance(db, &mysql.Config{})

	m, err := migrate.NewWithDatabaseInstance(
		migrationPath,
		"mysql",
		driver,
	)

	if err != nil {
		return err
	}

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

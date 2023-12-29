// utils/database.go

package utils

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func InitDB() *sql.DB {
	dbUser := GetEnv("DB_USER")
	dbPassword := GetEnv("DB_PASSWORD")
	dbHost := GetEnv("DB_HOST")
	dbPort := GetEnv("DB_PORT")
	dbName := GetEnv("DB_NAME")

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	conn, err := sql.Open("mysql", connectionString)
	if err != nil {
		fmt.Println("DB Connection Error:", err)
		os.Exit(1)
	}

	db = conn

	// Migrate the database
	MigrateDB()

	return db
}

func GetDB() *sql.DB {
	return db
}

func MigrateDB() {
	// Read and execute migration SQL script
	migrationScript, err := ioutil.ReadFile("migration.sql")
	if err != nil {
		fmt.Println("Error reading migration script:", err)
		os.Exit(1)
	}

	_, err = db.Exec(string(migrationScript))
	if err != nil {
		fmt.Println("Error executing migration script:", err)
		os.Exit(1)
	}
}

package database

import (
	"database/sql"
	"embed"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
)

//go:embed sql_migrations/*.sql
var dbMigrations embed.FS

var DB *sql.DB

func Connect() error {
	// Load .env file
	err := godotenv.Load("config/.env")
	if err != nil {
		fmt.Println("Warning: Error loading .env file, using system environment variables")
	}

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	DB = db
	fmt.Println("Successfully connected to database")
	return nil
}

// DBMigrate menjalankan database migration
func DBMigrate() {
	migrations := &migrate.EmbedFileSystemMigrationSource{
		FileSystem: dbMigrations,
		Root:       "sql_migrations",
	}

	n, err := migrate.Exec(DB, "postgres", migrations, migrate.Up)
	if err != nil {
		panic(err)
	}

	fmt.Println("Migration success, applied", n, "migrations!")
}
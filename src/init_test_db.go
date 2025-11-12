package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	dsn := "postgres://test_user:test_password@postgres:5432/test_db?sslmode=disable"

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatal("failed to connect to test database: ", err)
	}
	defer db.Close()

	if err := waitForDB(db); err != nil {
		log.Fatal("failed to wait for database: ", err)
	}

	if err := runSQLScriptsSequentially(db); err != nil {
		log.Fatal("failed to initialize database: ", err)
	}

	fmt.Println("Test database initialized successfully")
}

func runSQLScriptsSequentially(db *sqlx.DB) error {
	scripts := []string{
		"/app/internal/database/sql/01-create.sql",
		// "/app/internal/database/sql/02-constraints.sql",
		// "/app/internal/database/sql/03-inserts.sql",
		// "/app/internal/database/sql/trigger_accept.sql",
		// "/app/internal/database/sql/trigger_order.sql",
	}

	for _, script := range scripts {
		fmt.Printf("Running script: %s\n", script)
		content, err := os.ReadFile(script)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", script, err)
		}

		_, err = db.Exec(string(content))
		if err != nil {
			return fmt.Errorf("failed to execute %s: %w", script, err)
		}
		time.Sleep(100 * time.Millisecond)
	}
	return nil
}

func waitForDB(db *sqlx.DB) error {
	for i := 0; i < 30; i++ {
		err := db.Ping()
		if err == nil {
			return nil
		}
		fmt.Printf("Waiting for database... (attempt %d/30)\n", i+1)
		time.Sleep(time.Second)
	}
	return fmt.Errorf("database not ready after 30 attempts")
}

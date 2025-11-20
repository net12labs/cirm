package data

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type SqliteDb struct {
	db     *sql.DB
	Db     *sql.DB
	DbPath string
}

func NewDb() *SqliteDb {
	return &SqliteDb{}
}

func (db *SqliteDb) Init() error {
	var err error
	db.db, err = sql.Open("sqlite3", db.DbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	db.Db = db.db

	// Create test table
	createTableSQL := `CREATE TABLE IF NOT EXISTS test (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		value TEXT
	);`

	_, err = db.db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	// Insert sample data
	_, err = db.db.Exec("INSERT INTO test (name, value) VALUES (?, ?)", "sample", "test_value")
	if err != nil {
		return fmt.Errorf("failed to insert sample data: %w", err)
	}

	return nil
}

func (db *SqliteDb) Select(query string) error {
	rows, err := db.db.Query(query)
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("failed to get columns: %w", err)
	}

	fmt.Printf("Columns: %v\n", columns)

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return fmt.Errorf("failed to scan row: %w", err)
		}

		fmt.Printf("Row: %v\n", values)
	}

	return rows.Err()
}

func (db *SqliteDb) Close() error {
	if db.db != nil {
		return db.db.Close()
	}
	return nil
}

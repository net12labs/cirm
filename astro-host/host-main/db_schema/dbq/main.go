package dbq

import (
	"fmt"

	mdata "github.com/net12labs/cirm/mali/data"
)

type Database struct {
	Name   string
	Tables []*Table
	Db     *mdata.SqliteDb
}

type Table struct {
	Database *Database
	Name     string
	Columns  []*TableColumn
	Indexes  []*Index
}

type TableColumn struct {
	Name            string
	DataType        string
	Table           *Table
	IsPrimaryKey    bool
	IsAutoIncrement bool
	IsUnique        bool
	IsNullable      bool
	DefaultValue    string
	ForeignKey      *ForeignKey
}

type ForeignKey struct {
	ReferencedTable  string
	ReferencedColumn string
	OnDelete         string // CASCADE, SET NULL, RESTRICT, etc.
	OnUpdate         string
}

type Index struct {
	Name     string
	Table    *Table
	Columns  []string
	IsUnique bool
}

func NewDatabase(name string) *Database {
	return &Database{
		Name:   name,
		Tables: []*Table{},
	}
}

func (db *Database) AddTable(name string) *Table {
	table := &Table{
		Database: db,
		Name:     name,
		Columns:  []*TableColumn{},
		Indexes:  []*Index{},
	}
	db.Tables = append(db.Tables, table)
	return table
}

func (t *Table) AddColumn(name string, dataType string) *TableColumn {
	column := &TableColumn{
		Name:       name,
		DataType:   dataType,
		Table:      t,
		IsNullable: true, // default to nullable
	}
	t.Columns = append(t.Columns, column)
	return column
}

func (c *TableColumn) PrimaryKey() *TableColumn {
	c.IsPrimaryKey = true
	c.IsNullable = false
	return c
}

func (c *TableColumn) AutoIncrement() *TableColumn {
	c.IsAutoIncrement = true
	return c
}

func (c *TableColumn) Unique() *TableColumn {
	c.IsUnique = true
	return c
}

func (c *TableColumn) NotNull() *TableColumn {
	c.IsNullable = false
	return c
}

func (c *TableColumn) Default(value string) *TableColumn {
	c.DefaultValue = value
	return c
}

func (c *TableColumn) References(tableName, columnName string) *TableColumn {
	c.ForeignKey = &ForeignKey{
		ReferencedTable:  tableName,
		ReferencedColumn: columnName,
	}
	return c
}

func (c *TableColumn) OnDelete(action string) *TableColumn {
	if c.ForeignKey != nil {
		c.ForeignKey.OnDelete = action
	}
	return c
}

func (c *TableColumn) OnUpdate(action string) *TableColumn {
	if c.ForeignKey != nil {
		c.ForeignKey.OnUpdate = action
	}
	return c
}

func (t *Table) AddIndex(name string, columns []string, unique bool) *Index {
	index := &Index{
		Name:     name,
		Table:    t,
		Columns:  columns,
		IsUnique: unique,
	}
	t.Indexes = append(t.Indexes, index)
	return index
}
func (db *Database) GetTable(name string) *Table {
	for _, table := range db.Tables {
		if table.Name == name {
			return table
		}
	}
	return nil
}

func (t *Table) GetColumn(name string) *TableColumn {
	for _, column := range t.Columns {
		if column.Name == name {
			return column
		}
	}
	return nil
}

func (db *Database) ListTables() []string {
	tableNames := []string{}
	for _, table := range db.Tables {
		tableNames = append(tableNames, table.Name)
	}
	return tableNames
}

func (t *Table) ListColumns() []string {
	columnNames := []string{}
	for _, column := range t.Columns {
		columnNames = append(columnNames, column.Name)
	}
	return columnNames
}

func (t *Table) MakeCreateQuery() string {
	query := "CREATE TABLE IF NOT EXISTS " + t.Name + " (\n"

	// Add columns
	for i, column := range t.Columns {
		query += "  " + column.Name + " " + column.DataType

		if column.IsPrimaryKey {
			query += " PRIMARY KEY"
		}

		if column.IsAutoIncrement {
			query += " AUTOINCREMENT"
		}

		if !column.IsNullable {
			query += " NOT NULL"
		}

		if column.IsUnique && !column.IsPrimaryKey {
			query += " UNIQUE"
		}

		if column.DefaultValue != "" {
			query += " DEFAULT " + column.DefaultValue
		}

		if column.ForeignKey != nil {
			query += " REFERENCES " + column.ForeignKey.ReferencedTable +
				"(" + column.ForeignKey.ReferencedColumn + ")"
			if column.ForeignKey.OnDelete != "" {
				query += " ON DELETE " + column.ForeignKey.OnDelete
			}
			if column.ForeignKey.OnUpdate != "" {
				query += " ON UPDATE " + column.ForeignKey.OnUpdate
			}
		}

		if i < len(t.Columns)-1 {
			query += ","
		}
		query += "\n"
	}

	query += ");"
	return query
}

func (t *Table) MakeIndexQueries() []string {
	queries := []string{}
	for _, index := range t.Indexes {
		query := "CREATE "
		if index.IsUnique {
			query += "UNIQUE "
		}
		query += "INDEX " + index.Name + " ON " + t.Name + " ("
		for i, col := range index.Columns {
			query += col
			if i < len(index.Columns)-1 {
				query += ", "
			}
		}
		query += ");"
		queries = append(queries, query)
	}
	return queries
}

func (db *Database) InitDb() error {
	if db.Db == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	if db.Db.Db == nil {
		return fmt.Errorf("database connection Db field is nil")
	}

	tx, err := db.Db.Db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Create all tables
	for _, table := range db.Tables {
		createTableSQL := table.MakeCreateQuery()
		_, err := tx.Exec(createTableSQL)
		if err != nil {
			return fmt.Errorf("failed to create table %s: %w", table.Name, err)
		}

		// Create indexes for this table
		for _, indexQuery := range table.MakeIndexQueries() {
			_, err := tx.Exec(indexQuery)
			if err != nil {
				// Ignore error if index already exists
				// Could check specific error, but IF NOT EXISTS isn't standard for all DBs
				continue
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (db *Database) DropAllTables() error {
	if db.Db == nil || db.Db.Db == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	tx, err := db.Db.Db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Drop tables in reverse order to handle foreign key constraints
	for i := len(db.Tables) - 1; i >= 0; i-- {
		table := db.Tables[i]
		_, err := tx.Exec("DROP TABLE IF EXISTS " + table.Name)
		if err != nil {
			return fmt.Errorf("failed to drop table %s: %w", table.Name, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (db *Database) ResetDb() error {
	if err := db.DropAllTables(); err != nil {
		return fmt.Errorf("failed to drop tables: %w", err)
	}
	if err := db.InitDb(); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	return nil
}

func (db *Database) PrintCreateQueries() {
	for _, table := range db.Tables {
		fmt.Println(table.MakeCreateQuery())
		for _, indexQuery := range table.MakeIndexQueries() {
			fmt.Println(indexQuery)
		}
	}
}

package sitedb

import (
	"fmt"

	"github.com/net12labs/cirm/astro-host/host-main/db_schema/dbq"
)

var Dbo = &dbq.Database{
	Name:   "site",
	Tables: []*dbq.Table{},
}

func InitSchema() {
	// Account table
	accountTable := Dbo.AddTable("account")
	accountTable.AddColumn("id", "INTEGER").PrimaryKey().AutoIncrement()
	accountTable.AddColumn("username", "VARCHAR(255)").NotNull().Unique()
	accountTable.AddColumn("type", "VARCHAR(50)").NotNull().Default("'user'")
	accountTable.AddColumn("created_at", "DATETIME").NotNull()
	accountTable.AddColumn("updated_at", "DATETIME").NotNull()
	accountTable.AddIndex("idx_account_username", []string{"username"}, false)

	// User table
	userTable := Dbo.AddTable("user")
	userTable.AddColumn("id", "INTEGER").PrimaryKey().AutoIncrement()
	userTable.AddColumn("account_id", "INTEGER").NotNull().References("account", "id").OnDelete("CASCADE")
	userTable.AddColumn("username", "VARCHAR(255)").NotNull().Unique()
	userTable.AddColumn("created_at", "DATETIME").NotNull()
	userTable.AddColumn("updated_at", "DATETIME").NotNull()
	userTable.AddIndex("idx_user_account_id", []string{"account_id"}, false)
	userTable.AddIndex("idx_user_username", []string{"username"}, false)

	// Credential table
	credentialTable := Dbo.AddTable("credential")
	credentialTable.AddColumn("id", "INTEGER").PrimaryKey().AutoIncrement()
	credentialTable.AddColumn("user_id", "INTEGER").NotNull().References("user", "id").OnDelete("CASCADE")
	credentialTable.AddColumn("username", "VARCHAR(255)").NotNull()
	credentialTable.AddColumn("password", "VARCHAR(255)").NotNull()
	credentialTable.AddColumn("created_at", "DATETIME").NotNull()
	credentialTable.AddColumn("updated_at", "DATETIME").NotNull()
	credentialTable.AddIndex("idx_credential_user_id", []string{"user_id"}, false)
	credentialTable.AddIndex("idx_credential_username", []string{"username"}, false)

	// Session table
	sessionTable := Dbo.AddTable("session")
	sessionTable.AddColumn("id", "INTEGER").PrimaryKey().AutoIncrement()
	sessionTable.AddColumn("user_id", "INTEGER").NotNull().References("user", "id").OnDelete("CASCADE")
	sessionTable.AddColumn("token", "VARCHAR(255)").NotNull().Unique()
	sessionTable.AddColumn("expires_at", "DATETIME").NotNull()
	sessionTable.AddColumn("created_at", "DATETIME").NotNull()
	sessionTable.AddColumn("last_activity", "DATETIME").NotNull()
	sessionTable.AddIndex("idx_session_user_id", []string{"user_id"}, false)
	sessionTable.AddIndex("idx_session_token", []string{"token"}, true)
	sessionTable.AddIndex("idx_session_expires_at", []string{"expires_at"}, false)

	fmt.Printf("DATABASE SCHEMA INITIALIZED: %s\n", userTable.MakeCreateQuery())
}

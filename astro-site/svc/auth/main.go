package auth

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"

	args_auth "github.com/net12labs/cirm/astro-site/args/auth"
	odata "github.com/net12labs/cirm/ops/data"
)

type Auth struct {
	// Auth fields here
}

func NewService() *Auth {
	return &Auth{}
}

func (a *Auth) AccountCreate(args *args_auth.LoginCredentials) error {
	if !args.IsValid() {
		return fmt.Errorf("invalid account creation arguments")
	}

	tx, err := odata.Db("site").Db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Check if username already exists in account table
	var existingAccount int
	err = tx.QueryRow("SELECT COUNT(*) FROM account WHERE username = ?", args.Username).Scan(&existingAccount)
	if err != nil {
		return fmt.Errorf("failed to check existing account: %w", err)
	}
	if existingAccount > 0 {
		return fmt.Errorf("account with username '%s' already exists", args.Username)
	}

	// Check if username already exists in user table
	var existingUser int
	err = tx.QueryRow("SELECT COUNT(*) FROM user WHERE username = ?", args.Username).Scan(&existingUser)
	if err != nil {
		return fmt.Errorf("failed to check existing user: %w", err)
	}
	if existingUser > 0 {
		return fmt.Errorf("user with username '%s' already exists", args.Username)
	}

	now := time.Now().Format("2006-01-02 15:04:05")

	// Insert account and get account_id
	query := "INSERT INTO account (username, type, created_at, updated_at) VALUES (?, ?, ?, ?)"
	result, err := tx.Exec(query, args.Username, "user", now, now)
	if err != nil {
		return fmt.Errorf("failed to create account - query: %s, username: %s, error: %w", query, args.Username, err)
	}
	accountID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get account ID: %w", err)
	}

	// Insert user with account_id and get user_id
	result, err = tx.Exec("INSERT INTO user (account_id, username, created_at, updated_at) VALUES (?, ?, ?, ?)",
		accountID, args.Username, now, now)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	userID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get user ID: %w", err)
	}

	// Insert credential with user_id
	_, err = tx.Exec("INSERT INTO credential (user_id, username, password, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		userID, args.Username, args.Password, now, now)
	if err != nil {
		return fmt.Errorf("failed to create credential: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (a *Auth) UserLogin(cred *args_auth.LoginCredentials) (string, error) {
	if !cred.IsValid() {
		return "", fmt.Errorf("invalid login credentials")
	}

	tx, err := odata.Db("site").Db.Begin()
	if err != nil {
		return "", fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Verify credentials and get user_id
	var userID int
	err = tx.QueryRow("SELECT user_id FROM credential WHERE username = ? AND password = ?", cred.Username, cred.Password).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("invalid username or password")
		}
		return "", fmt.Errorf("failed to authenticate user: %w", err)
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	expiresAt := time.Now().Add(24 * time.Hour).Format("2006-01-02 15:04:05")

	// Check if user already has a valid session
	var existingToken string
	var sessionExpiresAt string
	err = tx.QueryRow("SELECT token, expires_at FROM session WHERE user_id = ? AND expires_at > ? ORDER BY expires_at DESC LIMIT 1",
		userID, now).Scan(&existingToken, &sessionExpiresAt)

	if err == nil {
		// Valid session exists, update last_activity and return existing token
		_, err = tx.Exec("UPDATE session SET last_activity = ? WHERE token = ?", now, existingToken)
		if err != nil {
			return "", fmt.Errorf("failed to update session activity: %w", err)
		}

		if err := tx.Commit(); err != nil {
			return "", fmt.Errorf("failed to commit transaction: %w", err)
		}

		return existingToken, nil
	} else if err != sql.ErrNoRows {
		return "", fmt.Errorf("failed to check existing session: %w", err)
	}

	// No valid session exists, create a new one
	token, err := generateToken()
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	_, err = tx.Exec("INSERT INTO session (user_id, token, expires_at, created_at, last_activity) VALUES (?, ?, ?, ?, ?)",
		userID, token, expiresAt, now, now)
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	return token, nil
}

func generateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// todo:
/*
-- create main ai agent instance for the user
-- create a profile for the user based on the role -- and then based on the page we go to - retrieve profile
-- ai agent should be attached to a profile
-- inirialize conversations for the user ai agent
-- add websockets to push agemt conversations out
-- then we add logic o create projects/missions etc. - this what ai will be filling in
*/

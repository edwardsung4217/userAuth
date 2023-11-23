// db.go
package types

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// OpenDB establishes a connection to the SQLite database
func OpenDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// SaveUser saves a user to the database
func SaveUser(db *sql.DB, user User) error {
	_, err := db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", user.Username, user.Password)
	if err != nil {
		log.Println("Error saving user:", err)
		return err
	}
	return nil
}

// GetUser retrieves a user from the database by ID
func GetUser(db *sql.DB, userID int) (User, error) {
	var user User
	err := db.QueryRow("SELECT * FROM users WHERE id = ?", userID).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		log.Println("Error retrieving user:", err)
		return User{}, err
	}
	return user, nil
}

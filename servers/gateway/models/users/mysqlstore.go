package users

import (
	"database/sql"
	"fmt"
	"os"

	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Create a MYSQL local database with password
var db, err = sql.Open(os.Getenv("MYSQL_DATABASE"), os.Getenv("DSN"))

//MainSQLStore represents an MySQL database store.
type MainSQLStore struct {
	database *sql.DB
}

//NewMainSQLStore constructs and returns a new MainSQLStore
func NewMainSQLStore(db *sql.DB) *MainSQLStore {
	return &MainSQLStore{
		database: db,
	}
}

// GetByID gets and returns the User from the MySQLStore database with the ID provided
func (mss *MainSQLStore) GetByID(id int64) (*User, error) {
	selectQuery := "select id,email,passHash,username,firstName,lastName,photoUrl from Users where id=?"
	row := mss.database.QueryRow(selectQuery, id)

	userInfo := &User{}
	if err := row.Scan(&userInfo.ID, &userInfo.Email, &userInfo.PassHash, &userInfo.UserName, &userInfo.FirstName,
		&userInfo.LastName, &userInfo.PhotoURL); err != nil {
		return nil, ErrUserNotFound
	}

	return userInfo, nil
}

// GetByEmail gets and returns the User from the MySQLStore database with the email provided
func (mss *MainSQLStore) GetByEmail(email string) (*User, error) {
	selectQuery := "select id,email,passHash,username,firstName,lastName,photoUrl from Users where email=?"
	row := mss.database.QueryRow(selectQuery, email)

	userInfo := &User{}
	if err := row.Scan(&userInfo.ID, &userInfo.Email, &userInfo.PassHash, &userInfo.UserName, &userInfo.FirstName,
		&userInfo.LastName, &userInfo.PhotoURL); err != nil {
		return nil, ErrUserNotFound
	}
	return userInfo, nil
}

// GetByUserName gets and returns the User from the MySQLStore database with the username provided
func (mss *MainSQLStore) GetByUserName(username string) (*User, error) {
	selectQuery := "select id,email,passHash,username,firstName,lastName,photoUrl from Users where username=?"
	row := mss.database.QueryRow(selectQuery, username)

	userInfo := &User{}
	if err := row.Scan(&userInfo.ID, &userInfo.Email, &userInfo.PassHash, &userInfo.UserName, &userInfo.FirstName,
		&userInfo.LastName, &userInfo.PhotoURL); err != nil {
		return nil, ErrUserNotFound
	}
	return userInfo, nil
}

// Insert values of the User into the MySQLStore and returns the same User
func (mss *MainSQLStore) Insert(user *User) (*User, error) {
	insq := "insert into Users(email,passHash,username,firstName,lastName,photoUrl) values(?,?,?,?,?,?)"
	_, err := mss.database.Exec(insq, user.Email, user.PassHash, user.UserName, user.FirstName, user.LastName, user.PhotoURL)
	if err != nil {
		fmt.Println(err.Error())
		return nil, fmt.Errorf("error inserting row: %v", err)
	}
	return user, nil
}

// Update FirstName and LastName of the User into the MySQLStore and returns the updated User
func (mss *MainSQLStore) Update(id int64, updates *Updates) (*User, error) {
	updateQuery := "update Users set firstName= ?, lastName= ? where id= ?"
	_, err := mss.database.Exec(updateQuery, updates.FirstName, updates.LastName, id)
	if err != nil {
		return nil, fmt.Errorf("error updating row: %v", err)
	}

	userInfo, _ := mss.GetByID(id)
	return userInfo, nil
}

// Delete the User from the MySQLStore database with the ID provided
func (mss *MainSQLStore) Delete(id int64) error {
	deleteQuery := "delete from Users where id= ?"
	_, err := mss.database.Exec(deleteQuery, id)
	if err != nil {
		return ErrUserNotFound
	}
	return nil
}

// LogSignIn logs all user sign-in attempts
func (mss *MainSQLStore) LogSignIn(user *User, dateTime time.Time, clientIP string) error {
	insertQuery := "insert into UserSignIn(UserID,LoginTime,clientIPAddress) values(?,?,?)"
	_, err := mss.database.Exec(insertQuery, user.ID, dateTime, clientIP)
	if err != nil {
		return fmt.Errorf("error logging User SignIn: %v", err)
	}
	return nil
}

package users

import (
	"errors"
	"time"
)

//FakeMySQLStore represents an MySQL database store.
type FakeMySQLStore struct {
	User *User
}

//GetByID returns the User with the given ID
func (s *FakeMySQLStore) GetByID(id int64) (*User, error) {
	if id == 0 {
		return &User{
			ID: 0,
		}, nil
	}
	return nil, errors.New("only id 0 can be found")
}

//GetByUserName returns the User with the given email
func (s *FakeMySQLStore) GetByUserName(email string) (*User, error) {
	if email == "username" {
		return &User{
			UserName: "username",
		}, nil
	}
	return nil, errors.New("only username \"username\" can be found")
}

//GetByEmail returns the User with the given username
func (s *FakeMySQLStore) GetByEmail(email string) (*User, error) {
	if email == "test@test.com" {
		return &User{}, nil
	}
	return nil, errors.New("only email \"test@test.com\" can be found")
}

// Insert values of the User into the FakeMySQLStore and returns the same User
func (s *FakeMySQLStore) Insert(user *User) (*User, error) {
	if user != nil {
		return &User{
			ID: 0,
		}, nil
	}
	return nil, errors.New("User could not be inserted")
}

// Update values of the User into the FakeMySQLStore and returns the same User
func (s *FakeMySQLStore) Update(id int64, updates *Updates) (*User, error) {
	if id == 0 {
		return &User{
			ID:        0,
			FirstName: updates.FirstName,
			LastName:  updates.LastName,
		}, nil
	}
	return nil, errors.New("User could not be updated")
}

// Delete the User from the FakeMySQLStore database with the ID provided
func (s *FakeMySQLStore) Delete(id int64) error {
	return nil
}

// LogSignIn logs all user sign-in attempts
func (s *FakeMySQLStore) LogSignIn(user *User, dateTime time.Time, clientIP string) error {
	return nil
}

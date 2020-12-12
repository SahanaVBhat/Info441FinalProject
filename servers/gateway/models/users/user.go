package users

import (
	"crypto/md5"
	"fmt"
	"net/mail"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

//gravatarBasePhotoURL is the base URL for Gravatar image requests.
//See https://id.gravatar.com/site/implement/images/ for details
const gravatarBasePhotoURL = "https://www.gravatar.com/avatar/"

//bcryptCost is the default bcrypt cost to use when hashing passwords
var bcryptCost = 13

//User represents a user account in the database
type User struct {
	ID        int64  `json:"id"`
	Email     string `json:"-"` //never JSON encoded/decoded
	PassHash  []byte `json:"-"` //never JSON encoded/decoded
	UserName  string `json:"userName"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	PhotoURL  string `json:"photoURL"`
}

//Credentials represents user sign-in credentials
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

//NewUser represents a new user signing up for an account
type NewUser struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	PasswordConf string `json:"passwordConf"`
	UserName     string `json:"userName"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
}

//Updates represents allowed updates to a user profile
type Updates struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

//Validate validates the new user and returns an error if
//any of the validation rules fail, or nil if its valid
func (nu *NewUser) Validate() error {
	e, err := mail.ParseAddress(nu.Email)
	if err != nil {
		return fmt.Errorf("%v is not a valid email", e)
	}

	if len(nu.Password) < 6 {
		return fmt.Errorf("Password must be at least 6 characters")
	}

	if strings.Compare(nu.Password, nu.PasswordConf) != 0 {
		return fmt.Errorf("Password and password confirmation do not match")
	}

	if len(nu.UserName) < 1 {
		return fmt.Errorf("Username must be non-zero length")
	}

	if strings.ContainsAny(nu.UserName, " ") {
		return fmt.Errorf("Username may not contain spaces")
	}

	return nil
}

//ToUser converts the NewUser to a User, setting the
//PhotoURL and PassHash fields appropriately
func (nu *NewUser) ToUser() (*User, error) {
	newUser := &User{}

	err := nu.Validate()
	if err != nil {
		return newUser, fmt.Errorf("Validation error")
	}

	newUser.ID = 0
	newUser.Email = nu.Email
	newUser.UserName = nu.UserName
	newUser.FirstName = nu.FirstName
	newUser.LastName = nu.LastName

	uniqueGravatar := md5.Sum([]byte(newUser.Email))
	newUser.PhotoURL = fmt.Sprintf("https://www.gravatar.com/avatar/%v.jpg", uniqueGravatar)

	//TODO: also call .SetPassword() to set the PassHash
	//field of the User to a hash of the NewUser.Password
	err2 := newUser.SetPassword(nu.Password)
	if err2 != nil {
		return nil, err2
	}
	return newUser, nil
}

//FullName returns the user's full name, in the form:
// "<FirstName> <LastName>"
//If either first or last name is an empty string, no
//space is put between the names. If both are missing,
//this returns an empty string
func (u *User) FullName() string {
	if len(u.FirstName) == 0 && len(u.LastName) == 0 {
		return ""
	} else if len(u.FirstName) == 0 {
		return u.LastName
	} else if len(u.LastName) == 0 {
		return u.FirstName
	} else {
		return u.FirstName + " " + u.LastName
	}

}

//SetPassword hashes the password and stores it in the PassHash field
func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 13)
	if err != nil {
		return fmt.Errorf("error generating bcrypt hash: %v\n", err)
	}
	u.PassHash = hash

	return nil
}

//Authenticate compares the plaintext password against the stored hash
//and returns an error if they don't match, or nil if they do
func (u *User) Authenticate(password string) error {
	return bcrypt.CompareHashAndPassword(u.PassHash, []byte(password))
}

//ApplyUpdates applies the updates to the user. An error
//is returned if the updates are invalid
func (u *User) ApplyUpdates(updates *Updates) error {
	update := &Updates{}
	if len(u.FirstName) == 0 || len(u.LastName) == 0 {
		return fmt.Errorf("Invalid first name or last name")
	}
	update.FirstName = u.FirstName
	update.LastName = u.LastName

	return nil
}

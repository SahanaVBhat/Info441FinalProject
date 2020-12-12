package users

import (
	"testing"
	"reflect"
)

// TestValidate tests the Validate function
func TestValidate(t *testing.T) {
	cases := []struct {
		hint        string
		testUser    *NewUser
		expectError bool
	}{
		{
			"NewUser is Valid",
			&NewUser{
				"test@test.com",
				"passhash123",
				"passhash123",
				"username",
				"firstname",
				"lastname",
			},
			false,
		},
		{
			"NewUser Not Valid. New User is Empty struct",
			&NewUser{},
			true,
		},
		{
			"Email is not valid",
			&NewUser{
				"test@test@com.gmail.com",
				"passhash123",
				"passhash123",
				"username",
				"firstname",
				"lastname",
			},
			true,
		},
		{
			"Password must be at least 6 characters",
			&NewUser{
				"test@test.com",
				"pass",
				"pass",
				"username",
				"firstname",
				"lastname",
			},
			true,
		},
		{
			"Password and PasswordConf must match",
			&NewUser{
				"test@test.com",
				"passhash123",
				"passhash123123",
				"username",
				"firstname",
				"lastname",
			},
			true,
		},
		{
			"UserName must be non-zero length",
			&NewUser{
				"test@test.com",
				"passhash123",
				"passhash123",
				"",
				"firstname",
				"lastname",
			},
			true,
		},
		{
			"UserName may not contain spaces",
			&NewUser{
				"test@test.com",
				"passhash123",
				"passhash123",
				"User name",
				"firstname",
				"lastname",
			},
			true,
		},
	}

	for _, c := range cases {
		err := c.testUser.Validate()
		if err != nil && c.expectError == false {
			t.Errorf("Unexpected error on valid New User. ERROR: %v ", err.Error())
		}
		if err == nil && c.expectError == true {
			t.Errorf("expected error but didn't get one. HINT: %v ", c.hint)
		}
	}
}

// TestToUser tests the function ToUser
func TestToUser(t *testing.T) {
	cases := []struct {
		hint        string
		testUser    *NewUser
		expectedUser *User
		expectError bool
	}{
		{
			"NewUser is Valid",
			&NewUser{
				"test@test.com",
				"passhash123",
				"passhash123",
				"username",
				"firstname",
				"lastname",
			},
			&User {
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"https://www.gravatar.com/avatar/f660ab912ec121d1b1e928a0bb4bc61b15f5ad44d5efdc4e1c92a25e99b8e44a.jpg",
			},
			false,
		},
		{
			"NewUser Email has capital and space.",
			&NewUser{
				" teSt@teSt.com ",
				"passhash123",
				"passhash123",
				"username",
				"firstname",
				"lastname",
			},
			&User {
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"https://www.gravatar.com/avatar/f660ab912ec121d1b1e928a0bb4bc61b15f5ad44d5efdc4e1c92a25e99b8e44a.jpg",
			},
			false,
		},
		{
			"Invalid user with invalid email",
			&NewUser{
				"test@test@com.gmail.com",
				"passhash123",
				"passhash123",
				"username",
				"firstname",
				"lastname",
			},
			&User {
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"photoURL",
			},
			true,
		},
	}
	for _, c := range cases {

		if c.expectError {
			_, err := c.testUser.ToUser()
			if err == nil {
				t.Errorf("Expected error but got %v instead. HINT: %v ", err, c.hint)
			}
		} else {
			usr, err := c.testUser.ToUser()
			if err != nil  {
				t.Errorf("Unexpected error on valid New User. HINT: %v ", c.hint)
			}

			if !reflect.DeepEqual(usr.FirstName, c.expectedUser.FirstName) {
				t.Errorf("Error, invalid match FirstName: %v", c.testUser)
			}

			if !reflect.DeepEqual(usr.LastName, c.expectedUser.LastName) {
				t.Errorf("Error, invalid match LastName: %v", c.testUser)
			}

			if !reflect.DeepEqual(usr.PhotoURL, c.expectedUser.PhotoURL) {
				t.Errorf("Error, invalid match photoURL: %v", c.testUser)
			}
			
		}

		
	}
}

// TestFullName tests the function FullName
func TestFullName(t *testing.T) {
	cases := []struct {
		hint             string
		testUser         *User
		expectedFullName string
	}{
		{
			"FirstName and LastName are set",
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			"firstname lastname",
		},
		{
			"FirstName or LastName, neither field set",
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"",
				"",
				"photourl",
			},
			"",
		},
		{
			"no FirstName found",
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"",
				"lastname",
				"photourl",
			},
			"lastname",
		},
		{
			"no LastName found",
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"",
				"photourl",
			},
			"firstname",
		},
	}
	for _, c := range cases {
		fullName := c.testUser.FullName()
		if fullName != c.expectedFullName {
			t.Errorf("Error, user fullname is not as expected . HINT: %v ", c.hint)
		}
	}
}
// TestSetPassword tests the function SetPassword
func TestSetPassword(t *testing.T) {
	cases := []struct {
		hint        string
		testUser    *User
		password string
		expectError bool
	}{
		{
			"Password is Valid",
			&User{
				1,
				"test@test.com",
				[]byte(""),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			"passhash123456",
			false,
		},
		{
			"Password is invalid",
			&User{
				1,
				"test@test.com",
				[]byte(""),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			"﷐﷑﷒﷓﷔",
			false,
		},
	}
	for _, c := range cases {
		if c.expectError {
			err := c.testUser.SetPassword(c.password)
			if err == nil {
				t.Errorf("expected error for %v but didn't get one. HINT: %v ", c.password, c.hint)
			}
			
		} else {
			SetPasswordErr := c.testUser.SetPassword(c.password)
			if SetPasswordErr != nil {
				t.Errorf("Unexpected error when setting password: %v", err)
			}
			
		}

	}

}

func TestAuthenticate(t *testing.T) {
	cases := []struct {
		hint        string
		testUser    *User
		testpassword string
		expectError bool
	}{
		{
			"Password Matches",
			&User{
				1,
				"test@test.com",
				[]byte("passhash123456"),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			"passhash123456",
			false,
		},
		{
			"Password Does Not Match",
			&User{
				1,
				"test@test.com",
				[]byte("passhash123456"),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			"diff",
			true,
		},
	}

	for _, c := range cases {
		if c.expectError {
			err := c.testUser.Authenticate(c.testpassword)
			if err == nil && c.expectError == true {
				t.Errorf("expected error for %v but didn't get one. HINT: %v ", c.testpassword, c.hint)
			}
			
		} else {
			SetPasswordErr := c.testUser.SetPassword(c.testpassword)
			if SetPasswordErr != nil {
				t.Errorf("Unexpected error when setting password: %v", err)
			}
			err := c.testUser.Authenticate(c.testpassword)
			if err != nil {
				t.Errorf("Unexpected error on authenticating password. HINT: %v ", c.hint)
			}
			
		}

		
	}

}

func TestApplyUpdates(t *testing.T) {
	cases := []struct {
		hint           string
		testUser       *User
		updateExpected bool
	}{
		{
			"FirstName is empty",
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"",
				"lastname",
				"photourl",
			},
			true,
		},
		{
			"LastName is empty",
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"",
				"photourl",
			},
			true,
		},
		{
			"FirstName & LastName are empty",
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"",
				"",
				"photourl",
			},
			true,
		},
		{
			"Updates are Valid",
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"Newfirstname",
				"Newlastname",
				"photourl",
			},
			false,
		},
	}
	for _, c := range cases {
		testupdate := &Updates{"", ""}
		err := c.testUser.ApplyUpdates(testupdate)
		if err != nil && c.updateExpected == false {
			t.Errorf("Unexpected error on valid updates. HINT: %v ", c.hint)
		}
		if err == nil && c.updateExpected == true {
			if c.testUser.FirstName != testupdate.FirstName {
				t.Errorf("Error: FirstName not updated")
			}
			if c.testUser.LastName != testupdate.LastName {
				t.Errorf("Error: LastName not updated")
			}
		}
	}
}

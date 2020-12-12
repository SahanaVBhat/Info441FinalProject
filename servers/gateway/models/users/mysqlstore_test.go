package users

import (
	"reflect"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"regexp"
)

// TestMainSQLStore tests the MainSQLStore object
func TestMainSQLStore(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("There was a problem opening a database connection: [%v]", err)
	}
	defer db.Close()

	store := NewMainSQLStore(db)

	row := mock.NewRows([]string{
		"ID",
		"Email",
		"PassHash",
		"UserName",
		"FirstName",
		"LastName",
		"PhotoURL"},
	).AddRow(
		1,
		"test@test.com",
		[]byte("passhash123"),
		"username",
		"firstname",
		"lastname",
		"photourl",
	)

	query := "select id,email,passHash,username,firstName,lastName,photoUrl from Users where id=?"
	mock.ExpectQuery(query).WithArgs(1).WillReturnRows(row)
	_, err = store.GetByID(1)
	if err != nil {
		t.Errorf("Problem with setting up SQLstore: %v", err)
	}

}

// TestGetByID is a test function for the MainSQLStore's GetByID
func TestGetByID(t *testing.T) {
	cases := []struct {
		name         string
		expectedUser *User
		idToGet      int64
		expectError  bool
	}{
		{
			"User Found",
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			1,
			false,
		},
		{
			"User Not Found",
			&User{},
			2,
			true,
		},
		{
			"User With Large ID Found",
			&User{
				1234567890,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			1234567890,
			false,
		},
	}

	for _, c := range cases {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("There was a problem opening a database connection: [%v]", err)
		}
		defer db.Close()

		mainSQLStore := &MainSQLStore{db}

		row := mock.NewRows([]string{
			"ID",
			"Email",
			"PassHash",
			"UserName",
			"FirstName",
			"LastName",
			"PhotoURL"},
		).AddRow(
			c.expectedUser.ID,
			c.expectedUser.Email,
			c.expectedUser.PassHash,
			c.expectedUser.UserName,
			c.expectedUser.FirstName,
			c.expectedUser.LastName,
			c.expectedUser.PhotoURL,
		)

		query := "select id,email,passHash,username,firstName,lastName,photoUrl from Users where id=?"

		if c.expectError {
			mock.ExpectQuery(query).WithArgs(c.idToGet).WillReturnError(ErrUserNotFound)

			user, err := mainSQLStore.GetByID(c.idToGet)
			if user != nil || err == nil {
				t.Errorf("Expected error [%v] but got [%v] instead", ErrUserNotFound, err)
			}
		} else {
			mock.ExpectQuery(query).WithArgs(c.idToGet).WillReturnRows(row)

			user, err := mainSQLStore.GetByID(c.idToGet)
			if err != nil {
				t.Errorf("Unexpected error on successful test [%s]: %v", c.name, err)
			}
			if !reflect.DeepEqual(user, c.expectedUser) {
				t.Errorf("Error, invalid match in test [%s]", c.name)
			}
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}

	}
}

// TestGetByEmail is a test function for the MainSQLStore's GetByEmail
func TestGetByEmail(t *testing.T) {
	cases := []struct {
		name         string
		expectedUser *User
		emailToGet      string
		expectError  bool
	}{
		{
			"User Found",
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			"test@test.com",
			false,
		},
		{
			"User Not Found",
			&User{},
			"",
			true,
		},
		{
			"User With long email Found",
			&User{
				1234567890,
				"abcdefghijklmnopqrstuvxyz1234567890@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			"abcdefghijklmnopqrstuvxyz1234567890@test.com",
			false,
		},
	}

	for _, c := range cases {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("There was a problem opening a database connection: [%v]", err)
		}
		defer db.Close()

		mainSQLStore := &MainSQLStore{db}

		row := mock.NewRows([]string{
			"ID",
			"Email",
			"PassHash",
			"UserName",
			"FirstName",
			"LastName",
			"PhotoURL"},
		).AddRow(
			c.expectedUser.ID,
			c.expectedUser.Email,
			c.expectedUser.PassHash,
			c.expectedUser.UserName,
			c.expectedUser.FirstName,
			c.expectedUser.LastName,
			c.expectedUser.PhotoURL,
		)

		query := "select id,email,passHash,username,firstName,lastName,photoUrl from Users where email=?"

		if c.expectError {
			mock.ExpectQuery(query).WithArgs(c.emailToGet).WillReturnError(ErrUserNotFound)

			user, err := mainSQLStore.GetByEmail(c.emailToGet)
			if user != nil || err == nil {
				t.Errorf("Expected error [%v] but got [%v] instead", ErrUserNotFound, err)
			}
		} else {
			mock.ExpectQuery(query).WithArgs(c.emailToGet).WillReturnRows(row)

			user, err := mainSQLStore.GetByEmail(c.emailToGet)
			if err != nil {
				t.Errorf("Unexpected error on successful test [%s]: %v", c.name, err)
			}
			if !reflect.DeepEqual(user, c.expectedUser) {
				t.Errorf("Error, invalid match in test [%s]", c.name)
			}
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}

	}
}

// TestGetByUserName is a test function for the MainSQLStore's GetByUserName
func TestGetByUserName(t *testing.T) {
	cases := []struct {
		name         string
		expectedUser *User
		userNameToGet      string
		expectError  bool
	}{
		{
			"User Found",
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			"username",
			false,
		},
		{
			"User Not Found",
			&User{},
			"",
			true,
		},
		{
			"User With weird username Found",
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"md4ufne*(#$90#*$7jd",
				"firstname",
				"lastname",
				"photourl",
			},
			"md4ufne*(#$90#*$7jd",
			false,
		},
	}

	for _, c := range cases {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("There was a problem opening a database connection: [%v]", err)
		}
		defer db.Close()

		mainSQLStore := &MainSQLStore{db}

		row := mock.NewRows([]string{
			"ID",
			"Email",
			"PassHash",
			"UserName",
			"FirstName",
			"LastName",
			"PhotoURL"},
		).AddRow(
			c.expectedUser.ID,
			c.expectedUser.Email,
			c.expectedUser.PassHash,
			c.expectedUser.UserName,
			c.expectedUser.FirstName,
			c.expectedUser.LastName,
			c.expectedUser.PhotoURL,
		)

		query := "select id,email,passHash,username,firstName,lastName,photoUrl from Users where username=?"

		if c.expectError {
			mock.ExpectQuery(query).WithArgs(c.userNameToGet).WillReturnError(ErrUserNotFound)

			user, err := mainSQLStore.GetByUserName(c.userNameToGet)
			if user != nil || err == nil {
				t.Errorf("Expected error [%v] but got [%v] instead", ErrUserNotFound, err)
			}
		} else {
			mock.ExpectQuery(query).WithArgs(c.userNameToGet).WillReturnRows(row)

			user, err := mainSQLStore.GetByUserName(c.userNameToGet)
			if err != nil {
				t.Errorf("Unexpected error on successful test [%s]: %v", c.name, err)
			}
			if !reflect.DeepEqual(user, c.expectedUser) {
				t.Errorf("Error, invalid match in test [%s]", c.name)
			}
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}

	}
}

// TestInsert is a test function for the MainSQLStore's Insert
func TestInsert(t *testing.T) {
	cases := []struct {
		name         string
		expectedUser *User
		userToInsert     *User
		expectError  bool
	}{
		{
			"User Inserted",
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			false,
		},
		{
			"Empty User",
			&User{},
			&User{},
			false,
		},
		{
			"User With longer data inserted",
			&User{
				1,
				"fgsdvbx@test.com",
				[]byte("dejifh9o4859rf"),
				"md4ufne*(#$90#*$7jd",
				"fname",
				"lname",
				"https://www.gravatar.com/avatar/",
			},
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			false,
		},
		{
			"User With Invalid items inserted",
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			&User{
				1,
				"wrong@wrong@wrong.com",
				[]byte("passhash123"),
				"username",
				"first name",
				"last name",
				"photourl",
			},
			true,
		},
	}

	for _, c := range cases {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("There was a problem opening a database connection: [%v]", err)
		}
		defer db.Close()

		mainSQLStore := &MainSQLStore{db}

		insert := regexp.QuoteMeta(`insert into user(email,passHash,username,firstName,lastName,photoUrl) values(?,?,?,?,?,?)`)

		if c.expectError {
			mock.ExpectExec(insert).WithArgs(c.userToInsert.Email, c.userToInsert.PassHash,c.userToInsert.UserName,
				c.userToInsert.FirstName,c.userToInsert.LastName, c.userToInsert.PhotoURL).WillReturnError(ErrUserNotFound)

			user, err := mainSQLStore.Insert(c.userToInsert)
			if user != nil || err == nil {
				t.Errorf("Expected error [%v] but got [%v] instead", ErrUserNotFound, err)
			}
		} else {
			mock.ExpectExec(insert).WithArgs(c.userToInsert.Email, c.userToInsert.PassHash,c.userToInsert.UserName,
				c.userToInsert.FirstName,c.userToInsert.LastName, c.userToInsert.PhotoURL).WillReturnResult(sqlmock.NewResult(1, 1))

			user, err := mainSQLStore.Insert(c.userToInsert)
			if err != nil {
				t.Errorf("Unexpected error on successful test [%s]: %v", c.name, err)
			}
			if !reflect.DeepEqual(user, c.userToInsert) {
				t.Errorf("Error, invalid match in test [%s]", c.name)
			}
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}

	}
}

// TestUpdate is a test function for the MainSQLStore's Update
func TestUpdate(t *testing.T) {
	cases := []struct {
		name         string
		userAfterUpdate *User
		userBeforeUpdate *User
		updates *Updates
		idToUpdate int64
		expectError  bool
	}{
		{
			"User Updated",
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"newfirstname",
				"newlastname",
				"photourl",
			},
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			&Updates{
				"newfirstname",
				"newlastname",
			},
			1,
			false,
		},
		{
			"Empty Updates",
			&User{},
			&User{},
			&Updates{},
			1,
			true,
		},
		{
			"Update Just First Name",
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			&Updates{
				"newfirstname",
				"",
			},
			1,
			true,
		},
		{
			"Update With Nonexistent ID",
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"newfirstname",
				"newlastname",
				"photourl",
			},
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			&Updates{
				"newfirstname",
				"newlastname",
			},
			100,
			true,
		},
		
	}

	for _, c := range cases {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("There was a problem opening a database connection: [%v]", err)
		}
		defer db.Close()

		mainSQLStore := &MainSQLStore{db}
		row := mock.NewRows([]string{
			"ID",
			"Email",
			"PassHash",
			"UserName",
			"FirstName",
			"LastName",
			"PhotoURL"},
		).AddRow(
			c.userAfterUpdate.ID,
			c.userAfterUpdate.Email,
			c.userAfterUpdate.PassHash,
			c.userAfterUpdate.UserName,
			c.userAfterUpdate.FirstName,
			c.userAfterUpdate.LastName,
			c.userAfterUpdate.PhotoURL,
		)

		insert := regexp.QuoteMeta(`insert into user(email,passHash,username,firstName,lastName,photoUrl) values(?,?,?,?,?,?)`)
		mock.ExpectExec(insert).WithArgs(c.userBeforeUpdate.Email, c.userBeforeUpdate.PassHash,c.userBeforeUpdate.UserName,
			c.userBeforeUpdate.FirstName,c.userBeforeUpdate.LastName, c.userBeforeUpdate.PhotoURL).WillReturnResult(sqlmock.NewResult(1, 1))

		_, err2 := mainSQLStore.Insert(c.userBeforeUpdate)
		if err2 != nil {
			t.Errorf("Unexpected error on inserting user to update %v", c.userBeforeUpdate)
		}

		updateQuery := regexp.QuoteMeta(`update user set firstName= ?, lastName= ? where id= ?`)
		selectQuery := "select id,email,passHash,username,firstName,lastName,photoUrl from Users where id=?"

		if c.expectError {
			mock.ExpectExec(updateQuery).WithArgs(c.updates.FirstName, c.updates.LastName, c.idToUpdate).WillReturnError(ErrUserNotFound)

			user, err := mainSQLStore.Update(c.idToUpdate, c.updates)
			if user != nil || err == nil {
				t.Errorf("Expected error [%v] but got [%v] instead", ErrUserNotFound, err)
			}
		} else {
			mock.ExpectExec(updateQuery).WithArgs(c.updates.FirstName, c.updates.LastName, c.idToUpdate).WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectQuery(selectQuery).WithArgs(c.idToUpdate).WillReturnRows(row)

			user, err := mainSQLStore.Update(c.idToUpdate, c.updates)
			if err != nil {
				t.Errorf("Unexpected error on successful test [%s]: %v", c.name, err)
			}
			if !reflect.DeepEqual(user, c.userAfterUpdate) {
				t.Errorf("Error, invalid match in test [%s]", c.name)
			}
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}

	}
}

// TestDelete is a test function for the MainSQLStore's Delete
func TestDelete(t *testing.T) {
	cases := []struct {
		name         string
		userToInsert *User
		idToDelete      int64
		expectError  bool
	}{
		{
			"User Found",
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			1,
			false,
		},
		{
			"User Not Found",
			&User{},
			2,
			true,
		},
		{
			"User With Large ID Deleted",
			&User{
				1234567890,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			1234567890,
			false,
		},
	}

	for _, c := range cases {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("There was a problem opening a database connection: [%v]", err)
		}
		defer db.Close()

		mainSQLStore := &MainSQLStore{db}

		insert := regexp.QuoteMeta(`insert into user(email,passHash,username,firstName,lastName,photoUrl) values(?,?,?,?,?,?)`)
		mock.ExpectExec(insert).WithArgs(c.userToInsert.Email, c.userToInsert.PassHash,c.userToInsert.UserName,
			c.userToInsert.FirstName,c.userToInsert.LastName, c.userToInsert.PhotoURL).WillReturnResult(sqlmock.NewResult(1, 1))

		_, err2 := mainSQLStore.Insert(c.userToInsert)
		if err2 != nil {
			t.Errorf("Unexpected error on inserting user to delete %v", err2)
		}

		query := "delete from users where id= ?"

		if c.expectError {
			mock.ExpectExec(query).WithArgs(c.idToDelete).WillReturnError(ErrUserNotFound)

			err := mainSQLStore.Delete(c.idToDelete)
			if err == nil {
				t.Errorf("Expected error [%v] but got [%v] instead", ErrUserNotFound, err)
			}
		} else {
			mock.ExpectExec(query).WithArgs(c.idToDelete).WillReturnResult(sqlmock.NewResult(c.idToDelete, 1))

			err := mainSQLStore.Delete(c.idToDelete)
			if err != nil {
				t.Errorf("Unexpected error on successful test [%s]: %v", c.name, err)
			}

			attemptSelect := "select id,email,passHash,username,firstName,lastName,photoUrl from Users where id=?"
			mock.ExpectQuery(attemptSelect).WithArgs(c.idToDelete).WillReturnError(ErrUserNotFound)
			res, err2 := mainSQLStore.GetByID(c.idToDelete)

			if res != nil || err2 == nil {
				t.Errorf("User with ID %v not deleted properly", c.idToDelete)
			}
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}

	}
}
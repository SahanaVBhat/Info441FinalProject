CREATE TABLE IF NOT EXISTS Users(
    ID INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	Email VARCHAR(128) NOT NULL UNIQUE,
	PassHash  VARCHAR(255) NOT NULL,
	UserName  VARCHAR(255) NOT NULL UNIQUE,
	FirstName VARCHAR(64) NOT NULL,
	LastName  VARCHAR(64) NOT NULL,
	PhotoURL  VARCHAR(128) NOT NULL
);

CREATE TABLE IF NOT EXISTS UserSignIn(
	UserSignInID INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	UserID INT NOT NULL,
	LoginTime DATETIME NOT NULL,
	ClientIPAddress VARCHAR(128) NOT NULL
);


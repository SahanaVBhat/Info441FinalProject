# Project Proposal

## Project Description 
<p> 
Currently there is no resource for UW students to see specific reviews and suggestions about a particular course. 
  Similar services, such as RateMyProfessor focuses on the professors themselves, while UW course evaluation only
  displays numerical ratings without much explanation for the course structure, details, or helpful suggestions. 
  As UW students and developers, we want to solve this issue by building an application to help UW students make 
  more informed decisions about the classes they want to take: by creating an application where they are able to 
  see reviews for a particular course and view the difficulty of the course, see the recommendations from other 
  students, and able to submit reviews/evaluations of a course anonymously.
</p>

## Technical Description
### Architectural diagram mapping 

User Client \
Gateway layer \
Authentication layer - verify student is UW student \
Redis session store \
mySQL database - information of all the evaluations/ratings of a course, information of students \
Search microservice - search for evaluations based on major/number \
Evaluation submission microservice - submit a form for evaluation of class

![](./InfrastructureDiagram.png)
https://lucid.app/invitations/accept/ceac9f69-1f24-4728-9327-78fdd456e2b0

### User Stories 

| Priority | User | Description | Technical Implementation Strategy |
| --- | --- | --- | --- |
| P0 | As a student | I want to be able to filter search for the course evaluations for a particular course submitted by other students | Access the MySQL database, return a response body that will be displayed in the front-end |
| P0 | As a student | I want to be able to post a evaluation of a course I took | Insert into the MySQL database, provide a request to the back-end |
| P0 | As a student | I want to be able to Sign in and Sign up as user | Store log in credentials in database, authenticate user and store sessions inside Redis database. |
| P1 | As a student | I want to be able to update my evaluation | User clicks on edit button, updates the fields,  save evaluation updates in SQL database |
| P1 | As a student | I want to be able to update the information in my profile after sign up | Update the MySQL database, return a response displayed on the front-end |

## API EndPoints
### User 
v1/User \
GET : to get a specific user account based on User ID. 
 * 200: Successful response with user information 
 * 401: Cannot verify User ID .
 * 415: Cannot decode body or receive unsupported body.
 * 500: Internal server error.

POST: to create new user account, Sign up
 * 201: created new user account
 * 400: bad request if user account already exists
 * 401: unauthorized user information access
 * 500: internal server error
 
PATCH: to update user email Id or password
 * 201: updated user information
 * 400: bad request if user account already exists
 * 401: unauthorized user information access
 * 500: internal server error

### Session
/v1/sessions \
POST
 * 201: created a new session for a user on sign in
 * 415: unsupported media
 * 500: internal server error
/v1/sessions/
GET
 * 400: bad request
 * 403: forbidden request if not user's session
DELETE
 * 400: bad request
 * 403: forbidden request if not user's session

### Course 
v1/Course :  to handle Course Information \
GET : to get a specific course based on Course ID. 
 * 200: Successful response with course information 
 * 401: Cannot verify Course ID .
 * 415: Cannot decode body or receive unsupported body.
 * 500: Internal server error.

POST: to add new Course for discussion
 * 201: created new Course Information
 * 400: bad request if course already exists
 * 401: unauthorized user information access
 * 500: internal server error

### Evaluation
v1/evaluation \
POST: to add evaluation to a course
 * 201: created new evaluation  
 * 401: unauthorized user information access
 * 500: internal server error
 
PATCH: to update evaluation text description, rating etc
201: updated evaluation parameters
* 401: unauthorized user information access
* 500: internal server error
 
DELETE: students who posted the evaluations will be able to delete their evaluations
* 200: Successfully deletes student evaluations.
* 401: Cannot verify review ID
* 500: Internal server error.

## Data Models
 * **Evaluation**: The courses are evaluated and reviewed by students. The evaluations are specific to a course. A single course can have multiple evaluations based on the quarter, year and instructors involved. One user can submit only one evaluation per course. User will be able to edit and delete evaluations.

	Evaluation: { \
		ID: int \
		studentID: int \
	courseID: int \
		instructors: string (multiple instructors can be included separated by ‘ , ’ or ‘ ; ‘) \
		year: date \
		quarter : string \
		creditType: string \
		credits: int \
		workload: int  (on a scale of 1-10) \
		gradingTechniques: int \
		description: string	\
		LikedUsers: string[] ( list of unique userIds who have liked the evaluation. The length of the array provides count of likes. Need to check if the userId exists before adding to the list, to maintain unique userIds) \
		DislikedUsers: string[] ( list of users who disliked the course) \
	}

 * **Course**: The evaluations are specific to a course. This data model contains course information such as Course Title, Course Code, Credits and a brief Description text.

	Course: { \
		ID: int \
		Title: string \
		Code: string \
		Description: string \
		Credits: int \
	}

 * **Student**: represents student information

	Student: { \
		ID :  int64 \
		Email:  string \
		PassHash: []byte \
		UserName: string \
		FirstName: string \
		LastName:  string \
		AvatarPhotoURL:  string \
	}

 * **NewStudent**: represents new student information, used during user registration

	NewStudent struct { \
		Email: string \
		Password: string \
		PasswordConfirmed: string \
		UserName: string \
		FirstName: string \
		LastName: string \
	}

 * **Student Credentials**: used during user login \
	Credentials: { \
		Email: string \ 
		Password: string \ 
	}






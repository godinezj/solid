package models

import (
	"database/sql"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID                uuid.UUID `json:"id" db:"id"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
	Email             string    `json:"email" db:"email"`
	PasswordHash      string    `json:"-" db:"password_hash"`
	Password          string    `json:"password" db:"-"`
	PasswordConfirm   string    `json:"password_confirm" db:"-"`
	ResetToken        uuid.UUID `json:"-" db:"reset_token"`
	ResetTokenConfirm uuid.UUID `json:"reset_token_confirm" db:"-"`
	// ResetTokenExpire time.Time `json:"-" db:"reset_token_expire"`
}

// String is not required by pop and may be deleted
func (u User) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Users is not required by pop and may be deleted
type Users []User

// String is not required by pop and may be deleted
func (u Users) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Create validates and creates a new User.
func (u *User) Create(tx *pop.Connection) (*validate.Errors, error) {
	u.Email = strings.ToLower(u.Email)
	pwdHash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return validate.NewErrors(), errors.WithStack(err)
	}
	u.PasswordHash = string(pwdHash)
	return tx.ValidateAndCreate(u)
}

// Update validates and Updates a new User.
func (u *User) Update(tx *pop.Connection) (*validate.Errors, error) {
	u.Email = strings.ToLower(u.Email)
	pwdHash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return validate.NewErrors(), errors.WithStack(err)
	}
	u.PasswordHash = string(pwdHash)
	return tx.ValidateAndUpdate(u)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (u *User) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: u.Email, Name: "Email"},
		&validators.EmailIsPresent{Field: u.Email, Name: "Email"},
		&validators.StringIsPresent{Field: u.Password, Name: "Password"},
		&validators.StringIsPresent{Field: u.Password, Name: "PasswordConfirm"},
		&validators.StringsMatch{Name: "Password", Field: u.Password, Field2: u.PasswordConfirm, Message: "Passwords do not match."},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (u *User) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&EmailNotTaken{Name: "Email", Field: u.Email, tx: tx},
	), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (u *User) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

type EmailNotTaken struct {
	Name  string
	Field string
	tx    *pop.Connection
}

func (v *EmailNotTaken) IsValid(errors *validate.Errors) {
	query := v.tx.Where("email=?", v.Field).Select("email")
	queryUser := User{}
	err := query.First(&queryUser)
	if err == nil { // found user with same email
		errors.Add(validators.GenerateKey(v.Name), "Account with that email aready exists")
	}
}

func (u *User) Load(tx *pop.Connection) error {
	// find the user by email
	err := tx.Where("email = ?", strings.ToLower(u.Email)).First(u)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			// couldn't find an user with that email address
			return errors.New("User not found.")
		}
		return errors.WithStack(err)
	}
	return nil
}

// Authorize checks user's password for logging in
func (u *User) Authorize(tx *pop.Connection) error {
	log.Println("Authenticating " + u.Email)
	if err := u.Load(tx); err != nil {
		return err
	}
	// confirm that the given password matches the hashed password from the db
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(u.Password))
	if err != nil {
		return errors.New("Wrong password.")
	}
	return nil
}

func (u *User) SendResetToken(tx *pop.Connection) error {
	if err := u.Load(tx); err != nil {
		return err
	}

	// set reset token
	log.Println("Setting reset token for " + u.Email)
	token, err := uuid.NewV4()
	if err == nil {
		// save token to user in db
		u.ResetToken = token
		tx.Update(u)
	}

	// TODO email token to user

	return nil
}

func (u *User) ChangePassword(tx *pop.Connection) (*validate.Errors, error) {
	// find user by email
	query := tx.Where("email = ?", u.Email)
	queryUser := User{}
	err := query.First(&queryUser)
	if err != nil {
		return nil, errors.New("User not found")
	}
	log.Printf("%s == %s\n", u.ResetTokenConfirm, queryUser.ResetToken)
	uuidsMatch := uuid.Equal(u.ResetTokenConfirm, queryUser.ResetToken)
	if !uuidsMatch {
		return nil, errors.New("Tokens do not match")
	}
	queryUser.Password = u.Password
	queryUser.PasswordConfirm = u.PasswordConfirm
	return queryUser.Update(tx)
}

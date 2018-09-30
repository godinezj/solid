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
	ID              uuid.UUID `json:"id" db:"id"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
	Email           string    `json:"email" db:"email"`
	PasswordHash    string    `json:"-" db:"password_hash"`
	Password        string    `json:"password" db:"-"`
	PasswordConfirm string    `json:"password_confirm" db:"-"`
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

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (u *User) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: u.Email, Name: "Email"},
		&validators.EmailIsPresent{Field: u.Email, Name: "Email"},
		&EmailNotTaken{Name: "Email", Field: u.Email, tx: tx},
		&validators.StringIsPresent{Field: u.Password, Name: "Password"},
		&validators.StringIsPresent{Field: u.Password, Name: "PasswordConfirm"},
		&validators.StringsMatch{Name: "Password", Field: u.Password, Field2: u.PasswordConfirm, Message: "Passwords do not match."},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (u *User) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(), nil
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

// Authorize checks user's password for logging in
func (u *User) Authorize(tx *pop.Connection) error {
	log.Println(u.Email)
	err := tx.Where("email = ?", strings.ToLower(u.Email)).First(u)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			// couldn't find an user with that email address
			return errors.New("User not found.")
		}
		return errors.WithStack(err)
	}
	// confirm that the given password matches the hashed password from the db
	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(u.Password))
	if err != nil {
		return errors.New("Wrong password.")
	}
	return nil
}

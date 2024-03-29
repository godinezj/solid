package models

import (
	"database/sql"
	"encoding/json"
	"strings"
	"time"
	"unicode"

	"bitbucket.org/godinezj/solid/ldap"
	"bitbucket.org/godinezj/solid/log"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

// User represents a user
type User struct {
	ID                uuid.UUID `json:"id" db:"id"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
	Username          string    `json:"username" db:"username"`
	FirstName         string    `json:"first_name" db:"first_name"`
	LastName          string    `json:"last_name" db:"last_name"`
	Email             string    `json:"email" db:"email"`
	Zip               string    `json:"zip" db:"zip"`
	Password          string    `json:"password" db:"-"`
	PasswordConfirm   string    `json:"password_confirm" db:"-"`
	ResetToken        uuid.UUID `json:"-" db:"reset_token"`
	ResetTokenConfirm uuid.UUID `json:"reset_token_confirm" db:"-"`
	ResetTokenExpire  time.Time `json:"-" db:"reset_token_expire"`
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
	u.Username = strings.ToLower(u.Username)
	verrs, err := tx.ValidateAndCreate(u)
	if err != nil || verrs.HasAny() {
		log.Errorf("verrs %v", verrs)
		log.Errorf("errs %v", err)
		return verrs, err
	}
	log.Info("User created in db")
	// only create ldap users in prod
	// var ENV = envy.Get("GO_ENV", "development")
	// if ENV != "production" {
	// 	log.Infof("ENV: %s, not creating users in LDAP", ENV)
	// 	return verrs, err
	// }
	// make admin connection
	client := ldap.Client{}
	defer client.Close() // close the admin connection
	client.Connect()
	if err != nil {
		log.Error("LDAP: An error occured creating user")
		return verrs, err
	}
	err = client.AdminAuth()
	if err != nil {
		verrs.Add("ldap", "An error occured creating user")
		return verrs, err
	}

	// add user
	_, err = client.AddUser(u.FirstName, u.LastName, u.Username, u.Password)
	if err != nil {
		verrs.Add("ldap", "An error occured creating user")
		return verrs, err
	}
	return verrs, err
}

// Update validates and Updates a new User.
func (u *User) Update(tx *pop.Connection) (*validate.Errors, error) {
	u.Email = strings.ToLower(u.Email)
	// TODO add ldap chpw functionality
	return tx.ValidateAndUpdate(u)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (u *User) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.RegexMatch{Name: "Username", Field: u.Username, Expr: "[a-zA-Z]+"},
		&validators.StringLengthInRange{Name: "Username", Field: u.Username, Min: 6, Max: 64},
		&validators.StringIsPresent{Field: u.Email, Name: "Email"},
		&validators.EmailIsPresent{Field: u.Email, Name: "Email"},
		&validators.StringIsPresent{Field: u.Password, Name: "Password"},
		&validators.StringIsPresent{Field: u.Password, Name: "PasswordConfirm"},
		&StrongPassword{Field: u.Password, Name: "Password"},
		&validators.StringLengthInRange{Field: u.Password, Name: "Password", Min: 6, Max: 64},
		&validators.StringsMatch{Name: "Password", Field: u.Password, Field2: u.PasswordConfirm, Message: "Passwords do not match."},
		&validators.RegexMatch{Name: "Zip", Field: u.Zip, Expr: "[0-9]+"},
		&validators.StringLengthInRange{Name: "Zip", Field: u.Zip, Min: 5, Max: 5},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
func (u *User) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&NotTaken{Name: "Email", Field: u.Email, tx: tx},
		&NotTaken{Name: "Username", Field: u.Username, tx: tx},
	), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (u *User) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

type StrongPassword struct {
	Name  string
	Field string
}

func (v *StrongPassword) IsValid(errors *validate.Errors) {
	mustHave := []func(rune) bool{
		unicode.IsUpper,
		unicode.IsLower,
		// unicode.IsPunct,
		unicode.IsDigit,
	}

	for _, testRune := range mustHave {
		found := false
		for _, r := range v.Field {
			if testRune(r) {
				found = true
			}
		}
		if !found {
			errors.Add(validators.GenerateKey(v.Name), "Password requires uppercase, lowercase, and a digit.")
			return
		}
	}
}

type NotTaken struct {
	Name  string
	Field string
	tx    *pop.Connection
}

func (v *NotTaken) IsValid(errors *validate.Errors) {
	query := v.tx.Where(v.Name+"=?", v.Field).Select(v.Name)
	queryUser := User{}
	err := query.First(&queryUser)
	if err == nil { // found user with same email
		errors.Add(validators.GenerateKey(v.Name), "Account with that "+v.Name+" aready exists")
	}
}

func (u *User) Load(tx *pop.Connection) error {
	// find the user by email
	err := tx.Where("email = ?", strings.ToLower(u.Email)).First(u)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			// couldn't find an user with that email address
			return errors.New("user not found")
		}
		return errors.WithStack(err)
	}
	return nil
}

// Authenticate checks user's password for logging in
func (u *User) Authenticate(tx *pop.Connection) error {
	log.Info("Authenticating " + u.Username)
	if err := u.Load(tx); err != nil {
		return err
	}

	ldap := ldap.Client{}
	if err := ldap.Connect(); err != nil {
		return err
	}
	err := ldap.Authenticate(u.Username, u.Password)
	if err != nil {
		return err
	}
	return nil
}

// SendResetToken sends the reset token via email
func (u *User) SendResetToken(tx *pop.Connection) error {
	if err := u.Load(tx); err != nil {
		return err
	}

	// set reset token
	log.Info("Setting reset token for " + u.Email)
	token, err := uuid.NewV4()
	if err == nil {
		// save token to user in db
		u.ResetToken = token
		tx.Update(u)
	}

	// TODO email token to user

	return nil
}

// ChangePassword used to change a user's password
func (u *User) ChangePassword(tx *pop.Connection) (*validate.Errors, error) {
	// find user by email
	query := tx.Where("email = ?", u.Email)
	queryUser := User{}
	err := query.First(&queryUser)
	if err != nil {
		return nil, errors.New("User not found")
	}
	log.Infof("%s == %s\n", u.ResetTokenConfirm, queryUser.ResetToken)
	if u.ResetTokenConfirm != queryUser.ResetToken {
		return nil, errors.New("Tokens do not match")
	}
	queryUser.Password = u.Password
	queryUser.PasswordConfirm = u.PasswordConfirm
	return queryUser.Update(tx)
}

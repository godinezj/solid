package actions

import (
	"log"

	"bitbucket.org/godinezj/solid/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
)

// UsersResource is the resource for the User model
type UsersResource struct {
	buffalo.Resource
}

// Create adds a User to the DB. This function is mapped to the
// path POST /users
func Create(c buffalo.Context) error {
	// Allocate an empty User
	user := &models.User{}

	// Bind user to the html form elements
	errMessage := "Could not create user."
	if err := c.Bind(user); err != nil {
		return c.Render(422, r.JSON(map[string]string{"message": errMessage}))
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		log.Println("No transaction found")
		return c.Render(422, r.JSON(map[string]string{"message": errMessage}))
	}

	// Validate the data
	verrs, err := user.Create(tx)
	if err != nil {
		log.Println("Validation failed")

		return c.Render(422, r.JSON(map[string]string{"message": errMessage}))
	}

	if verrs.HasAny() {
		// Respond with errors to user
		log.Println("Validation failed")
		for verr := range verrs.Errors {
			log.Println("Failed validation: " + verr)
		}
		return c.Render(422, r.JSON(verrs))
	}

	// render success message to user
	return c.Render(200, r.JSON(map[string]string{"message": "Account created, please login."}))
}

// Login logs in a user.
func Login(c buffalo.Context) error {
	user := &models.User{}
	// Bind the user to the html form elements
	errMessage := "Invalid email or password."
	if err := c.Bind(user); err != nil {
		log.Println(err)
		r.JSON(map[string]string{"message": errMessage})
	}
	tx := c.Value("tx").(*pop.Connection)
	err := user.Authenticate(tx)
	if err != nil {
		log.Println(err)
		// Return invalid email/password wether users exists or not
		return c.Render(422, r.JSON(map[string]string{"message": errMessage}))
	} else {
		c.Session().Set("current_user_id", user.ID)
		return c.Render(201, r.JSON(map[string]string{"message": "Login success."}))
	}
}

func GenPassResetToken(c buffalo.Context) error {
	user := &models.User{}
	// Bind the user to the html form elements
	if err := c.Bind(user); err != nil {
		log.Println(err)
		return c.Render(422, r.JSON(map[string]string{"message": "Could not reset password"}))
	}
	tx := c.Value("tx").(*pop.Connection)
	err := user.SendResetToken(tx)
	if err != nil {
		log.Println(err)
	}

	// reply successfully even if email/user does not exist
	return c.Render(200, r.JSON(map[string]string{"message": "Check your email."}))
}

func ValidatePassResetToken(c buffalo.Context) error {
	user := &models.User{}
	// Bind the user to the html form elements
	if err := c.Bind(user); err != nil {
		return c.Render(422, r.JSON(map[string]string{"message": "Could not validate password."}))
	}
	tx := c.Value("tx").(*pop.Connection)
	verrs, err := user.ChangePassword(tx)
	if err != nil {
		log.Println(err)
	}
	if verrs.HasAny() {
		return c.Render(422, r.JSON(verrs))
	}

	// requires email, token, password, & password_confirmation
	return c.Render(200, r.JSON(map[string]string{"message": "Password reset."}))
}

package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gofrs/uuid"
)

// HomeHandler is a default handler to serve up
// a home page.
func HomeHandler(c buffalo.Context) error {
	st := c.Session().Get("session_test")
	if st == nil {
		return c.Render(501, r.String("Session test failed"))
	}
	c.Set("session_test", st.(uuid.UUID).String())
	return c.Render(200, r.HTML("index.html"))
}

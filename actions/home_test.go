package actions

import (
	"github.com/gobuffalo/uuid"
)

func (as *ActionSuite) Test_HomeHandler() {
	u, _ := uuid.NewV4()
	as.Session.Set("session_test", u)
	res := as.HTML("/").Get()
	as.Equal(200, res.Code)
	as.Contains(res.Body.String(), "Hello from Buffalo "+u.String())
}

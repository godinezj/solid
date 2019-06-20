package actions

import (
	"bitbucket.org/godinezj/solid/log"
	"bitbucket.org/godinezj/solid/models"
	"bitbucket.org/godinezj/solid/test"
)

func (as *ActionSuite) Test_Vpn_Create() {
	test.SetupTest()
	log.Info("Running VPN Create test")
	resp := as.JSON("/vpn/create").Post("")
	as.Equal(302, resp.Code)

	// create user, and set userID
	u := &models.User{
		FirstName:       "John",
		LastName:        "Doe",
		Email:           "jd@example.com",
		Zip:             "90210",
		Password:        "P@ssw0rd!",
		PasswordConfirm: "P@ssw0rd!",
	}
	verrs, err := u.Create(as.DB)
	as.NoError(err)
	as.False(verrs.HasAny())
	as.Session.Set("current_user_id", u.ID)

	log.Infof("Tesing with user UUID: %v", as.Session.Get("current_user_id"))
	resph := as.HTML("/vpn/create").Post("")
	log.Info(resph.Body.String())
	as.Equal(200, resph.Code)
}

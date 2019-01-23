package actions

import "bitbucket.org/godinezj/solid/models"

// func (as *ActionSuite) Test_UsersResource_List() {
// 	as.Fail("Not Implemented!")
// }

// func (as *ActionSuite) Test_UsersResource_Show() {
// 	as.Fail("Not Implemented!")
// }

// func (as *ActionSuite) Test_UsersResource_New() {
// 	as.Fail("Not Implemented!")
// }

func (as *ActionSuite) Test_UsersResource_Create() {
	u := &models.User{
		FirstName:       "John",
		LastName:        "Doe",
		Email:           "jd@example.com",
		Zip:             "90210",
		Password:        "P@ssw0rd!",
		PasswordConfirm: "P@ssw0rd!",
	}
	as.JSON("/users").Post(u)
	// as.Equal(resp.Code, 200)
	err := as.DB.First(u)
	as.NoError(err)
	as.NotZero(u.ID)
	as.Equal("John", u.FirstName)
}

// func (as *ActionSuite) Test_UsersResource_Edit() {
// 	as.Fail("Not Implemented!")
// }

// func (as *ActionSuite) Test_UsersResource_Update() {
// 	as.Fail("Not Implemented!")
// }

// func (as *ActionSuite) Test_UsersResource_Destroy() {
// 	as.Fail("Not Implemented!")
// }

package actions

import (
	"io/ioutil"
	"os"

	"bitbucket.org/godinezj/solid/log"
	"bitbucket.org/godinezj/solid/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
)

func VPNShow(c buffalo.Context) error {

	return c.Render(200, r.Plain("client_config.txt"))
}

// VPNCreate gives requesting user their vpn config.
func VPNCreate(c buffalo.Context) error {
	// check user id set in session
	uid := c.Session().Get("current_user_id")
	if uid == nil {
		return c.Redirect(302, "/login")
	}

	errMessage := "An unexpected error occured"
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		log.Error("No transaction found")
		return c.Render(422, r.JSON(map[string]string{"message": errMessage}))
	}
	u := &models.User{}
	err := tx.Find(u, uid.(uuid.UUID))
	if err != nil {
		return c.Render(422, r.JSON(map[string]string{"message": errMessage}))
	}

	log.Infof("Creating VPN config %s", u.ID)
	vpn := &models.VPN{}
	vpn.UserID = u.ID

	// Crete the vpn config
	verrs, err := vpn.Create(tx)
	if err != nil {
		log.Error(err)
		return c.Render(422, r.JSON(map[string]string{"message": errMessage}))
	}

	if verrs.HasAny() {
		return c.Redirect(302, "/vpn/show")
	}
	// TODO move into model
	caCertData, err := ioutil.ReadFile(os.Getenv("EASYRSA_DIR") + "/pki/ca.crt")
	if err != nil {
		log.Error(err)
	}
	c.Set("CACert", string(caCertData))
	c.Set("vpn", vpn)
	return c.Render(200, r.Plain("vpn/client_config.txt"))
}

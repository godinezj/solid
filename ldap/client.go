package ldap

import (
	"fmt"
	"os"

	ldap "gopkg.in/ldap.v2"
)

type Client struct {
	Conn *ldap.Conn
}

func (c *Client) Connect() error {
	// connect
	conn, err := ldap.Dial("tcp", os.Getenv("LDAP_BIND_HOST"))
	if err != nil {
		return err
	}

	// authenticate
	err = conn.Bind(os.Getenv("LDAP_BIND_USER"), os.Getenv("LDAP_BIND_PASS"))
	if err != nil {
		return err
	}

	c.Conn = conn
	return nil
}

func (c *Client) Authenticate(username, password string) error {
	userDN := fmt.Sprintf("uid=%s,ou=Users,dc=solidly,dc=io", username)
	err := c.Conn.Bind(userDN, password)
	return err
}

func (c *Client) AddUser(username string) (string, error) {
	return "", nil
}

func (c *Client) Close() {
	c.Conn.Close()
	c.Conn = nil
}

package ldap

import (
	"crypto/tls"
	"fmt"
	"os"

	"bitbucket.org/godinezj/solid/log"
	ldap "gopkg.in/ldap.v2"
)

type Client struct {
	Conn *ldap.Conn
}

// Opens a TCP connection to LDAP
func (c *Client) Connect() error {
	config := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         os.Getenv("LDAP_SSL_HOSTNAME"),
	}
	// connect
	conn, err := ldap.DialTLS("tcp", os.Getenv("LDAP_BIND_HOST"), config)
	if err != nil {
		return err
	}
	c.Conn = conn
	return nil
}

// Authenticates with admin credentials, use this for administrative functions
func (c *Client) AdminAuth() error {
	if c.Conn == nil {
		return fmt.Errorf("Connection not bound to LDAP")
	}
	// authenticate
	err := c.Conn.Bind(os.Getenv("LDAP_BIND_USER"), os.Getenv("LDAP_BIND_PASS"))
	if err != nil {
		return err
	}
	return nil
}

// Adds user, requires an admin connection
func (c *Client) AddUser(first, last, username, password string) (string, error) {
	if c.Conn == nil {
		return "", fmt.Errorf("Connection not bound to LDAP")
	}
	dn := fmt.Sprintf(os.Getenv("LDAP_USERS_DN"), username)
	log.Infof("Adding dn: %s", dn)
	addReq := ldap.NewAddRequest(dn)
	cn := fmt.Sprintf("%s %s", first, last)
	log.Infof("Adding cn: %s", cn)
	addReq.Attribute("cn", []string{cn})
	addReq.Attribute("sn", []string{last})
	addReq.Attribute("uid", []string{username})
	addReq.Attribute("ou", []string{"Users"})
	addReq.Attribute("objectClass", []string{"organizationalPerson", "inetOrgPerson"})
	err := c.Conn.Add(addReq)
	if err != nil {
		return "", err
	}

	passModReq := ldap.NewPasswordModifyRequest(dn, "", password)
	passModRes, err := c.Conn.PasswordModify(passModReq)
	if err != nil {
		return "", err
	}
	return passModRes.GeneratedPassword, nil
}

// Authenticates users using simple username and password
func (c *Client) Authenticate(username, password string) error {
	if c.Conn == nil {
		return fmt.Errorf("Connection not bound to LDAP")
	}
	userDN := fmt.Sprintf(os.Getenv("LDAP_USERS_DN"), username)
	log.Infof("Authenticating user: %s with password %s", userDN, password)
	err := c.Conn.Bind(userDN, password)
	return err
}

// Closes the connection
func (c *Client) Close() {
	if c.Conn != nil {
		c.Conn.Close()
	}
	c.Conn = nil
}

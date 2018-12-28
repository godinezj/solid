package ldap

import (
	"os"
	"testing"
)

func Test_Client(t *testing.T) {
	os.Setenv("LDAP_BIND_HOST", "127.0.0.1:389")
	os.Setenv("LDAP_BIND_USER", "cn=admin,dc=solidly,dc=io")
	os.Setenv("LDAP_BIND_PASS", "8RlDnnSb1Kce")
	client := Client{}
	err := client.Connect()
	if err != nil {
		t.Error(err)
	}

	err = client.Authenticate("mikeg", "1I.SaeCh")
	if err != nil {
		t.Error(err)
	}
	client.Close()
}

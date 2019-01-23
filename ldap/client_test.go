package ldap

import (
	"fmt"
	"testing"

	"github.com/joho/godotenv"
)

func Test_Client(t *testing.T) {
	// Load Client dependencies
	err := godotenv.Load("../.env")
	if err != nil {
		t.Fatal("Error loading .env file")
	}
	client := Client{}
	err = client.Connect()
	if err != nil {
		t.Error(err)
	}
	err = client.AdminAuth()
	if err != nil {
		t.Error(err)
	}
	// client.Conn.Debug = true
	// add user
	username := "testuser8"
	pass, err := client.AddUser("test", "user", username, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(pass) < 1 {
		t.Error("Received empty password")
	} else {
		fmt.Printf("New user password is %s", pass)
	}
	client.Close() // close the admin connection
	client.Connect()
	// authenticate user
	err = client.Authenticate(username, pass)
	if err != nil {
		t.Error(err)
	}
	client.Close()
}

package models

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	"bitbucket.org/godinezj/solid/log"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

type VPN struct {
	ID          uuid.UUID `json:"id" db:"id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	PrivateKey  string    `json:"private_key" db:"private_key"`
	Certificate string    `json:"certificate" db:"certificate"`
	CACert      string    `json:"ca_cert" db:"-"`
	TLSAuthKey  string    `json:"tls_auth_key" db:"-"`
}

// String is not required by pop and may be deleted
func (v VPN) String() string {
	jv, _ := json.Marshal(v)
	return string(jv)
}

// Vpns is not required by pop and may be deleted
type Vpns []VPN

// String is not required by pop and may be deleted
func (v Vpns) String() string {
	jv, _ := json.Marshal(v)
	return string(jv)
}

func (v *VPN) buildClient(userID string) error {
	log.Infof("Creating VPN client config for %s", userID)

	// sanity check
	_, err := exec.LookPath("easyrsa")
	if err != nil {
		return errors.New("easyrsa not installed")
	}
	// generate client certificates
	cmd := exec.Command("./easyrsa", "build-client-full", userID, "nopass")
	cmd.Dir = os.Getenv("EASYRSA_DIR")
	var buf bytes.Buffer
	bufWriter := bufio.NewWriter(&buf)
	cmd.Stdout = bufWriter
	cmd.Stderr = bufWriter
	log.Info(cmd.Args)
	err = cmd.Run()
	if err != nil {
		log.Error(buf.String())
		log.Error(err)
		return err
	}
	log.Info(buf.String())
	log.Infof("VPN client config created for %s", userID)
	return nil
}

func (v *VPN) cert(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	beginCert := false
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Text() == "-----BEGIN CERTIFICATE-----" {
			beginCert = true
		}
		if beginCert {
			buf.WriteString(scanner.Text() + "\n")
		}
	}
	return buf.String(), nil
}

// Create client config and saves it to the backend
func (v *VPN) Create(tx *pop.Connection) (*validate.Errors, error) {
	log.Info("Creating VPN client config")
	// check to if config already exists
	userID := v.UserID.String()
	query := tx.Where("user_id=?", userID).Select("user_id")
	queryVPN := VPN{}
	err := query.First(&queryVPN)
	if err == nil {
		log.Infof("VPN config already exists for %s", userID)
		verrs := validate.NewErrors()
		verrs.Add("user_id", "VPN config already exists")
		return verrs, errors.WithStack(err)
	}
	// build-client-full
	err = v.buildClient(userID)
	if err != nil {
		return validate.NewErrors(), err
	}
	// load private key into model
	pkiPath := os.Getenv("EASYRSA_DIR") + "/pki"
	keyData, err := ioutil.ReadFile(pkiPath + "/private/" + userID + ".key")
	if err != nil {
		log.Errorf("Could not read private key for %s", userID)
		return validate.NewErrors(), err
	}
	v.PrivateKey = string(keyData)
	// load cert into model
	v.Certificate, err = v.cert(pkiPath + "/issued/" + userID + ".crt")
	if err != nil {
		log.Errorf("Could not read certificate for %s", userID)
		log.Error(err)
		return validate.NewErrors(), err
	}
	// load tls-auth key into model
	tlsAuthData, err := ioutil.ReadFile(pkiPath + "/ta.key")
	if err != nil {
		log.Errorf("Could not read private key for %s", userID)
		return validate.NewErrors(), err
	}
	v.TLSAuthKey = string(tlsAuthData)

	// save config
	return tx.ValidateAndCreate(v)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (v *VPN) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (v *VPN) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (v *VPN) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

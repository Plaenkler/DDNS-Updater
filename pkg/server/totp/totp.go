package totp

import (
	"os"

	log "github.com/plaenkler/ddns-updater/pkg/logging"

	"github.com/pquerna/otp/totp"
)

const (
	filePerm   = 0644
	secretPath = "./data/TOTPSecret"
)

var (
	keySecret string
)

func init() {
	secret, err := readKeySecret()
	if err != nil {
		log.Fatalf("[totp-init-1] could not load secret: %v", err)
	}
	keySecret = secret
}

func readKeySecret() (string, error) {
	secret, err := os.ReadFile(secretPath)
	if err != nil && !os.IsNotExist(err) {
		return "", err
	}
	if os.IsNotExist(err) {
		secret, err = createKeySecret()
		if err != nil {
			return "", err
		}
	}
	return string(secret), nil
}

func createKeySecret() ([]byte, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "DDNS-Updater",
		AccountName: "Administrator",
	})
	if err != nil {
		return nil, err
	}
	secret := []byte(key.Secret())
	err = os.WriteFile(secretPath, secret, filePerm)
	if err != nil {
		return nil, err
	}
	return secret, nil
}

func VerifiyTOTP(otp string) bool {
	return totp.Validate(otp, keySecret)
}

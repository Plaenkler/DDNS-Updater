package totp

import (
	"bytes"
	"encoding/base64"
	"image/png"
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

func GetKeyAsQR() (string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "DDNS-Updater",
		AccountName: "Administrator",
		Secret:      []byte(keySecret),
	})
	if err != nil {
		return "", err
	}
	img, err := key.Image(200, 200)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = png.Encode(&buf, img)
	if err != nil {
		return "", err
	}
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func VerifiyTOTP(otp string) bool {
	return totp.Validate(otp, keySecret)
}

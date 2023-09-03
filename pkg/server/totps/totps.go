package totps

import (
	"bytes"
	"encoding/base32"
	"encoding/base64"
	"image/png"
	"os"
	"path/filepath"

	log "github.com/plaenkler/ddns-updater/pkg/logging"

	"github.com/pquerna/otp/totp"
)

const (
	dirPerm     = 0755
	filePerm    = 0644
	secretPath  = "./data/TOTPSecret"
	issuer      = "DDNS-Updater"
	accountName = "Administrator"
)

var (
	keySecret string
)

func init() {
	var err error
	keySecret, err = read()
	if err != nil {
		log.Fatalf("[totp-init-1] could not load secret: %v", err)
	}
}

func read() (string, error) {
	secret, err := os.ReadFile(secretPath)
	if err != nil && !os.IsNotExist(err) {
		return "", err
	}
	if os.IsNotExist(err) {
		secret, err = create()
		if err != nil {
			return "", err
		}
	}
	return string(secret), nil
}

func create() ([]byte, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: accountName,
	})
	if err != nil {
		return nil, err
	}
	err = os.MkdirAll(filepath.Dir(secretPath), dirPerm)
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
	secret, err := base32.StdEncoding.DecodeString(keySecret)
	if err != nil {
		return "", err
	}
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: accountName,
		Secret:      secret,
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
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func Verify(otp string) bool {
	return totp.Validate(otp, keySecret)
}

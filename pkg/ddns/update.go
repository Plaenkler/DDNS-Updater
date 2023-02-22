package ddns

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/plaenkler/ddns/pkg/model"
)

type Update func(job model.SyncJob, ipAddr string) error

var updaters = map[string]Update{
	"Strato":       updateStrato,
	"DDNSS":        updateDDNSS,
	"Dynu":         updateDynu,
	"Aliyun":       updateAliyun,
	"AllInkl":      updateAllInkl,
	"Cloudflare":   updateCloudflare,
	"DD24":         updateDD24,
	"DigitalOcean": updateDigitalOcean,
	"DonDominio":   updateDonDominio,
	"DNSOMatic":    updateDNSOMatic,
	"DNSPod":       updateDNSPod,
	"Dreamhost":    updateDreamhost,
	"DuckDNS":      updateDuckDNS,
	"DynDNS":       updateDynDNS,
	"FreeDNS":      updateFreeDNS,
	"Gandi":        updateGandi,
	"GCP":          updateGCP,
	"GoDaddy":      updateGoDaddy,
	"Google":       updateGoogle,
	"He":           updateHe,
	"Infomaniak":   updateInfomaniak,
	"INWX":         updateINWX,
	"Linode":       updateLinode,
	"LuaDNS":       updateLuaDNS,
	"Namecheap":    updateNamecheap,
	"NoIP":         updateNoIP,
	"Njalla":       updateNjalla,
	"OpenDNS":      updateOpenDNS,
	"OVH":          updateOVH,
	"Porkbun":      updatePorkbun,
	"Selfhost":     updateSelfhost,
	"Servercow":    updateServercow,
	"Spdyn":        updateSpdyn,
	"Variomedia":   updateVariomedia,
}

func IsUpdaterSupported(updater string) bool {
	_, ok := updaters[updater]
	return ok
}

func sendHTTPRequest(method string, url string, auth *url.Userinfo) (string, error) {
	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %v", err)
	}
	if auth != nil {
		password, _ := auth.Password()
		req.SetBasicAuth(auth.Username(), password)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to update DDNS entry: %s", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read HTTP response body: %v", err)
	}
	return string(bytes.TrimSpace(body)), nil
}

func updateStrato(job model.SyncJob, ipAddr string) error {
	urlStr := fmt.Sprintf("https://%s:%s@dyndns.strato.com/nic/update?hostname=%s&myip=%s", job.User, job.Password, job.Domain, ipAddr)
	resp, err := sendHTTPRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return err
	}
	switch {
	case strings.Contains(resp, "good"):
		return nil
	case strings.Contains(resp, "nochg"):
		return nil
	default:
		return fmt.Errorf("failed to update DDNS entry: %s", resp)
	}
}

func updateDDNSS(job model.SyncJob, ipAddr string) error {
	urlStr := fmt.Sprintf("https://www.ddnss.de/upd.php?user=%s&pwd=%s&host=%s&ip=%s&ip6=%s", job.User, job.Password, job.Domain, ipAddr, "")
	resp, err := sendHTTPRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return err
	}
	if strings.Contains(resp, "Error Occurred While Processing Request") {
		return fmt.Errorf("failed to update DDNS entry: %s", resp)
	}
	return nil
}

func updateDynu(job model.SyncJob, ipAddr string) error {
	urlStr := fmt.Sprintf("https://%s:%s@api.dynu.com/nic/update?myip=%s&myipv6=", job.User, job.Password, ipAddr)
	resp, err := sendHTTPRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return err
	}
	switch {
	case strings.HasPrefix(resp, "good"):
		return nil
	case strings.HasPrefix(resp, "nochg"):
		return nil
	default:
		return fmt.Errorf("failed to update DDNS entry: %s", resp)
	}
}

func updateAliyun(job model.SyncJob, ipAddr string) error {
	return nil
}

func updateAllInkl(job model.SyncJob, ipAddr string) error {
	return nil
}

func updateCloudflare(job model.SyncJob, ipAddr string) error {
	return nil
}

func updateDD24(job model.SyncJob, ipAddr string) error {
	return nil
}

func updateDigitalOcean(job model.SyncJob, ipAddr string) error {
	return nil
}

func updateDonDominio(job model.SyncJob, ipAddr string) error {
	return nil
}

func updateDNSOMatic(job model.SyncJob, ipAddr string) error {
	return nil
}

func updateDNSPod(job model.SyncJob, ipAddr string) error {
	return nil
}

func updateDreamhost(job model.SyncJob, ipAddr string) error {
	return nil
}

func updateDuckDNS(job model.SyncJob, ipAddr string) error {
	return nil
}

func updateDynDNS(job model.SyncJob, ipAddr string) error {
	return nil
}

func updateFreeDNS(job model.SyncJob, ipAddr string) error {
	return nil
}

func updateGandi(job model.SyncJob, ipAddr string) error {
	return nil
}

func updateGCP(job model.SyncJob, ipAddr string) error {
	return nil
}

func updateGoDaddy(job model.SyncJob, ipAddr string) error {
	return nil
}

func updateGoogle(job model.SyncJob, ipAddr string) error {
	return nil
}

func updateHe(job model.SyncJob, ipAddr string) error {
	return nil
}

func updateInfomaniak(job model.SyncJob, ipAddr string) error {
	return nil
}

func updateINWX(job model.SyncJob, ipAddr string) error {
	return nil
}

func updateLinode(job model.SyncJob, ipAddr string) error {
	return nil
}

func updateLuaDNS(job model.SyncJob, ipAddr string) error {
	return nil
}

func updateNamecheap(job model.SyncJob, ipAddr string) error {
	return nil
}

func updateNoIP(job model.SyncJob, ipAddr string) error {
	return nil
}

func updateNjalla(job model.SyncJob, ipAddr string) error {
	return nil
}

func updateOpenDNS(job model.SyncJob, ipAddr string) error {
	return nil
}

func updateOVH(job model.SyncJob, ipAddr string) error {
	return nil
}

func updatePorkbun(job model.SyncJob, ipAddr string) error {
	return nil
}

func updateSelfhost(job model.SyncJob, ipAddr string) error {
	return nil
}

func updateServercow(job model.SyncJob, ipAddr string) error {
	return nil
}

func updateSpdyn(job model.SyncJob, ipAddr string) error {
	return nil
}

func updateVariomedia(job model.SyncJob, ipAddr string) error {
	return nil
}

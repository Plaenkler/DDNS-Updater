package ddns

import (
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/plaenkler/ddns-updater/pkg/config"
	log "github.com/plaenkler/ddns-updater/pkg/logging"
	"golang.org/x/net/dns/dnsmessage"
)

var dnsResolvers = map[string]string{
	"cf-primary":   "1.1.1.1:853",
	"cf-secondary": "1.0.0.1:853",
}

var httpResolvers = map[string]string{
	"ipify":  "https://api4.ipify.org",
	"ipych":  "https://api.ipy.ch",
	"intel":  "https://nms.intellitrend.de",
	"ident":  "https://4.ident.me",
	"ifconf": "https://ifconfig.me/ip",
	"icanh":  "https://ipv4.icanhazip.com",
	"aws":    "https://checkip.amazonaws.com",
	"ipinfo": "https://ipinfo.io/ip",
	"ipsb":   "https://api-ipv4.ip.sb/ip",
	"seeip":  "https://api.seeip.org",
	"wtf":    "https://wtfismyip.com/text",
	"mullv":  "https://am.i.mullvad.net/ip",
}

func GetPublicIP() (string, error) {
	cRes := config.Get().Resolver
	if cRes != "" {
		addr, err := resolveIPAddressHTTP(cRes)
		if err != nil {
			return "", fmt.Errorf("resolver %s failed: %s", cRes, err)
		}
		log.Infof("%s succeeded: %s", cRes, addr)
		return addr, nil
	}
	for r := range dnsResolvers {
		addr, err := resolveIPAddressDNS(dnsResolvers[r])
		if err != nil {
			log.Errorf("resolver %s failed: %s", r, err)
			continue
		}
		log.Infof("%s succeeded: %s", r, addr)
		return addr, nil
	}
	for r := range httpResolvers {
		addr, err := resolveIPAddressHTTP(httpResolvers[r])
		if err != nil {
			log.Errorf("resolver %s failed: %s", r, err)
			continue
		}
		log.Infof("%s succeeded: %s", r, addr)
		return addr, nil
	}
	return "", fmt.Errorf("all resolvers failed")
}

func resolveIPAddressHTTP(url string) (string, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer log.ErrorClose(resp.Body)
	bytes, err := io.ReadAll(io.LimitReader(resp.Body, 15))
	if err != nil {
		return "", err
	}
	addr := strings.TrimSpace(string(bytes))
	if !isValidIPAddress(addr) {
		return "", fmt.Errorf("invalid IP address: %q", addr)
	}
	return addr, nil
}

func resolveIPAddressDNS(server string) (string, error) {
	msg, err := exchangeDNS(server)
	if err != nil {
		return "", err
	}
	for _, answer := range msg.Answers {
		record, ok := answer.Body.(*dnsmessage.TXTResource)
		if !ok {
			continue
		}
		addr := strings.Join(record.TXT, "")
		if isValidIPAddress(addr) {
			return addr, nil
		}
	}
	return "", fmt.Errorf("no valid IP address in response")
}

func exchangeDNS(server string) (*dnsmessage.Message, error) {
	question := dnsmessage.Question{
		Name:  dnsmessage.MustNewName("whoami.cloudflare."),
		Type:  dnsmessage.TypeTXT,
		Class: dnsmessage.ClassCHAOS,
	}
	query := dnsmessage.Message{Questions: []dnsmessage.Question{question}}
	req, err := query.AppendPack(make([]byte, 2))
	if err != nil {
		return nil, err
	}
	binary.BigEndian.PutUint16(req, uint16(len(req)-2))
	deadline := time.Now().Add(5 * time.Second)
	dialer := &net.Dialer{Deadline: deadline}
	tlsConfig := &tls.Config{
		ServerName: "cloudflare-dns.com",
		MinVersion: tls.VersionTLS12,
	}
	conn, err := tls.DialWithDialer(dialer, "tcp4", server, tlsConfig)
	if err != nil {
		return nil, err
	}
	defer log.ErrorClose(conn)
	err = conn.SetDeadline(deadline)
	if err != nil {
		return nil, err
	}
	_, err = conn.Write(req)
	if err != nil {
		return nil, err
	}
	_, err = io.ReadFull(conn, req[:2])
	if err != nil {
		return nil, err
	}
	resp := make([]byte, binary.BigEndian.Uint16(req[:2]))
	_, err = io.ReadFull(conn, resp)
	if err != nil {
		return nil, err
	}
	var msg dnsmessage.Message
	err = msg.Unpack(resp)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

func isValidIPAddress(ip string) bool {
	addr := net.ParseIP(ip)
	if addr == nil {
		return false
	}
	if addr.IsUnspecified() {
		return false
	}
	if addr.IsPrivate() {
		return false
	}
	if addr.IsLoopback() {
		return false
	}
	if addr.IsMulticast() {
		return false
	}
	return true
}

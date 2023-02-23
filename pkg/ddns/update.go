package ddns

import "github.com/plaenkler/ddns/pkg/ddns/providers"

type Updater func(request interface{}) error

type Provider struct {
	Updater Updater
	Request interface{}
}

var updaters = map[string]Provider{
	"Strato":       {Updater: providers.UpdateStrato, Request: providers.UpdateStratoRequest{}},
	"DDNSS":        {Updater: providers.UpdateDDNSS, Request: providers.UpdateDDNSSRequest{}},
	"Dynu":         {Updater: providers.UpdateDynu, Request: providers.UpdateDynuRequest{}},
	"Aliyun":       {Updater: providers.UpdateAliyun, Request: providers.UpdateAliyunRequest{}},
	"AllInkl":      {Updater: providers.UpdateAllInkl, Request: providers.UpdateAllInklRequest{}},
	"Cloudflare":   {Updater: providers.UpdateCloudflare, Request: providers.UpdateCloudflareRequest{}},
	"DD24":         {Updater: providers.UpdateDD24, Request: providers.UpdateDD24Request{}},
	"DigitalOcean": {Updater: providers.UpdateDigitalOcean, Request: providers.UpdateDigitalOceanRequest{}},
	"DonDominio":   {Updater: providers.UpdateDonDominio, Request: providers.UpdateDonDominioRequest{}},
	"DNSOMatic":    {Updater: providers.UpdateDNSOMatic, Request: providers.UpdateDNSOMaticRequest{}},
	"DNSPod":       {Updater: providers.UpdateDNSPod, Request: providers.UpdateDNSPodRequest{}},
	"Dreamhost":    {Updater: providers.UpdateDreamhost, Request: providers.UpdateDreamhostRequest{}},
	"DuckDNS":      {Updater: providers.UpdateDuckDNS, Request: providers.UpdateDuckDNSRequest{}},
	"DynDNS":       {Updater: providers.UpdateDynDNS, Request: providers.UpdateDynDNSRequest{}},
	"FreeDNS":      {Updater: providers.UpdateFreeDNS, Request: providers.UpdateFreeDNSRequest{}},
	"Gandi":        {Updater: providers.UpdateGandi, Request: providers.UpdateGandiRequest{}},
	"GCP":          {Updater: providers.UpdateGCP, Request: providers.UpdateGCPRequest{}},
	"GoDaddy":      {Updater: providers.UpdateGoDaddy, Request: providers.UpdateGoDaddyRequest{}},
	"Google":       {Updater: providers.UpdateGoogle, Request: providers.UpdateGoogleRequest{}},
	"He":           {Updater: providers.UpdateHe, Request: providers.UpdateHeRequest{}},
	"Infomaniak":   {Updater: providers.UpdateInfomaniak, Request: providers.UpdateInfomaniakRequest{}},
	"INWX":         {Updater: providers.UpdateINWX, Request: providers.UpdateINWXRequest{}},
	"Linode":       {Updater: providers.UpdateLinode, Request: providers.UpdateLinodeRequest{}},
	"LuaDNS":       {Updater: providers.UpdateLuaDNS, Request: providers.UpdateLuaDNSRequest{}},
	"Namecheap":    {Updater: providers.UpdateNamecheap, Request: providers.UpdateNamecheapRequest{}},
	"NoIP":         {Updater: providers.UpdateNoIP, Request: providers.UpdateNoIPRequest{}},
	"Njalla":       {Updater: providers.UpdateNjalla, Request: providers.UpdateNjallaRequest{}},
	"OpenDNS":      {Updater: providers.UpdateOpenDNS, Request: providers.UpdateOpenDNSRequest{}},
	"OVH":          {Updater: providers.UpdateOVH, Request: providers.UpdateOVHRequest{}},
	"Porkbun":      {Updater: providers.UpdatePorkbun, Request: providers.UpdatePorkbunRequest{}},
	"Selfhost":     {Updater: providers.UpdateSelfhost, Request: providers.UpdateSelfhostRequest{}},
	"Servercow":    {Updater: providers.UpdateServercow, Request: providers.UpdateServercowRequest{}},
	"Spdyn":        {Updater: providers.UpdateSpdyn, Request: providers.UpdateSpdynRequest{}},
	"Variomedia":   {Updater: providers.UpdateVariomedia, Request: providers.UpdateVariomediaRequest{}},
}

func IsSupported(p string) bool {
	_, ok := updaters[p]
	return ok
}

func GetProviders() []string {
	var p []string
	for k := range updaters {
		p = append(p, k)
	}
	return p
}

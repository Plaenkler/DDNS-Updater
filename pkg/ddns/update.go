package ddns

import "github.com/plaenkler/ddns/pkg/ddns/providers"

type Updater func(request interface{}) error

type Provider struct {
	Updater Updater
	Request func() interface{}
}

var updaters = map[string]Provider{
	"Strato":       {Updater: providers.UpdateStrato, Request: func() interface{} { return &providers.UpdateStratoRequest{} }},
	"DDNSS":        {Updater: providers.UpdateDDNSS, Request: func() interface{} { return &providers.UpdateDDNSSRequest{} }},
	"Dynu":         {Updater: providers.UpdateDynu, Request: func() interface{} { return &providers.UpdateDynuRequest{} }},
	"Aliyun":       {Updater: providers.UpdateAliyun, Request: func() interface{} { return &providers.UpdateAliyunRequest{} }},
	"AllInkl":      {Updater: providers.UpdateAllInkl, Request: func() interface{} { return &providers.UpdateAllInklRequest{} }},
	"Cloudflare":   {Updater: providers.UpdateCloudflare, Request: func() interface{} { return &providers.UpdateCloudflareRequest{} }},
	"DD24":         {Updater: providers.UpdateDD24, Request: func() interface{} { return &providers.UpdateDD24Request{} }},
	"DigitalOcean": {Updater: providers.UpdateDigitalOcean, Request: func() interface{} { return &providers.UpdateDigitalOceanRequest{} }},
	"DonDominio":   {Updater: providers.UpdateDonDominio, Request: func() interface{} { return &providers.UpdateDonDominioRequest{} }},
	"DNSOMatic":    {Updater: providers.UpdateDNSOMatic, Request: func() interface{} { return &providers.UpdateDNSOMaticRequest{} }},
	"DNSPod":       {Updater: providers.UpdateDNSPod, Request: func() interface{} { return &providers.UpdateDNSPodRequest{} }},
	"Dreamhost":    {Updater: providers.UpdateDreamhost, Request: func() interface{} { return &providers.UpdateDreamhostRequest{} }},
	"DuckDNS":      {Updater: providers.UpdateDuckDNS, Request: func() interface{} { return &providers.UpdateDuckDNSRequest{} }},
	"DynDNS":       {Updater: providers.UpdateDynDNS, Request: func() interface{} { return &providers.UpdateDynDNSRequest{} }},
	"FreeDNS":      {Updater: providers.UpdateFreeDNS, Request: func() interface{} { return &providers.UpdateFreeDNSRequest{} }},
	"Gandi":        {Updater: providers.UpdateGandi, Request: func() interface{} { return &providers.UpdateGandiRequest{} }},
	"GCP":          {Updater: providers.UpdateGCP, Request: func() interface{} { return &providers.UpdateGCPRequest{} }},
	"GoDaddy":      {Updater: providers.UpdateGoDaddy, Request: func() interface{} { return &providers.UpdateGoDaddyRequest{} }},
	"Google":       {Updater: providers.UpdateGoogle, Request: func() interface{} { return &providers.UpdateGoogleRequest{} }},
	"He":           {Updater: providers.UpdateHe, Request: func() interface{} { return &providers.UpdateHeRequest{} }},
	"Infomaniak":   {Updater: providers.UpdateInfomaniak, Request: func() interface{} { return &providers.UpdateInfomaniakRequest{} }},
	"INWX":         {Updater: providers.UpdateINWX, Request: func() interface{} { return &providers.UpdateINWXRequest{} }},
	"Linode":       {Updater: providers.UpdateLinode, Request: func() interface{} { return &providers.UpdateLinodeRequest{} }},
	"LuaDNS":       {Updater: providers.UpdateLuaDNS, Request: func() interface{} { return &providers.UpdateLuaDNSRequest{} }},
	"Namecheap":    {Updater: providers.UpdateNamecheap, Request: func() interface{} { return &providers.UpdateNamecheapRequest{} }},
	"NoIP":         {Updater: providers.UpdateNoIP, Request: func() interface{} { return &providers.UpdateNoIPRequest{} }},
	"Njalla":       {Updater: providers.UpdateNjalla, Request: func() interface{} { return &providers.UpdateNjallaRequest{} }},
	"OpenDNS":      {Updater: providers.UpdateOpenDNS, Request: func() interface{} { return &providers.UpdateOpenDNSRequest{} }},
	"OVH":          {Updater: providers.UpdateOVH, Request: func() interface{} { return &providers.UpdateOVHRequest{} }},
	"Porkbun":      {Updater: providers.UpdatePorkbun, Request: func() interface{} { return &providers.UpdatePorkbunRequest{} }},
	"Selfhost":     {Updater: providers.UpdateSelfhost, Request: func() interface{} { return &providers.UpdateSelfhostRequest{} }},
	"Servercow":    {Updater: providers.UpdateServercow, Request: func() interface{} { return &providers.UpdateServercowRequest{} }},
	"Spdyn":        {Updater: providers.UpdateSpdyn, Request: func() interface{} { return &providers.UpdateSpdynRequest{} }},
	"Variomedia":   {Updater: providers.UpdateVariomedia, Request: func() interface{} { return &providers.UpdateVariomediaRequest{} }},
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

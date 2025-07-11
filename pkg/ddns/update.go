package ddns

import (
	"sort"

	"github.com/plaenkler/ddns-updater/pkg/ddns/providers"
)

type updater func(request interface{}, ipAddr string) error

type requestFactory func() interface{}

type provider struct {
	Updater updater
	Factory requestFactory
}

var updaters = map[string]provider{
	"Custom":       {Updater: providers.UpdateCustom, Factory: func() interface{} { return &providers.UpdateCustomRequest{} }},
	"Strato":       {Updater: providers.UpdateStrato, Factory: func() interface{} { return &providers.UpdateStratoRequest{} }},
	"DDNSS":        {Updater: providers.UpdateDDNSS, Factory: func() interface{} { return &providers.UpdateDDNSSRequest{} }},
	"Dynu":         {Updater: providers.UpdateDynu, Factory: func() interface{} { return &providers.UpdateDynuRequest{} }},
	"Aliyun":       {Updater: providers.UpdateAliyun, Factory: func() interface{} { return &providers.UpdateAliyunRequest{} }},
	"AllInkl":      {Updater: providers.UpdateAllInkl, Factory: func() interface{} { return &providers.UpdateAllInklRequest{} }},
	"Cloudflare":   {Updater: providers.UpdateCloudflare, Factory: func() interface{} { return &providers.UpdateCloudflareRequest{} }},
	"DD24":         {Updater: providers.UpdateDD24, Factory: func() interface{} { return &providers.UpdateDD24Request{} }},
	"DigitalOcean": {Updater: providers.UpdateDigitalOcean, Factory: func() interface{} { return &providers.UpdateDigitalOceanRequest{} }},
	"DonDominio":   {Updater: providers.UpdateDonDominio, Factory: func() interface{} { return &providers.UpdateDonDominioRequest{} }},
	"DNSOMatic":    {Updater: providers.UpdateDNSOMatic, Factory: func() interface{} { return &providers.UpdateDNSOMaticRequest{} }},
	"DNSPod":       {Updater: providers.UpdateDNSPod, Factory: func() interface{} { return &providers.UpdateDNSPodRequest{} }},
	"Dreamhost":    {Updater: providers.UpdateDreamhost, Factory: func() interface{} { return &providers.UpdateDreamhostRequest{} }},
	"DuckDNS":      {Updater: providers.UpdateDuckDNS, Factory: func() interface{} { return &providers.UpdateDuckDNSRequest{} }},
	"DynDNS":       {Updater: providers.UpdateDynDNS, Factory: func() interface{} { return &providers.UpdateDynDNSRequest{} }},
	"FreeDNS":      {Updater: providers.UpdateFreeDNS, Factory: func() interface{} { return &providers.UpdateFreeDNSRequest{} }},
	"Gandi":        {Updater: providers.UpdateGandi, Factory: func() interface{} { return &providers.UpdateGandiRequest{} }},
	"GCP":          {Updater: providers.UpdateGCP, Factory: func() interface{} { return &providers.UpdateGCPRequest{} }},
	"GoDaddy":      {Updater: providers.UpdateGoDaddy, Factory: func() interface{} { return &providers.UpdateGoDaddyRequest{} }},
	"Google":       {Updater: providers.UpdateGoogle, Factory: func() interface{} { return &providers.UpdateGoogleRequest{} }},
	"He":           {Updater: providers.UpdateHe, Factory: func() interface{} { return &providers.UpdateHeRequest{} }},
	"Hetzner":      {Updater: providers.UpdateHetzner, Factory: func() interface{} { return &providers.UpdateHetznerRequest{} }},
	"Infomaniak":   {Updater: providers.UpdateInfomaniak, Factory: func() interface{} { return &providers.UpdateInfomaniakRequest{} }},
	"INWX":         {Updater: providers.UpdateINWX, Factory: func() interface{} { return &providers.UpdateINWXRequest{} }},
	"Linode":       {Updater: providers.UpdateLinode, Factory: func() interface{} { return &providers.UpdateLinodeRequest{} }},
	"LuaDNS":       {Updater: providers.UpdateLuaDNS, Factory: func() interface{} { return &providers.UpdateLuaDNSRequest{} }},
	"Namecheap":    {Updater: providers.UpdateNamecheap, Factory: func() interface{} { return &providers.UpdateNamecheapRequest{} }},
	"NoIP":         {Updater: providers.UpdateNoIP, Factory: func() interface{} { return &providers.UpdateNoIPRequest{} }},
	"Njalla":       {Updater: providers.UpdateNjalla, Factory: func() interface{} { return &providers.UpdateNjallaRequest{} }},
	"OpenDNS":      {Updater: providers.UpdateOpenDNS, Factory: func() interface{} { return &providers.UpdateOpenDNSRequest{} }},
	"OVH":          {Updater: providers.UpdateOVH, Factory: func() interface{} { return &providers.UpdateOVHRequest{} }},
	"Porkbun":      {Updater: providers.UpdatePorkbun, Factory: func() interface{} { return &providers.UpdatePorkbunRequest{} }},
	"Selfhost":     {Updater: providers.UpdateSelfhost, Factory: func() interface{} { return &providers.UpdateSelfhostRequest{} }},
	"Servercow":    {Updater: providers.UpdateServercow, Factory: func() interface{} { return &providers.UpdateServercowRequest{} }},
	"Spdyn":        {Updater: providers.UpdateSpdyn, Factory: func() interface{} { return &providers.UpdateSpdynRequest{} }},
	"Variomedia":   {Updater: providers.UpdateVariomedia, Factory: func() interface{} { return &providers.UpdateVariomediaRequest{} }},
	"MaxiHoster":   {Updater: providers.UpdateMaxiHoster, Factory: func() interface{} { return &providers.UpdateMaxiHosterRequest{} }},
}

func GetProviders() []string {
	var p []string
	for k := range updaters {
		p = append(p, k)
	}
	sort.Strings(p)
	return p
}

func GetUpdaters() map[string]provider {
	return updaters
}

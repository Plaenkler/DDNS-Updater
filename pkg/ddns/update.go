package ddns

import (
	"sort"

	"github.com/plaenkler/ddns-updater/pkg/ddns/providers"
)

type updater func(request interface{}, ipAddr string) error

type provider struct {
	Updater updater
	Request interface{}
}

var updaters = map[string]provider{
	"Custom":     {Updater: providers.UpdateCustom, Request: providers.UpdateCustomRequest{}},
	"Strato":     {Updater: providers.UpdateStrato, Request: providers.UpdateStratoRequest{}},
	"DDNSS":      {Updater: providers.UpdateDDNSS, Request: providers.UpdateDDNSSRequest{}},
	"Dynu":       {Updater: providers.UpdateDynu, Request: providers.UpdateDynuRequest{}},
	"Aliyun":     {Updater: providers.UpdateAliyun, Request: providers.UpdateAliyunRequest{}},
	"DD24":       {Updater: providers.UpdateDD24, Request: providers.UpdateDD24Request{}},
	"Hetzner":    {Updater: providers.UpdateHetzner, Request: providers.UpdateHetznerRequest{}},
	"Infomaniak": {Updater: providers.UpdateInfomaniak, Request: providers.UpdateInfomaniakRequest{}},
	"INWX":       {Updater: providers.UpdateINWX, Request: providers.UpdateINWXRequest{}},
	"IONOS":      {Updater: providers.UpdateIONOS, Request: providers.UpdateIONOSRequest{}},
	"NoIP":       {Updater: providers.UpdateNoIP, Request: providers.UpdateNoIPRequest{}},
	"MaxiHoster": {Updater: providers.UpdateMaxiHoster, Request: providers.UpdateMaxiHosterRequest{}},
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

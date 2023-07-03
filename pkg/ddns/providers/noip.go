//package providers
//
//import (
//	"fmt"
//)
//
//type UpdateNoIPRequest struct {
//    Domain        string
//	  Host          string
//	  Username      string
//	  Password      string
//	  UseProviderIP bool
//}
//
//func UpdateNoIP(request interface{}, ipAddr string) error {
//	  r, ok := request.(*UpdateNoIPRequest)
//	  if !ok {
//	  	return fmt.Errorf("invalid request type: %T", request)
//	  }
//	  return fmt.Errorf("not implemented %s", r.Domain)
//}

//Tom dem Seine Notizen von der Knowledgebase von NOIP.COM:

//1: ExampleUpdateRequestString: "http://username:password@dynupdate.no-ip.com/nic/update?hostname=mytest.example.com&myip=192.0.2.25"

//2: Raw HTTP GET Request:
//"GET
///nic/update?hostname=mytestyourhost.example.com&myip=192.0.2.25
//HTTP/1.1
//Host: dynupdate.no-ip.com
//Authorization: Basic base-64-authorization
//User-Agent: ##SEE SETTING the User-Agent Below##""

//3: For software clients use the format:
//"User-Agent: Company NameOfProgram/OSVersion-ReleaseVersion maintainer-contact@example.com"
//User Agent is important! It helps us identify what Update Client is accessing us. This allows us to provide better technical support to YOUR USERS.

package providers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/plaenkler/ddns/pkg/util"
)

type UpdateNoIPRequest struct {
	// Domain        string
	Host     string
	Username string
	Password string
	// UseProviderIP bool
}

func UpdateNoIP(request interface{}, ipAddr string) error {
	r, ok := request.(*UpdateNoIPRequest)
	if !ok {
		return fmt.Errorf("invalid request type: %T", request)
	}

	//UpdateRequestString:
	urlStr := fmt.Sprintf("http://%s:%s@dynupdate.no-ip.com/nic/update?hostname=%s&myip=%s", r.Username, r.Password, r.Host, ipAddr)
	resp, err := util.SendHTTPRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return err
	}
	if strings.Contains(resp, "Error Occurred While Processing Request") {
		return fmt.Errorf("failed to update DDNS entry: %s", resp)
	}
	return nil
}

package providers

import (
	"fmt"
	"net/http"
	"strings"
	"io/ioutil"
)

type UpdateNoIPRequest struct {
	
	Host          string
	Username      string
	Password      string
	
}

func UpdateNoIP(request interface{}, ipAddr string) error {
	r, ok := request.(*UpdateNoIPRequest)
	if !ok {
		return fmt.Errorf("invalid request type: %T", request)
	}
	return fmt.Errorf("not implemented %s", r.Host)
}
	UsAg := "User-Agent: Plaenkler DDNS-Updater/V0 info@plaenkler.com" 
	urlStr := fmt.Sprintf("http://%s:%s@dynupdate.no-ip.com/nic/update?hostname=%s&myip=%s", r.Username, r.Password, r.Host, ipAddr)
	resp, err := SendHTTPRequest(http.MethodGet, urlStr, UsAg)
	if err != nil {
		return err
	}
	
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	
	ErrMsg := fmt.Sprintf("When not succsessful, please resolve before trying again: %s/%s", err, body)

	fmt.Println(ErrMsg)
	

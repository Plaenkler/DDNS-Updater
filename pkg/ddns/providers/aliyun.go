package providers

import (
	"fmt"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
)

type UpdateAliyunRequest struct {
	Domain       string `json:"Domain"`
	Host         string `json:"Host"`
	AccessKeyID  string `json:"AccessKeyID"`
	AccessSecret string `json:"AccessSecret"`
	Region       string `json:"Region"`
}

func UpdateAliyun(request interface{}, ipAddr string) error {
	r, ok := request.(UpdateAliyunRequest)
	if !ok {
		return fmt.Errorf("invalid request type: %T", request)
	}
	client, err := alidns.NewClientWithAccessKey(r.Region, r.AccessKeyID, r.AccessSecret)
	if err != nil {
		return err
	}
	listRequest := alidns.CreateDescribeDomainRecordsRequest()
	listRequest.Scheme = "https"
	listRequest.DomainName = r.Domain
	listRequest.RRKeyWord = r.Host
	resp, err := client.DescribeDomainRecords(listRequest)
	if err != nil {
		return err
	}
	recordID := ""
	for _, record := range resp.DomainRecords.Record {
		if strings.EqualFold(record.RR, r.Host) {
			recordID = record.RecordId
			break
		}
	}
	if recordID == "" {
		return fmt.Errorf("no record found for host %s", r.Host)
	}
	updateRequest := alidns.CreateUpdateDomainRecordRequest()
	updateRequest.Scheme = "https"
	updateRequest.Value = ipAddr
	updateRequest.Type = "A"
	updateRequest.RR = r.Host
	updateRequest.RecordId = recordID
	_, err = client.UpdateDomainRecord(updateRequest)
	if err != nil {
		return err
	}
	return nil
}

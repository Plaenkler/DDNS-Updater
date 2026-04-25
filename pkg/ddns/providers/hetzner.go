package providers

import (
	"context"
	"fmt"
	"math"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/plaenkler/ddns-updater/pkg/config"
)

type UpdateHetznerRequest struct {
	APIToken     string
	RecordZoneID string
	RecordType   string
	RecordName   string
	RecordTTL    uint32 `json:",string"`
}

func UpdateHetzner(request any, ipAddr string) error {
	r, ok := request.(*UpdateHetznerRequest)
	if !ok {
		return fmt.Errorf("invalid request type: %T", request)
	}

	ctx := context.Background()
	c := hcloud.NewClient(hcloud.WithToken(r.APIToken))

	zone, _, err := c.Zone.Get(ctx, r.RecordZoneID)
	if err != nil {
		return fmt.Errorf("get zone: %w", err)
	}
	if zone == nil {
		return fmt.Errorf("zone not found: %s", r.RecordZoneID)
	}

	rrsetType := hcloud.ZoneRRSetType(r.RecordType)
	rrset, _, err := c.Zone.GetRRSetByNameAndType(ctx, zone, r.RecordName, rrsetType)
	if err != nil {
		return fmt.Errorf("find rrset: %w", err)
	}

	if rrset == nil {
		ttl, err := recordTTL(r.RecordTTL)
		if err != nil {
			return err
		}
		result, _, err := c.Zone.CreateRRSet(ctx, zone, hcloud.ZoneRRSetCreateOpts{
			Name:    r.RecordName,
			Type:    rrsetType,
			TTL:     &ttl,
			Records: []hcloud.ZoneRRSetRecord{{Value: ipAddr}},
		})
		if err != nil {
			return fmt.Errorf("create rrset: %w", err)
		}
		if result.Action != nil {
			if err := c.Action.WaitFor(ctx, result.Action); err != nil {
				return fmt.Errorf("wait for create rrset action: %w", err)
			}
		}
		return nil
	}

	action, _, err := c.Zone.SetRRSetRecords(ctx, rrset, hcloud.ZoneRRSetSetRecordsOpts{
		Records: []hcloud.ZoneRRSetRecord{{Value: ipAddr}},
	})
	if err != nil {
		return fmt.Errorf("set rrset records: %w", err)
	}
	if action != nil {
		if err := c.Action.WaitFor(ctx, action); err != nil {
			return fmt.Errorf("wait for set rrset records action: %w", err)
		}
	}
	return nil
}

// recordTTL returns the TTL to use for a DNS record. When the configured TTL
// is zero, the application update interval is used as a default. Values that
// exceed the maximum safe signed 32-bit integer are rejected with an error to
// prevent silent data corruption on 32-bit platforms.
func recordTTL(configured uint32) (int, error) {
	if configured == 0 {
		configured = config.Get().Interval
	}
	if configured > math.MaxInt32 {
		return 0, fmt.Errorf("record TTL value %d exceeds maximum allowed (%d)", configured, math.MaxInt32)
	}
	return int(configured), nil
}

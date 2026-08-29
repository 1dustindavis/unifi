package unifi

import "fmt"

// APIActiveClientsPath returns currently active clients for a site using the
// Network v2 API. Unlike the legacy stat/sta endpoint, this endpoint includes
// client types such as Teleport connections on current UniFi Network releases.
const APIActiveClientsPath = "/v2/api/site/%s/clients/active"

// GetActiveClients returns the currently active clients reported by the UniFi
// Network v2 API for the requested sites.
func (u *Unifi) GetActiveClients(sites []*Site) ([]*ActiveClient, error) {
	if u == nil {
		return nil, ErrNilUnifi
	}

	data := make([]*ActiveClient, 0)

	for _, site := range sites {
		response := make([]*ActiveClient, 0)

		u.DebugLog("Polling Controller, retrieving active UniFi Clients, site %s", site.SiteName)

		clientPath := fmt.Sprintf(APIActiveClientsPath, site.Name)
		if err := u.GetData(clientPath, &response); err != nil {
			return nil, err
		}

		for _, client := range response {
			client.SourceName = u.URL
			client.SiteName = site.SiteName
		}

		data = append(data, response...)
	}

	return data, nil
}

// ActiveClient defines the common fields returned by the Network v2 active
// clients endpoint. The endpoint includes client types not present in the
// legacy stat/sta response, including Teleport connections.
type ActiveClient struct {
	AssocTime        FlexInt `json:"assoc_time"`
	DisplayName      string  `json:"display_name"`
	ExternalClientID string  `json:"external_client_id"`
	ID               string  `json:"id"`
	IP               string  `json:"ip"`
	LastSeen         FlexInt `json:"last_seen"`
	Mac              string  `json:"mac,omitempty"`
	Name             string  `json:"name"`
	NetworkID        string  `json:"network_id"`
	RxBytes          FlexInt `json:"rx_bytes"`
	RxPackets        FlexInt `json:"rx_packets"`
	SiteID           string  `json:"site_id"`
	Status           string  `json:"status"`
	TokenID          string  `json:"token_id"`
	TxBytes          FlexInt `json:"tx_bytes"`
	TxPackets        FlexInt `json:"tx_packets"`
	Type             string  `json:"type"`
	Uptime           FlexInt `json:"uptime"`
	SiteName         string  `json:"-"`
	SourceName       string  `json:"-"`
}

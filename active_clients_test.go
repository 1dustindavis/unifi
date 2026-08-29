package unifi_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/unpoller/unifi/v6"
)

const activeClientResponse = `[{"assoc_time":1700000000,"display_name":"Teleport Client 02:00:00:00:00:01","external_client_id":"11111111-2222-4333-8444-555555555555","id":"aaaaaaaaaaaaaaaaaaaaaaaa","ip":"192.0.2.10","last_seen":1700000123,"name":"Teleport Client 02:00:00:00:00:01","network_id":"bbbbbbbbbbbbbbbbbbbbbbbb","rx_bytes":123456789,"rx_packets":123456,"site_id":"cccccccccccccccccccccccc","status":"online","token_id":"dddddddddddddddddddddddd","tx_bytes":987654321,"tx_packets":654321,"type":"TELEPORT","uptime":123}]`

func TestParseTeleportActiveClient(t *testing.T) {
	t.Parallel()

	clients := make([]*unifi.ActiveClient, 0)
	err := json.Unmarshal([]byte(activeClientResponse), &clients)

	assert.NoError(t, err)
	if assert.Len(t, clients, 1) {
		client := clients[0]
		assert.Equal(t, "TELEPORT", client.Type)
		assert.Equal(t, "online", client.Status)
		assert.Equal(t, "192.0.2.10", client.IP)
		assert.Equal(t, "bbbbbbbbbbbbbbbbbbbbbbbb", client.NetworkID)
		assert.Equal(t, "11111111-2222-4333-8444-555555555555", client.ExternalClientID)
		assert.Equal(t, "dddddddddddddddddddddddd", client.TokenID)
		assert.Equal(t, float64(1700000000), client.AssocTime.Val)
		assert.Equal(t, float64(1700000123), client.LastSeen.Val)
		assert.Equal(t, float64(123456789), client.RxBytes.Val)
		assert.Equal(t, float64(987654321), client.TxBytes.Val)
		assert.Equal(t, float64(123456), client.RxPackets.Val)
		assert.Equal(t, float64(654321), client.TxPackets.Val)
		assert.Equal(t, float64(123), client.Uptime.Val)
	}
}

func TestGetActiveClients(t *testing.T) {
	t.Parallel()

	requested := make([]string, 0, 2)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requested = append(requested, r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(activeClientResponse))
	}))
	defer server.Close()

	u := &unifi.Unifi{
		Client: server.Client(),
		Config: &unifi.Config{URL: server.URL},
	}
	sites := []*unifi.Site{
		{Name: "default", SiteName: "Default (default)"},
		{Name: "remote", SiteName: "Remote (remote)"},
	}

	clients, err := u.GetActiveClients(sites)
	require.NoError(t, err)
	require.Len(t, clients, 2)
	assert.Equal(t, []string{
		"/v2/api/site/default/clients/active",
		"/v2/api/site/remote/clients/active",
	}, requested)
	assert.Equal(t, "Default (default)", clients[0].SiteName)
	assert.Equal(t, server.URL, clients[0].SourceName)
	assert.Equal(t, "Remote (remote)", clients[1].SiteName)
	assert.Equal(t, server.URL, clients[1].SourceName)
	assert.Equal(t, "TELEPORT", clients[0].Type)
	assert.Equal(t, "online", clients[0].Status)
}

func TestGetActiveClientsSiteValidation(t *testing.T) {
	t.Parallel()

	u := &unifi.Unifi{Config: &unifi.Config{}}

	_, err := u.GetActiveClientsSite(nil)
	assert.ErrorIs(t, err, unifi.ErrNoSiteProvided)

	_, err = u.GetActiveClientsSite(&unifi.Site{})
	assert.ErrorIs(t, err, unifi.ErrNoSiteProvided)

	var nilClient *unifi.Unifi
	_, err = nilClient.GetActiveClientsSite(&unifi.Site{Name: "default"})
	assert.True(t, errors.Is(err, unifi.ErrNilUnifi))
}

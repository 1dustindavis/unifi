package unifi_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/unpoller/unifi/v6"
)

func TestParseTeleportActiveClient(t *testing.T) {
	t.Parallel()

	const response = `[{"assoc_time":1700000000,"display_name":"Teleport Client 02:00:00:00:00:01","external_client_id":"11111111-2222-4333-8444-555555555555","id":"aaaaaaaaaaaaaaaaaaaaaaaa","ip":"192.0.2.10","last_seen":1700000123,"name":"Teleport Client 02:00:00:00:00:01","network_id":"bbbbbbbbbbbbbbbbbbbbbbbb","rx_bytes":123456789,"rx_packets":123456,"site_id":"cccccccccccccccccccccccc","status":"online","token_id":"dddddddddddddddddddddddd","tx_bytes":987654321,"tx_packets":654321,"type":"TELEPORT","uptime":123}]`

	clients := make([]*unifi.ActiveClient, 0)
	err := json.Unmarshal([]byte(response), &clients)

	assert.NoError(t, err)
	if assert.Len(t, clients, 1) {
		client := clients[0]
		assert.Equal(t, "TELEPORT", client.Type)
		assert.Equal(t, "online", client.Status)
		assert.Equal(t, "192.0.2.10", client.IP)
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

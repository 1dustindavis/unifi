package mocks

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/unpoller/unifi/v6"
)

// GetActiveClients returns mocked active client data.
func (m *MockUnifi) GetActiveClients(_ []*unifi.Site) ([]*unifi.ActiveClient, error) {
	return m.GetActiveClientsSite(nil)
}

// GetActiveClientsSite returns mocked active client data for a site.
func (m *MockUnifi) GetActiveClientsSite(_ *unifi.Site) ([]*unifi.ActiveClient, error) {
	results := make([]*unifi.ActiveClient, numItemsMocked)

	for i := 0; i < numItemsMocked; i++ {
		var client unifi.ActiveClient

		if err := gofakeit.Struct(&client); err != nil {
			return results, err
		}

		results[i] = &client
	}

	return results, nil
}

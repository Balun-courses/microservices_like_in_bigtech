//go:build integration
// +build integration

package warehouses_management_system

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_ReserveStocks(t *testing.T) {
	// prepare

	r := &Client{}

	t.Run("Test 1.", func(t *testing.T) {
		if err := r.ReserveStocks(context.Background(), 1, nil); err == nil {
			assert.NotNil(t, err)
		}
	})
}

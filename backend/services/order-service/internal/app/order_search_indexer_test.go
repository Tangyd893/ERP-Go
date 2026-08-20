package app

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Tangyd893/ERP-Go/backend/services/order-service/internal/domain"
	opensearchx "github.com/Tangyd893/ERP-Go/backend/shared/opensearch"
	"github.com/stretchr/testify/require"
)

func TestOrderSearchIndexer_IndexOrder(t *testing.T) {
	var gotPath string
	var gotDoc opensearchx.OrderDoc
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		_ = json.NewDecoder(r.Body).Decode(&gotDoc)
		w.WriteHeader(http.StatusCreated)
	}))
	defer srv.Close()

	t.Setenv("OPENSEARCH_URL", srv.URL)
	osClient := opensearchx.NewFromEnv()
	require.NotNil(t, osClient)

	idx := NewOrderSearchIndexer(osClient)
	err := idx.IndexOrder(context.Background(), &domain.SalesOrder{
		ID: "o1", TenantID: "t1", PlatformOrderNo: "P-1",
		Status: domain.OrderApproved, TotalAmount: 12.5,
	})
	require.NoError(t, err)
	require.Equal(t, "/erp-orders/_doc/o1", gotPath)
	require.Equal(t, "o1", gotDoc.OrderID)
	require.Equal(t, "approved", gotDoc.Status)
	require.Equal(t, 12.5, gotDoc.TotalAmount)
}

func TestOrderSearchIndexer_NilSafe(t *testing.T) {
	var idx *OrderSearchIndexer
	require.NoError(t, idx.IndexOrder(context.Background(), &domain.SalesOrder{ID: "x"}))
	require.NoError(t, NewOrderSearchIndexer(nil).IndexOrder(context.Background(), &domain.SalesOrder{ID: "x"}))
}

package opensearchx

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIndexDoc_AndEnsureIndex(t *testing.T) {
	var putCount, headCount int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodHead:
			headCount++
			w.WriteHeader(http.StatusNotFound)
		case r.Method == http.MethodPut && r.URL.Path == "/erp-orders":
			w.WriteHeader(http.StatusOK)
		case r.Method == http.MethodPut:
			putCount++
			var doc OrderDoc
			_ = json.NewDecoder(r.Body).Decode(&doc)
			require.Equal(t, "o1", doc.OrderID)
			w.WriteHeader(http.StatusCreated)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()

	c := &Client{baseURL: srv.URL, httpClient: srv.Client()}
	require.NoError(t, c.EnsureIndex(context.Background(), IndexOrders))
	require.Equal(t, 1, headCount)
	require.NoError(t, c.IndexDoc(context.Background(), IndexOrders, "o1", OrderDoc{
		OrderID: "o1", TenantID: "default", Status: "pending",
	}))
	require.Equal(t, 1, putCount)
}

package sql

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"google.golang.org/api/sqladmin/v1"

	"github.com/zopdev/zopdev/api/resources/providers/models"
)

func getServer(t *testing.T, resp any, isError bool) *httptest.Server {
	t.Helper()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if isError {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		err := json.NewEncoder(w).Encode(resp)
		if err != nil {
			http.Error(w, "unable to marshal response", http.StatusInternalServerError)
			return
		}
	}))

	return srv
}

func Test_GetAllInstances(t *testing.T) {
	resp := &sqladmin.InstancesListResponse{
		Items: []*sqladmin.DatabaseInstance{
			{Name: "test-instance", Project: "test-project"},
		},
	}

	srv := getServer(t, resp, false)
	defer srv.Close()

	instSvc, err := sqladmin.NewService(context.Background(), option.WithoutAuthentication(), option.WithEndpoint(srv.URL))
	require.NoError(t, err)

	c := Client{SQL: instSvc.Instances}

	admin, err := c.GetAllInstances(nil, "test-project")

	require.NoError(t, err)
	require.NotNil(t, admin)
	assert.Equal(t, []models.Instance{{Name: "test-instance", Type: "SQL", ProviderID: "test-project"}}, admin)
}

func Test_GetAllInstances_Error(t *testing.T) {
	srv := getServer(t, nil, true)
	defer srv.Close()

	expected := &googleapi.Error{
		Code: http.StatusInternalServerError,
		Body: "Internal server error\n",
	}

	instSvc, err := sqladmin.NewService(context.Background(), option.WithoutAuthentication(), option.WithEndpoint(srv.URL))
	require.NoError(t, err)

	c := Client{SQL: instSvc.Instances}

	admin, err := c.GetAllInstances(nil, "test-project")

	require.Error(t, err)
	require.Nil(t, admin)
	assert.Equal(t, expected.Error(), err.Error())
}

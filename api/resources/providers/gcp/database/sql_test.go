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

	"github.com/zopdev/zopdev/api/resources/models"
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
			{Name: "test-instance1", Project: "test-project", Settings: &sqladmin.Settings{ActivationPolicy: ALWAYS}},
			{Name: "test-instance2", Project: "test-project", Settings: &sqladmin.Settings{ActivationPolicy: NEVER}},
			{Name: "test-instance3", Project: "test-project", Settings: &sqladmin.Settings{ActivationPolicy: "ON_DEMAND"}},
		}}
	result := []models.Instance{
		{Name: "test-instance1", UID: "test-project/test-instance1", Type: "SQL", Status: RUNNING},
		{Name: "test-instance2", UID: "test-project/test-instance2", Type: "SQL", Status: SUSPENDED},
		{Name: "test-instance3", UID: "test-project/test-instance3", Type: "SQL", Status: SUSPENDED},
	}

	srv := getServer(t, resp, false)
	defer srv.Close()

	instSvc, err := sqladmin.NewService(context.Background(), option.WithoutAuthentication(), option.WithEndpoint(srv.URL))
	require.NoError(t, err)

	c := Client{SQL: instSvc.Instances}

	admin, err := c.GetAllInstances(nil, "test-project")

	require.NoError(t, err)
	require.NotNil(t, admin)
	assert.Equal(t, result, admin)
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

func TestClient_StartInstance(t *testing.T) {
	// Success case
	srv1 := getServer(t, nil, false)
	defer srv1.Close()

	instSvc, err := sqladmin.NewService(context.Background(), option.WithoutAuthentication(), option.WithEndpoint(srv1.URL))
	require.NoError(t, err)

	c := Client{SQL: instSvc.Instances}

	err = c.StartInstance(nil, "test-project", "test-instance")
	require.NoError(t, err)

	// Error case
	srv2 := getServer(t, nil, true)
	defer srv2.Close()

	instSvc, err = sqladmin.NewService(context.Background(), option.WithoutAuthentication(), option.WithEndpoint(srv2.URL))
	require.NoError(t, err)

	c = Client{SQL: instSvc.Instances}

	err = c.StartInstance(nil, "test-project", "test-instance")
	require.Error(t, err)
	assert.Equal(t, &InternalServerError{}, err)
}

func TestClient_StopInstance(t *testing.T) {
	// Success case
	srv1 := getServer(t, nil, false)
	defer srv1.Close()

	instSvc, err := sqladmin.NewService(context.Background(), option.WithoutAuthentication(), option.WithEndpoint(srv1.URL))
	require.NoError(t, err)

	c := Client{SQL: instSvc.Instances}

	err = c.StopInstance(nil, "test-project", "test-instance")
	require.NoError(t, err)

	// Error case
	srv2 := getServer(t, nil, true)
	defer srv2.Close()

	instSvc, err = sqladmin.NewService(context.Background(), option.WithoutAuthentication(), option.WithEndpoint(srv2.URL))
	require.NoError(t, err)

	c = Client{SQL: instSvc.Instances}

	err = c.StopInstance(nil, "test-project", "test-instance")
	require.Error(t, err)
}

func Test_getError(t *testing.T) {
	err := &googleapi.Error{
		Code:    http.StatusConflict,
		Message: "Conflict error",
	}

	errRes := getError(err)
	expected := &ErrConflict{
		Message: "Conflict error",
	}

	assert.Equal(t, expected, errRes)

	err = &googleapi.Error{
		Code:    http.StatusInternalServerError,
		Message: "Internal server error",
	}

	errRes = getError(err)
	expected2 := &InternalServerError{}

	assert.Equal(t, expected2, errRes)

	errRes = getError(nil)
	assert.NoError(t, errRes)
}

func Test_Errors(t *testing.T) {
	e := &ErrConflict{Message: "Conflict error"}
	assert.Equal(t, "Conflict error", e.Error())
	assert.Equal(t, http.StatusConflict, e.StatusCode())

	e2 := &InternalServerError{}
	assert.Equal(t, "Internal server error!", e2.Error())
}

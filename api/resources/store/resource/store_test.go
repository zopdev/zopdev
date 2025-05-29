package resource

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/container"

	"github.com/zopdev/zopdev/api/resources/models"
)

func TestStore_InsertResource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockContainer, mocks := container.NewMockContainer(t)

	mockInput := &models.Resource{
		UID:          "zopdev/test-instance",
		Name:         "test-instance",
		Status:       "RUNNING",
		CloudAccount: models.CloudAccount{ID: 1, Type: "GCP"},
		Type:         "SQL",
	}
	ctx := &gofr.Context{Container: mockContainer, Context: context.Background()}
	store := New()

	testCases := []struct {
		name      string
		input     *models.Resource
		expErr    error
		mockCalls func()
	}{
		{
			name:  "Successful Insert",
			input: mockInput,
			mockCalls: func() {
				mocks.SQL.Sqlmock.ExpectExec(
					`INSERT INTO resources (resource_uid, name, state, cloud_account_id, cloud_provider, resource_type) VALUES (?, ?, ?, ?, ?, ?)`).
					WithArgs(mockInput.UID, mockInput.Name, mockInput.Status, mockInput.CloudAccount.ID, mockInput.CloudAccount.Type, mockInput.Type).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name:   "Insert Error",
			input:  mockInput,
			expErr: assert.AnError,
			mockCalls: func() {
				mocks.SQL.Sqlmock.ExpectExec(
					`INSERT INTO resources (resource_uid, name, state, cloud_account_id, cloud_provider, resource_type) VALUES (?, ?, ?, ?, ?, ?)`).
					WithArgs(mockInput.UID, mockInput.Name, mockInput.Status, mockInput.CloudAccount.ID, mockInput.CloudAccount.Type, mockInput.Type).
					WillReturnError(assert.AnError)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockCalls()

			err := store.InsertResource(ctx, tc.input)

			assert.Equal(t, tc.expErr, err)
		})
	}
}

func TestStore_GetResources(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockContainer, mocks := container.NewMockContainer(t)
	mockTime := time.Now()
	query := `SELECT id, resource_uid, name, state, cloud_account_id, 
       cloud_provider, resource_type, created_at, updated_at 
		FROM resources WHERE cloud_account_id = ? AND resource_type IN (?, ?) ORDER BY resource_uid`
	ctx := &gofr.Context{Container: mockContainer, Context: context.Background()}
	store := New()

	testCases := []struct {
		name         string
		cloudAccID   int64
		resourceType []string
		expErr       error
		expResp      []models.Resource
		mockCalls    func()
	}{
		{
			name:         "Fetch All Resources",
			cloudAccID:   123,
			resourceType: []string{"SQL", "VM"},
			mockCalls: func() {
				mocks.SQL.Sqlmock.ExpectQuery(query).WithArgs(123, "SQL", "VM").
					WillReturnRows(sqlmock.NewRows([]string{"id", "resource_uid", "name", "state",
						"cloud_account_id", "cloud_provider", "resource_type", "created_at", "updated_at"}).
						AddRow(1, "zopdev/sql-instance-1", "sql-instance-1", "RUNNING", 123, "GCP",
							"SQL", mockTime, mockTime).
						AddRow(2, "zopdev/vm-instance-1", "vm-instance-1", "STOPPED", 123, "GCP",
							"VM", mockTime, mockTime))
			},
			expResp: []models.Resource{
				{ID: 1, UID: "zopdev/sql-instance-1", CloudAccount: models.CloudAccount{ID: 123, Type: "GCP"},
					Type: "SQL", Name: "sql-instance-1", Status: "RUNNING", CreatedAt: mockTime, UpdatedAt: mockTime},
				{ID: 2, UID: "zopdev/vm-instance-1", CloudAccount: models.CloudAccount{ID: 123, Type: "GCP"},
					Type: "VM", Name: "vm-instance-1", Status: "STOPPED", CreatedAt: mockTime, UpdatedAt: mockTime},
			},
		},
		{
			name:         "Fetch Resources with No Type",
			cloudAccID:   123,
			resourceType: []string{"SQL"},
			expErr:       nil,
			expResp: []models.Resource{
				{ID: 1, Name: "sql-instance-1", Type: "SQL", Status: "RUNNING", CreatedAt: mockTime, UpdatedAt: mockTime,
					UID: "zopdev/sql-instance-1", CloudAccount: models.CloudAccount{ID: 123, Type: "GCP"}},
			},
			mockCalls: func() {
				mocks.SQL.Sqlmock.ExpectQuery(`SELECT id, resource_uid, name, state, cloud_account_id, 
       cloud_provider, resource_type, created_at, updated_at 
		FROM resources WHERE cloud_account_id = ? AND resource_type IN (?) ORDER BY resource_uid`).WithArgs(123, "SQL").
					WillReturnRows(sqlmock.NewRows([]string{"id", "resource_uid", "name", "state",
						"cloud_account_id", "cloud_provider", "resource_type", "created_at", "updated_at"}).
						AddRow(1, "zopdev/sql-instance-1", "sql-instance-1", "RUNNING", 123, "GCP",
							"SQL", mockTime, mockTime))
			},
		},
		{
			name:         "Fetch Resources with Empty Type",
			cloudAccID:   123,
			resourceType: []string{"SQL", "VM"},
			expErr:       assert.AnError,
			mockCalls: func() {
				mocks.SQL.Sqlmock.ExpectQuery(query).WithArgs(123, "SQL", "VM").
					WillReturnError(assert.AnError)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockCalls()

			resources, err := store.GetResources(ctx, tc.cloudAccID, tc.resourceType)

			assert.Equal(t, tc.expErr, err)
			assert.Equal(t, tc.expResp, resources)
		})
	}
}

func TestStore_UpdateResource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockContainer, mocks := container.NewMockContainer(t)
	ctx := &gofr.Context{Container: mockContainer, Context: context.Background()}
	store := New()

	testCases := []struct {
		name      string
		resource  *models.Resource
		expErr    error
		mockCalls func()
	}{
		{
			name: "Successful Update",
			resource: &models.Resource{
				ID:           1,
				Name:         "test-instance",
				Status:       "RUNNING",
				CloudAccount: models.CloudAccount{ID: 1, Type: "GCP"},
				Type:         "SQL",
			},
			mockCalls: func() {
				mocks.SQL.Sqlmock.ExpectExec(`UPDATE resources SET state = ? WHERE id = ?`).
					WithArgs("RUNNING", 1).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name: "Update Error",
			resource: &models.Resource{
				ID:           2,
				Name:         "test-instance-2",
				Status:       "STOPPED",
				CloudAccount: models.CloudAccount{ID: 2, Type: "AWS"},
				Type:         "VM",
			},
			expErr: assert.AnError,
			mockCalls: func() {
				mocks.SQL.Sqlmock.ExpectExec(`UPDATE resources SET state = ? WHERE id = ?`).
					WithArgs("STOPPED", 2).WillReturnError(assert.AnError)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockCalls()

			err := store.UpdateResource(ctx, tc.resource)

			assert.Equal(t, tc.expErr, err)
		})
	}
}

func TestStore_RemoveResource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockContainer, mocks := container.NewMockContainer(t)
	ctx := &gofr.Context{Container: mockContainer, Context: context.Background()}
	store := New()

	testCases := []struct {
		name       string
		resourceID int64
		expErr     error
		mockCalls  func()
	}{
		{
			name:       "Successful Removal",
			resourceID: 1,
			mockCalls: func() {
				mocks.SQL.Sqlmock.ExpectExec(`DELETE FROM resources WHERE id = ?`).
					WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name:       "Removal Error",
			resourceID: 2,
			expErr:     assert.AnError,
			mockCalls: func() {
				mocks.SQL.Sqlmock.ExpectExec(`DELETE FROM resources WHERE id = ?`).
					WithArgs(2).WillReturnError(assert.AnError)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockCalls()

			err := store.RemoveResource(ctx, tc.resourceID)

			assert.Equal(t, tc.expErr, err)
		})
	}
}

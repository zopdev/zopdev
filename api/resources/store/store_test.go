package store

import (
	"context"
	"database/sql"
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

	mockInput := &models.Instance{
		UID:          "zopdev/test-instance",
		Name:         "test-instance",
		Status:       "RUNNING",
		CloudAccount: models.CloudAccount{ID: 1, Type: "GCP"},
		Settings: models.Settings{
			"region": "us-central1",
			"zone":   "us-central1-a",
		},
		Type:   "SQL",
		Region: "us-central1",
	}
	ctx := &gofr.Context{Container: mockContainer, Context: context.Background()}
	store := New()

	testCases := []struct {
		name      string
		input     *models.Instance
		expErr    error
		mockCalls func()
	}{
		{
			name:  "Successful Insert",
			input: mockInput,
			mockCalls: func() {
				mocks.SQL.Sqlmock.ExpectExec(
					`INSERT INTO resources (resource_uid, name, state, cloud_account_id, cloud_provider, resource_type, 
settings, region) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`).
					WithArgs(mockInput.UID, mockInput.Name, mockInput.Status, mockInput.CloudAccount.ID,
						mockInput.CloudAccount.Type, mockInput.Type, mockInput.Settings, mockInput.Region).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name:   "Insert Error",
			input:  mockInput,
			expErr: assert.AnError,
			mockCalls: func() {
				mocks.SQL.Sqlmock.ExpectExec(
					`INSERT INTO resources (resource_uid, name, state, cloud_account_id, cloud_provider, resource_type,  
settings, region) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`).
					WithArgs(mockInput.UID, mockInput.Name, mockInput.Status, mockInput.CloudAccount.ID,
						mockInput.CloudAccount.Type, mockInput.Type, mockInput.Settings, mockInput.Region).
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
	settings := models.Settings{"region": "us-central1", "zone": "us-central1-a"}
	query := `SELECT id, resource_uid, name, state, cloud_account_id, 
       cloud_provider, resource_type, created_at, updated_at, settings, region
		FROM resources WHERE cloud_account_id = ? AND resource_type IN (?, ?) ORDER BY resource_uid`
	ctx := &gofr.Context{Container: mockContainer, Context: context.Background()}
	store := New()

	testCases := []struct {
		name         string
		cloudAccID   int64
		resourceType []string
		expErr       error
		expResp      []models.Instance
		mockCalls    func()
	}{
		{
			name:         "Fetch All Resources",
			cloudAccID:   123,
			resourceType: []string{"SQL", "VM"},
			mockCalls: func() {
				mocks.SQL.Sqlmock.ExpectQuery(query).WithArgs(123, "SQL", "VM").
					WillReturnRows(sqlmock.NewRows([]string{"id", "resource_uid", "name", "state",
						"cloud_account_id", "cloud_provider", "resource_type", "created_at", "updated_at", "settings", "region"}).
						AddRow(1, "zopdev/sql-instance-1", "sql-instance-1", "RUNNING", 123, "GCP",
							"SQL", mockTime, mockTime, &settings, "us-central1").
						AddRow(2, "zopdev/vm-instance-1", "vm-instance-1", "STOPPED", 123, "GCP",
							"VM", mockTime, mockTime, &settings, "us-central1"))
			},
			expResp: []models.Instance{
				{ID: 1, UID: "zopdev/sql-instance-1", CloudAccount: models.CloudAccount{ID: 123, Type: "GCP"},
					Type: "SQL", Name: "sql-instance-1", Status: "RUNNING", CreatedAt: mockTime, UpdatedAt: mockTime,
					Settings: settings, Region: "us-central1"},
				{ID: 2, UID: "zopdev/vm-instance-1", CloudAccount: models.CloudAccount{ID: 123, Type: "GCP"},
					Type: "VM", Name: "vm-instance-1", Status: "STOPPED", CreatedAt: mockTime, UpdatedAt: mockTime,
					Settings: settings, Region: "us-central1"},
			},
		},
		{
			name:         "Fetch Resources with No Type",
			cloudAccID:   123,
			resourceType: []string{"SQL"},
			expErr:       nil,
			expResp: []models.Instance{
				{ID: 1, Name: "sql-instance-1", Type: "SQL", Status: "RUNNING", CreatedAt: mockTime, UpdatedAt: mockTime, Region: "us-central1",
					Settings: settings, UID: "zopdev/sql-instance-1", CloudAccount: models.CloudAccount{ID: 123, Type: "GCP"}},
			},
			mockCalls: func() {
				mocks.SQL.Sqlmock.ExpectQuery(`SELECT id, resource_uid, name, state, cloud_account_id, 
       cloud_provider, resource_type, created_at, updated_at, settings, region
		FROM resources WHERE cloud_account_id = ? AND resource_type IN (?) ORDER BY resource_uid`).WithArgs(123, "SQL").
					WillReturnRows(sqlmock.NewRows([]string{"id", "resource_uid", "name", "state", "cloud_account_id",
						"cloud_provider", "resource_type", "created_at", "updated_at", "settings", "region"}).
						AddRow(1, "zopdev/sql-instance-1", "sql-instance-1", "RUNNING", 123, "GCP",
							"SQL", mockTime, mockTime, &settings, "us-central1"))
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
		status    string
		id        int64
		expErr    error
		mockCalls func()
	}{
		{
			name:   "Successful Update",
			id:     1,
			status: "RUNNING",
			mockCalls: func() {
				mocks.SQL.Sqlmock.ExpectExec(`UPDATE resources SET state = ? WHERE id = ?`).
					WithArgs("RUNNING", 1).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name:   "Update Error",
			id:     2,
			status: "STOPPED",
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

			err := store.UpdateStatus(ctx, tc.status, tc.id)

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

func TestStore_GetResourceByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTime := time.Now()
	mockContainer, mocks := container.NewMockContainer(t)
	settings := models.Settings{"region": "us-central1", "zone": "us-central1-a"}
	ctx := &gofr.Context{Container: mockContainer, Context: context.Background()}
	query := `SELECT id, resource_uid, name, state, cloud_account_id,
		cloud_provider, resource_type, created_at, updated_at, settings, region
	FROM resources WHERE id = ?`
	mockResp := &models.Instance{
		ID:     1,
		UID:    "zopdev/sql-instance-1",
		Name:   "sql-instance-1",
		Status: "RUNNING",
		Type:   "SQL",
		CloudAccount: models.CloudAccount{
			ID:   123,
			Type: "GCP",
		},
		CreatedAt: mockTime,
		UpdatedAt: mockTime,
		Settings:  settings,
		Region:    "us-central1",
	}
	store := New()

	testCases := []struct {
		name      string
		inputID   int64
		expResp   *models.Instance
		expErr    error
		mockCalls func()
	}{
		{
			name:    "Successful Fetch by ID",
			inputID: 1,
			expResp: mockResp,
			mockCalls: func() {
				mocks.SQL.ExpectQuery(query).WithArgs(int64(1)).
					WillReturnRows(sqlmock.NewRows([]string{"id", "resource_uid", "name", "state", "cloud_account_id",
						"cloud_provider", "resource_type", "created_at", "updated_at", "settings", "region"}).
						AddRow(1, "zopdev/sql-instance-1", "sql-instance-1", "RUNNING", 123, "GCP",
							"SQL", mockTime, mockTime, &settings, "us-central1").
						AddRow(2, "zopdev/vm-instance-1", "vm-instance-1", "STOPPED", 123, "GCP",
							"VM", mockTime, mockTime, &settings, "us-central1"))
			},
		},
		{
			name:    "Fetch by Non-Existent ID",
			inputID: 999,
			expResp: nil,
			expErr:  sql.ErrNoRows,
			mockCalls: func() {
				mocks.SQL.ExpectQuery(query).WithArgs(int64(999)).
					WillReturnError(sql.ErrNoRows)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockCalls()

			resource, err := store.GetResourceByID(ctx, tc.inputID)

			assert.Equal(t, tc.expErr, err)
			assert.Equal(t, tc.expResp, resource)
		})
	}
}

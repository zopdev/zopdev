package resourcegroup

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/container"

	"github.com/zopdev/zopdev/api/resources/models"
)

func TestStore_GetAllResourceGroups(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockContainer, mocks := container.NewMockContainer(t)
	query := `SELECT id, name, description FROM resource_groups`
	rows := sqlmock.NewRows([]string{"id", "name", "description"}).
		AddRow(1, "group1", "Description for group 1").
		AddRow(2, "group2", "Description for group 2")

	testCases := []struct {
		name      string
		expErr    error
		expRes    []models.ResourceGroup
		mockSetup func()
	}{
		{
			name: "GetAllResourceGroups_Success",
			expRes: []models.ResourceGroup{
				{ID: 1, Name: "group1", Description: "Description for group 1"},
				{ID: 2, Name: "group2", Description: "Description for group 2"},
			},
			mockSetup: func() {
				mocks.SQL.ExpectQuery(query).
					WillReturnRows(rows)
			},
		},
		{
			name:   "GetAllResourceGroups_Error",
			expErr: sqlmock.ErrCancelled,
			mockSetup: func() {
				mocks.SQL.ExpectQuery(query).
					WillReturnError(sqlmock.ErrCancelled)
			},
		},
		{
			name:   "GetAllResourceGroups_NoRows",
			expRes: nil,
			mockSetup: func() {
				mocks.SQL.ExpectQuery(query).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description"}))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()

			store := New()
			ctx := &gofr.Context{
				Context:   context.Background(),
				Container: mockContainer,
			}

			res, err := store.GetAllResourceGroups(ctx, 1)

			assert.Equal(t, tc.expErr, err)
			assert.Equal(t, tc.expRes, res)
		})
	}
}

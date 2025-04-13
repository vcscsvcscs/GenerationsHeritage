package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	memgraphMock "github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/internal/memgraph/mock"
)

func TestCouldManagePerson(t *testing.T) {
	ctx := context.Background()
	mockSession := new(memgraphMock.SessionWithContext)

	t.Run("adminId equals XUserID", func(t *testing.T) {
		err := CouldManagePerson(ctx, mockSession, 1, 2, 2)
		require.NoError(t, err)
	})

	t.Run("adminId not equal XUserID, calls CouldManagePersonUnknownAdmin", func(t *testing.T) {
		mockSession.On("ExecuteRead", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil).Once()

		err := CouldManagePerson(ctx, mockSession, 1, 2, 3)
		require.NoError(t, err)

		mockSession.AssertExpectations(t)
	})
}

func TestCouldManagePersonUnknownAdmin(t *testing.T) {
	ctx := context.Background()
	mockSession := new(memgraphMock.SessionWithContext)

	t.Run("userId equals XUserID", func(t *testing.T) {
		err := CouldManagePersonUnknownAdmin(ctx, mockSession, 1, 1)
		require.NoError(t, err)
	})

	t.Run("userId not equal XUserID, ExecuteRead returns no error", func(t *testing.T) {
		mockSession.On("ExecuteRead", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil).Once()

		err := CouldManagePersonUnknownAdmin(ctx, mockSession, 1, 2)
		require.NoError(t, err)

		mockSession.AssertExpectations(t)
	})

	t.Run("userId not equal XUserID, ExecuteRead returns error", func(t *testing.T) {
		mockSession.On("ExecuteRead", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("some error")).Once()

		err := CouldManagePersonUnknownAdmin(ctx, mockSession, 1, 2)
		require.Error(t, err)
		require.EqualError(t, err, "some error")

		mockSession.AssertExpectations(t)
	})
}

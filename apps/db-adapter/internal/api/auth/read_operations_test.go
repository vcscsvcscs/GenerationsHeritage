package auth

import (
	"fmt"
	"sync"
	"testing"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	memgraphMock "github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/internal/memgraph/mock"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/pkg/api"
)

func TestCouldSeePersonsProfile(t *testing.T) {
	ctx := t.Context()

	t.Run("User can manage person", func(t *testing.T) {
		mockResult := new(memgraphMock.Result)
		mockResult.On("Single", mock.Anything).Return(&neo4j.Record{}, nil)
		mockSession := new(memgraphMock.SessionWithContext)
		mockSession.On("ExecuteRead", mock.Anything, mock.Anything, mock.Anything).Return(mockResult, nil).Once()
		err := CouldSeePersonsProfile(ctx, mockSession, 1, 2)
		require.NoError(t, err)
	})

	t.Run("User cannot manage person but is in family tree", func(t *testing.T) {
		mockResult := new(memgraphMock.Result)
		mockResult.On("Single", mock.Anything).Return(nil, fmt.Errorf("no permission"))
		mockSession := &memgraphMock.SessionWithContext{
			ReturnOnce: &sync.Once{},
		}
		mockSession.On("ExecuteRead", mock.Anything, mock.Anything, mock.Anything).Return(mockResult, nil, map[string]any{
			"people": []api.OptimizedPersonNode{
				{Id: api.IntPtr(1)},
			},
		}, nil)

		err := CouldSeePersonsProfile(ctx, mockSession, 1, 3)
		require.NoError(t, err)
		mockSession.AssertExpectations(t)
	})

	t.Run("User not in family tree", func(t *testing.T) {
		mockSession := memgraphMock.SessionWithContext{
			ReturnOnce: &sync.Once{},
		}
		mockSession.On("ExecuteRead", mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("no permission"), map[string]any{
			"people": []api.OptimizedPersonNode{
				{Id: api.IntPtr(4)},
			},
		}, nil)

		err := CouldSeePersonsProfile(ctx, &mockSession, 1, 3)
		require.Error(t, err)
		require.EqualError(t, err, "user 3 does not have permission to see user 1")
		mockSession.AssertExpectations(t)
	})

	t.Run("Error during ExecuteRead", func(t *testing.T) {
		mockSession := &memgraphMock.SessionWithContext{
			ReturnOnce: &sync.Once{},
		}
		mockSession.On("ExecuteRead", mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("read error"), nil, fmt.Errorf("read error"))

		err := CouldSeePersonsProfile(ctx, mockSession, 1, 3)
		require.Error(t, err)
		require.EqualError(t, err, "read error")
		mockSession.AssertExpectations(t)
	})

	t.Run("Invalid result format from ExecuteRead", func(t *testing.T) {
		mockSession := &memgraphMock.SessionWithContext{
			ReturnOnce: &sync.Once{},
		}
		mockSession.On("ExecuteRead", mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("invalid"), "invalid", nil)

		err := CouldSeePersonsProfile(ctx, mockSession, 1, 3)
		require.Error(t, err)
		require.EqualError(t, err, "could not convert result to map[string]any")
		mockSession.AssertExpectations(t)
	})

	t.Run("Invalid people format in result", func(t *testing.T) {
		mockSession := &memgraphMock.SessionWithContext{
			ReturnOnce: &sync.Once{},
		}
		mockSession.On("ExecuteRead", mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("invalid"), map[string]any{
			"people": "invalid",
		}, nil)

		err := CouldSeePersonsProfile(ctx, mockSession, 1, 3)
		require.Error(t, err)
		require.EqualError(t, err, "could not convert people to []map[string]any: unexpected type: string")
		mockSession.AssertExpectations(t)
	})
}

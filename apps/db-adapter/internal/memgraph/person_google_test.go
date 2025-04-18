package memgraph

import (
	"context"
	"errors"
	"testing"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	tmock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/internal/memgraph/mock"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/pkg/api"
)

func TestGetPersonByGoogleId(t *testing.T) {
	tests := []struct {
		mockRunResult any
		mockRunError  error
		mockSingleErr error
		expectedRes   any
		mockSingle    *neo4j.Record
		name          string
		expectedErr   bool
	}{
		{
			name:          "Success",
			mockRunResult: new(mock.Result),
			mockRunError:  nil,
			mockSingle: &neo4j.Record{
				Values: []any{"value"},
				Keys:   []string{"key"},
			},
			mockSingleErr: nil,
			expectedErr:   false,
			expectedRes:   map[string]any{"key": "value"},
		},
		{
			name:          "RunError",
			mockRunResult: nil,
			mockRunError:  errors.New("run error"),
			mockSingle:    nil,
			mockSingleErr: nil,
			expectedErr:   true,
			expectedRes:   nil,
		},
		{
			name:          "SingleError",
			mockRunResult: new(mock.Result),
			mockRunError:  nil,
			mockSingle:    nil,
			mockSingleErr: errors.New("single error"),
			expectedErr:   true,
			expectedRes:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockTx := new(mock.Transaction)

			mockTx.On("Run", ctx, GetPersonByGoogleIdCypherQuery, tmock.Anything).Return(tt.mockRunResult, tt.mockRunError)
			if tt.mockRunResult != nil {
				mockResult := tt.mockRunResult.(*mock.Result)
				mockResult.On("Single", ctx).Return(tt.mockSingle, tt.mockSingleErr)
			}

			work := GetPersonByGoogleId(ctx, "test-google-id")
			result, err := work(mockTx)

			if tt.expectedErr {
				require.Error(t, err)
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				require.Equal(t, tt.expectedRes, result)
			}

			mockTx.AssertExpectations(t)
			if tt.mockRunResult != nil {
				tt.mockRunResult.(*mock.Result).AssertExpectations(t)
			}
		})
	}
}

func TestUpdatePersonByGoogleID(t *testing.T) {
	tests := []struct {
		mockRunError error
		name         string
		expectedErr  bool
	}{
		{
			name:         "Success",
			mockRunError: nil,
			expectedErr:  false,
		},
		{
			name:         "RunError",
			mockRunError: errors.New("run error"),
			expectedErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockTx := new(mock.Transaction)

			mockTx.On("Run", ctx, UpdatePersonCypherQuery).Return(nil, tt.mockRunError)

			personProps := &api.PersonProperties{}
			work := UpdatePersonByInviteCode(ctx, "test-person-id", personProps)
			result, err := work(mockTx)

			if tt.expectedErr {
				require.Error(t, err)
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
			}

			mockTx.AssertExpectations(t)
		})
	}
}

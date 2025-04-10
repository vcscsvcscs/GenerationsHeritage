package memgraph

import (
	"context"
	"errors"
	"testing"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/stretchr/testify/assert"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/internal/memgraph/mock"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/pkg/api"
)

func TestGetPersonByGoogleId(t *testing.T) {
	tests := []struct {
		name          string
		mockRunResult any
		mockRunError  error
		mockSingle    *neo4j.Record
		mockSingleErr error
		expectedErr   bool
		expectedRes   any
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
			expectedRes:   "value",
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
			mockResult := tt.mockRunResult.(*mock.Result)

			mockTx.On("Run", ctx, GetPersonByGoogleIdCypherQuery).Return(tt.mockRunResult, tt.mockRunError)
			if mockResult != nil {
				mockResult.On("Single", ctx).Return(tt.mockSingle, tt.mockSingleErr)
			}

			work := GetPersonByGoogleId(ctx, "test-google-id")
			result, err := work(mockTx)

			if tt.expectedErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedRes, result)
			}

			mockTx.AssertExpectations(t)
			if mockResult != nil {
				mockResult.AssertExpectations(t)
			}
		})
	}
}

func TestUpdatePerson(t *testing.T) {
	tests := []struct {
		name         string
		mockRunError error
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
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}

			mockTx.AssertExpectations(t)
		})
	}
}

package memgraph

import (
	"context"
	"errors"
	"testing"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/stretchr/testify/assert"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/internal/memgraph/mock"
)

func TestGetFamilyTreeById(t *testing.T) {
	testCases := []struct {
		name           string
		mockTxSetup    func() *mock.Transaction
		expectedResult map[string]any
		expectedError  error
	}{
		{
			name: "Successful case",
			mockTxSetup: func() *mock.Transaction {
				mockTx := new(mock.Transaction)
				mockResult := new(mock.Result)
				mockRecord := &neo4j.Record{
					Values: []any{"value"},
					Keys:   []string{"key"},
				}
				mockResult.On("Single", context.Background()).Return(mockRecord, nil)
				mockTx.On("Run", context.Background(), GetBloodRelativesCypherQuery, map[string]any{"id": 123}).Return(mockResult, nil)
				return mockTx
			},
			expectedResult: map[string]any{"key": "value"},
			expectedError:  nil,
		},
		{
			name: "Error during Run",
			mockTxSetup: func() *mock.Transaction {
				mockTx := new(mock.Transaction)
				mockTx.On("Run", context.Background(), GetBloodRelativesCypherQuery, map[string]any{"id": 123}).Return(nil, errors.New("run error"))
				return mockTx
			},
			expectedResult: nil,
			expectedError:  errors.New("run error"),
		},
		{
			name: "Error during Single",
			mockTxSetup: func() *mock.Transaction {
				mockTx := new(mock.Transaction)
				mockResult := new(mock.Result)
				mockResult.On("Single", context.Background()).Return(nil, errors.New("single error"))
				mockTx.On("Run", context.Background(), GetBloodRelativesCypherQuery, map[string]any{"id": 123}).Return(mockResult, nil)
				return mockTx
			},
			expectedResult: nil,
			expectedError:  errors.New("single error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			userId := 123

			mockTx := tc.mockTxSetup()
			work := GetFamilyTreeById(ctx, userId)
			result, err := work(mockTx)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, result)
			}
		})
	}
}

func TestGetFamilyTreeWithSpousesById(t *testing.T) {
	testCases := []struct {
		name           string
		mockTxSetup    func() *mock.Transaction
		expectedResult map[string]any
		expectedError  error
	}{
		{
			name: "Successful case",
			mockTxSetup: func() *mock.Transaction {
				mockTx := new(mock.Transaction)
				mockResult := new(mock.Result)
				mockRecord := &neo4j.Record{
					Values: []any{"value"},
					Keys:   []string{"key"},
				}
				mockResult.On("Single", context.Background()).Return(mockRecord, nil)
				mockTx.On("Run", context.Background(), GetFamilyTreeWithSpousesCypherQuery, map[string]any{"id": 123}).Return(mockResult, nil)
				return mockTx
			},
			expectedResult: map[string]any{"key": "value"},
			expectedError:  nil,
		},
		{
			name: "Error during Run",
			mockTxSetup: func() *mock.Transaction {
				mockTx := new(mock.Transaction)
				mockTx.On("Run", context.Background(), GetFamilyTreeWithSpousesCypherQuery, map[string]any{"id": 123}).Return(nil, errors.New("run error"))
				return mockTx
			},
			expectedResult: nil,
			expectedError:  errors.New("run error"),
		},
		{
			name: "Error during Single",
			mockTxSetup: func() *mock.Transaction {
				mockTx := new(mock.Transaction)
				mockResult := new(mock.Result)
				mockResult.On("Single", context.Background()).Return(nil, errors.New("single error"))
				mockTx.On("Run", context.Background(), GetFamilyTreeWithSpousesCypherQuery, map[string]any{"id": 123}).Return(mockResult, nil)
				return mockTx
			},
			expectedResult: nil,
			expectedError:  errors.New("single error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			userId := 123

			mockTx := tc.mockTxSetup()
			work := GetFamilyTreeWithSpousesById(ctx, userId)
			result, err := work(mockTx)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, result)
			}
		})
	}
}

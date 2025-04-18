package memgraph

import (
	"context"
	"errors"
	"testing"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/stretchr/testify/require"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/internal/memgraph/mock"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/pkg/api"
)

func TestGetRelationship(t *testing.T) {
	testCases := []struct {
		mockTxSetup    func() *mock.Transaction
		expectedResult map[string]any
		expectedError  error
		name           string
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
				mockTx.On("Run", context.Background(), GetRelationshipCypherQuery, map[string]any{"id1": 1, "id2": 2}).Return(mockResult, nil)
				return mockTx
			},
			expectedResult: map[string]any{"key": "value"},
			expectedError:  nil,
		},
		{
			name: "Error during Run",
			mockTxSetup: func() *mock.Transaction {
				mockTx := new(mock.Transaction)
				mockTx.On("Run", context.Background(), GetRelationshipCypherQuery, map[string]any{"id1": 1, "id2": 2}).Return(nil, errors.New("run error"))
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
				mockTx.On("Run", context.Background(), GetRelationshipCypherQuery, map[string]any{"id1": 1, "id2": 2}).Return(mockResult, nil)
				return mockTx
			},
			expectedResult: nil,
			expectedError:  errors.New("single error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			id1, id2 := 1, 2

			mockTx := tc.mockTxSetup()
			work := GetRelationship(ctx, id1, id2)
			result, err := work(mockTx)

			if tc.expectedError != nil {
				require.Error(t, err)
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedResult, result)
			}
		})
	}
}

func TestDeleteRelationship(t *testing.T) {
	testCases := []struct {
		mockTxSetup    func() *mock.Transaction
		expectedError  error
		name           string
		expectedResult bool
	}{
		{
			name: "Successful case",
			mockTxSetup: func() *mock.Transaction {
				mockTx := new(mock.Transaction)
				mockResult := new(mock.Result)
				mockResult.On("Peek", context.Background()).Return(false, nil)
				mockTx.On("Run", context.Background(), DeleteRelationshipCypherQuery, map[string]any{"id1": 1, "id2": 2}).Return(mockResult, nil)
				return mockTx
			},
			expectedResult: true,
			expectedError:  nil,
		},
		{
			name: "Error during Run",
			mockTxSetup: func() *mock.Transaction {
				mockTx := new(mock.Transaction)
				mockTx.On("Run", context.Background(), DeleteRelationshipCypherQuery, map[string]any{"id1": 1, "id2": 2}).Return(nil, errors.New("run error"))
				return mockTx
			},
			expectedResult: false,
			expectedError:  errors.New("run error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			id1, id2 := 1, 2

			mockTx := tc.mockTxSetup()
			work := DeleteRelationship(ctx, id1, id2)
			result, err := work(mockTx)

			if tc.expectedError != nil {
				require.Error(t, err)
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedResult, result)
			}
		})
	}
}

func TestUpdateRelationship(t *testing.T) {
	testCases := []struct {
		mockTxSetup    func() *mock.Transaction
		expectedResult map[string]any
		expectedError  error
		name           string
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
				mockTx.On("Run", context.Background(), UpdateRelationshipCypherQuery, map[string]any{
					"id1":          1,
					"id2":          2,
					"relationship": api.FamilyRelationship{},
				}).Return(mockResult, nil)
				return mockTx
			},
			expectedResult: map[string]any{"key": "value"},
			expectedError:  nil,
		},
		{
			name: "Error during Run",
			mockTxSetup: func() *mock.Transaction {
				mockTx := new(mock.Transaction)
				mockTx.On("Run", context.Background(), UpdateRelationshipCypherQuery, map[string]any{
					"id1":          1,
					"id2":          2,
					"relationship": api.FamilyRelationship{},
				}).Return(nil, errors.New("run error"))
				return mockTx
			},
			expectedResult: nil,
			expectedError:  errors.New("run error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			id1, id2 := 1, 2
			relationship := api.FamilyRelationship{}

			mockTx := tc.mockTxSetup()
			work := UpdateRelationship(ctx, id1, id2, relationship)
			result, err := work(mockTx)

			if tc.expectedError != nil {
				require.Error(t, err)
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedResult, result)
			}
		})
	}
}

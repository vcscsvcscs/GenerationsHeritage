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

func TestCreatePerson(t *testing.T) {
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
				mockTx.On("Run", context.Background(), CreatePersonCypherQuery, map[string]any{
					"Person": map[string]any{},
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
				mockTx.On("Run", context.Background(), CreatePersonCypherQuery, map[string]any{
					"Person": map[string]any{},
				}).Return(nil, errors.New("run error"))
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
				mockTx.On("Run", context.Background(), CreatePersonCypherQuery, map[string]any{
					"Person": map[string]any{},
				}).Return(mockResult, nil)
				return mockTx
			},
			expectedResult: nil,
			expectedError:  errors.New("single error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			person := &api.PersonProperties{}

			mockTx := tc.mockTxSetup()
			work := CreatePerson(ctx, person)
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

func TestHardDeletePerson(t *testing.T) {
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
				mockResult.On("Peek", context.Background()).Return(false, nil)
				mockTx.On("Run", context.Background(), HardDeletePersonCypherQuery, map[string]any{"id": 123}).Return(mockResult, nil)
				return mockTx
			},
			expectedResult: nil,
			expectedError:  nil,
		},
		{
			name: "Error during Run",
			mockTxSetup: func() *mock.Transaction {
				mockTx := new(mock.Transaction)
				mockTx.On("Run", context.Background(), HardDeletePersonCypherQuery, map[string]any{"id": 123}).Return(nil, errors.New("run error"))
				return mockTx
			},
			expectedResult: nil,
			expectedError:  errors.New("run error"),
		},
		{
			name: "Unexpected record returned",
			mockTxSetup: func() *mock.Transaction {
				mockTx := new(mock.Transaction)
				mockResult := new(mock.Result)
				mockResult.On("Peek", context.Background()).Return(true, nil)
				mockTx.On("Run", context.Background(), HardDeletePersonCypherQuery, map[string]any{"id": 123}).Return(mockResult, nil)
				return mockTx
			},
			expectedResult: nil,
			expectedError:  errors.New("record was returned when it wasn't supposed to happen"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			id := 123

			mockTx := tc.mockTxSetup()
			work := HardDeletePerson(ctx, id)
			result, err := work(mockTx)

			if tc.expectedError != nil {
				require.Error(t, err)
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.Nil(t, result)
			}
		})
	}
}

func TestUpdatePerson(t *testing.T) {
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
					Values: []any{map[string]any{"updatedKey": "updatedValue"}},
					Keys:   []string{"person"},
				}
				mockResult.On("Single", context.Background()).Return(mockRecord, nil)
				mockTx.On("Run", context.Background(), UpdatePersonCypherQuery, map[string]any{
					"id":    123,
					"props": map[string]any{},
				}).Return(mockResult, nil)
				return mockTx
			},
			expectedResult: map[string]any{"updatedKey": "updatedValue"},
			expectedError:  nil,
		},
		{
			name: "Error during Run",
			mockTxSetup: func() *mock.Transaction {
				mockTx := new(mock.Transaction)
				mockTx.On("Run", context.Background(), UpdatePersonCypherQuery, map[string]any{
					"id":    123,
					"props": map[string]any{},
				}).Return(nil, errors.New("run error"))
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
				mockTx.On("Run", context.Background(), UpdatePersonCypherQuery, map[string]any{
					"id":    123,
					"props": map[string]any{},
				}).Return(mockResult, nil)
				return mockTx
			},
			expectedResult: nil,
			expectedError:  errors.New("single error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			id := 123
			person := &api.PersonProperties{}

			mockTx := tc.mockTxSetup()
			work := UpdatePerson(ctx, id, person)
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

func TestSoftDeletePerson(t *testing.T) {
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
					Values: []any{"deletedValue"},
					Keys:   []string{"deletedKey"},
				}
				mockResult.On("Single", context.Background()).Return(mockRecord, nil)
				mockTx.On("Run", context.Background(), SoftDeletePersonCypherQuery, map[string]any{
					"id": 123,
				}).Return(mockResult, nil)
				return mockTx
			},
			expectedResult: map[string]any{"deletedKey": "deletedValue"},
			expectedError:  nil,
		},
		{
			name: "Error during Run",
			mockTxSetup: func() *mock.Transaction {
				mockTx := new(mock.Transaction)
				mockTx.On("Run", context.Background(), SoftDeletePersonCypherQuery, map[string]any{
					"id": 123,
				}).Return(nil, errors.New("run error"))
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
				mockTx.On("Run", context.Background(), SoftDeletePersonCypherQuery, map[string]any{
					"id": 123,
				}).Return(mockResult, nil)
				return mockTx
			},
			expectedResult: nil,
			expectedError:  errors.New("single error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			id := 123

			mockTx := tc.mockTxSetup()
			work := SoftDeletePerson(ctx, id)
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

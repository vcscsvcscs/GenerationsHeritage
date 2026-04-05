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

func TestCreateRecipeForPerson(t *testing.T) {
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
					Values: []any{"recipeValue", "relValue"},
					Keys:   []string{"recipe", "relationship"},
				}
				mockResult.On("Single", context.Background()).Return(mockRecord, nil)
				mockTx.On("Run", context.Background(), CreateRecipeWithRelationshipCypherQuery, map[string]any{
					"personId":     123,
					"Recipe":       map[string]any{},
					"Relationship": map[string]any{},
				}).Return(mockResult, nil)
				return mockTx
			},
			expectedResult: map[string]any{"recipe": "recipeValue", "relationship": "relValue"},
			expectedError:  nil,
		},
		{
			name: "Error during Run",
			mockTxSetup: func() *mock.Transaction {
				mockTx := new(mock.Transaction)
				mockTx.On("Run", context.Background(), CreateRecipeWithRelationshipCypherQuery, map[string]any{
					"personId":     123,
					"Recipe":       map[string]any{},
					"Relationship": map[string]any{},
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
				mockTx.On("Run", context.Background(), CreateRecipeWithRelationshipCypherQuery, map[string]any{
					"personId":     123,
					"Recipe":       map[string]any{},
					"Relationship": map[string]any{},
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
			recipe := &api.RecipeProperties{}

			mockTx := tc.mockTxSetup()
			work := CreateRecipeForPerson(ctx, 123, recipe, nil)
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

func TestGetRecipesByPersonId(t *testing.T) {
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
					Values: []any{[]any{"recipe1"}, []any{"rel1"}},
					Keys:   []string{"recipes", "recipeRelations"},
				}
				mockResult.On("Single", context.Background()).Return(mockRecord, nil)
				mockTx.On("Run", context.Background(), GetRecipesByPersonIdCypherQuery, map[string]any{
					"id": 123,
				}).Return(mockResult, nil)
				return mockTx
			},
			expectedResult: map[string]any{"recipes": []any{"recipe1"}, "recipeRelations": []any{"rel1"}},
			expectedError:  nil,
		},
		{
			name: "Error during Run",
			mockTxSetup: func() *mock.Transaction {
				mockTx := new(mock.Transaction)
				mockTx.On("Run", context.Background(), GetRecipesByPersonIdCypherQuery, map[string]any{
					"id": 123,
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

			mockTx := tc.mockTxSetup()
			work := GetRecipesByPersonId(ctx, 123)
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

func TestUpdateRecipe(t *testing.T) {
	testCases := []struct {
		mockTxSetup    func() *mock.Transaction
		expectedResult any
		expectedError  error
		name           string
	}{
		{
			name: "Successful case",
			mockTxSetup: func() *mock.Transaction {
				mockTx := new(mock.Transaction)
				mockResult := new(mock.Result)
				mockRecord := &neo4j.Record{
					Values: []any{map[string]any{"name": "Updated Recipe"}},
					Keys:   []string{"recipe"},
				}
				mockResult.On("Single", context.Background()).Return(mockRecord, nil)
				mockTx.On("Run", context.Background(), UpdateRecipeCypherQuery, map[string]any{
					"id":    123,
					"props": map[string]any{},
				}).Return(mockResult, nil)
				return mockTx
			},
			expectedResult: map[string]any{"name": "Updated Recipe"},
			expectedError:  nil,
		},
		{
			name: "Error during Run",
			mockTxSetup: func() *mock.Transaction {
				mockTx := new(mock.Transaction)
				mockTx.On("Run", context.Background(), UpdateRecipeCypherQuery, map[string]any{
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
				mockTx.On("Run", context.Background(), UpdateRecipeCypherQuery, map[string]any{
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
			recipe := &api.RecipeProperties{}

			mockTx := tc.mockTxSetup()
			work := UpdateRecipe(ctx, 123, recipe)
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

func TestSoftDeleteRecipe(t *testing.T) {
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
				mockTx.On("Run", context.Background(), SoftDeleteRecipeCypherQuery, map[string]any{
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
				mockTx.On("Run", context.Background(), SoftDeleteRecipeCypherQuery, map[string]any{
					"id": 123,
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

			mockTx := tc.mockTxSetup()
			work := SoftDeleteRecipe(ctx, 123)
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

func TestHardDeleteRecipe(t *testing.T) {
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
				mockTx.On("Run", context.Background(), HardDeleteRecipeCypherQuery, map[string]any{"id": 123}).Return(mockResult, nil)
				return mockTx
			},
			expectedResult: nil,
			expectedError:  nil,
		},
		{
			name: "Error during Run",
			mockTxSetup: func() *mock.Transaction {
				mockTx := new(mock.Transaction)
				mockTx.On("Run", context.Background(), HardDeleteRecipeCypherQuery, map[string]any{"id": 123}).Return(nil, errors.New("run error"))
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
				mockTx.On("Run", context.Background(), HardDeleteRecipeCypherQuery, map[string]any{"id": 123}).Return(mockResult, nil)
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
			work := HardDeleteRecipe(ctx, id)
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

func TestCreateRecipeRelationship(t *testing.T) {
	testCases := []struct {
		mockTxSetup    func() *mock.Transaction
		expectedResult any
		expectedError  error
		name           string
	}{
		{
			name: "Successful case",
			mockTxSetup: func() *mock.Transaction {
				mockTx := new(mock.Transaction)
				mockResult := new(mock.Result)
				mockRecord := &neo4j.Record{
					Values: []any{"relValue"},
					Keys:   []string{"relationship"},
				}
				mockResult.On("Single", context.Background()).Return(mockRecord, nil)
				mockTx.On("Run", context.Background(), CreateRecipeRelationshipCypherQuery, map[string]any{
					"personId":     1,
					"recipeId":     2,
					"Relationship": map[string]any{},
				}).Return(mockResult, nil)
				return mockTx
			},
			expectedResult: "relValue",
			expectedError:  nil,
		},
		{
			name: "Error during Run",
			mockTxSetup: func() *mock.Transaction {
				mockTx := new(mock.Transaction)
				mockTx.On("Run", context.Background(), CreateRecipeRelationshipCypherQuery, map[string]any{
					"personId":     1,
					"recipeId":     2,
					"Relationship": map[string]any{},
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

			mockTx := tc.mockTxSetup()
			work := CreateRecipeRelationship(ctx, 1, 2, nil)
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

func TestDeleteRecipeRelationship(t *testing.T) {
	testCases := []struct {
		mockTxSetup    func() *mock.Transaction
		expectedResult any
		expectedError  error
		name           string
	}{
		{
			name: "Successful case",
			mockTxSetup: func() *mock.Transaction {
				mockTx := new(mock.Transaction)
				mockResult := new(mock.Result)
				mockResult.On("Peek", context.Background()).Return(false, nil)
				mockTx.On("Run", context.Background(), DeleteRecipeRelationshipCypherQuery, map[string]any{
					"personId": 1,
					"recipeId": 2,
				}).Return(mockResult, nil)
				return mockTx
			},
			expectedResult: nil,
			expectedError:  nil,
		},
		{
			name: "Error during Run",
			mockTxSetup: func() *mock.Transaction {
				mockTx := new(mock.Transaction)
				mockTx.On("Run", context.Background(), DeleteRecipeRelationshipCypherQuery, map[string]any{
					"personId": 1,
					"recipeId": 2,
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

			mockTx := tc.mockTxSetup()
			work := DeleteRecipeRelationship(ctx, 1, 2)
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

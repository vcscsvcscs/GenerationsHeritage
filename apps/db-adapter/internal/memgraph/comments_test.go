package memgraph

import (
	"context"
	"errors"
	"testing"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	mmock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/internal/memgraph/mock"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/pkg/api"
)

func TestUpsertCommentOnProfile(t *testing.T) {
	ctx := context.Background()
	comment := &api.Message{Message: api.StringPtr("Hello!")}
	commentMap := StructToMap(comment)

	testCases := []struct {
		expectedResult map[string]any
		mockTxSetup    func() *mock.Transaction
		expectedError  error
		name           string
	}{
		{
			name: "Successful case",
			mockTxSetup: func() *mock.Transaction {
				mockTx := new(mock.Transaction)
				mockResult := new(mock.Result)
				mockRecord := &neo4j.Record{Values: []any{"val"}, Keys: []string{"out"}}
				mockResult.On("Single", ctx).Return(mockRecord, nil)
				mockTx.On("Run", ctx, CommentCypherQuery, map[string]any{"id1": 1, "id2": 2, "comment": commentMap}).Return(mockResult, nil)
				return mockTx
			},
			expectedResult: map[string]any{"out": "val"},
			expectedError:  nil,
		},
		{
			name: "Run error",
			mockTxSetup: func() *mock.Transaction {
				mockTx := new(mock.Transaction)
				mockTx.On("Run", ctx, CommentCypherQuery, mmock.Anything).Return(nil, errors.New("run error"))
				return mockTx
			},
			expectedResult: nil,
			expectedError:  errors.New("run error"),
		},
		{
			name: "Single error",
			mockTxSetup: func() *mock.Transaction {
				mockTx := new(mock.Transaction)
				mockResult := new(mock.Result)
				mockResult.On("Single", ctx).Return(nil, errors.New("single error"))
				mockTx.On("Run", ctx, CommentCypherQuery, mmock.Anything).Return(mockResult, nil)
				return mockTx
			},
			expectedResult: nil,
			expectedError:  errors.New("single error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			work := UpsertCommentOnProfile(ctx, 1, 2, comment)
			result, err := work(tc.mockTxSetup())

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

func TestGetCommentsOnProfile(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		expectedResult map[string]any
		mockTxSetup    func() *mock.Transaction
		expectedError  error
		name           string
	}{
		{
			name: "Successful case",
			mockTxSetup: func() *mock.Transaction {
				mockTx := new(mock.Transaction)
				mockResult := new(mock.Result)
				mockRecord := &neo4j.Record{Values: []any{"some"}, Keys: []string{"comments"}}
				mockResult.On("Single", ctx).Return(mockRecord, nil)
				mockTx.On("Run", ctx, CommentsOnProfileCypherQuery, map[string]any{"id": 2}).Return(mockResult, nil)
				return mockTx
			},
			expectedResult: map[string]any{"comments": "some"},
			expectedError:  nil,
		},
		{
			name: "Run error",
			mockTxSetup: func() *mock.Transaction {
				mockTx := new(mock.Transaction)
				mockTx.On("Run", ctx, CommentsOnProfileCypherQuery, mmock.Anything).Return(nil, errors.New("run error"))
				return mockTx
			},
			expectedResult: nil,
			expectedError:  errors.New("run error"),
		},
		{
			name: "Single error",
			mockTxSetup: func() *mock.Transaction {
				mockTx := new(mock.Transaction)
				mockResult := new(mock.Result)
				mockResult.On("Single", ctx).Return(nil, errors.New("single error"))
				mockTx.On("Run", ctx, CommentsOnProfileCypherQuery, mmock.Anything).Return(mockResult, nil)
				return mockTx
			},
			expectedResult: nil,
			expectedError:  errors.New("single error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			work := GetCommentsOnProfile(ctx, 2)
			result, err := work(tc.mockTxSetup())

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

func TestDeleteComment(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		mockTxSetup   func() *mock.Transaction
		expectedError error
		name          string
	}{
		{
			name: "Successful deletion",
			mockTxSetup: func() *mock.Transaction {
				mockTx := new(mock.Transaction)
				mockResult := new(mock.Result)
				mockResult.On("Peek", ctx).Return(false)
				mockTx.On("Run", ctx, DeleteCommentCypherQuery, map[string]any{"id1": 1, "id2": 2}).Return(mockResult, nil)
				return mockTx
			},
			expectedError: nil,
		},
		{
			name: "Run error",
			mockTxSetup: func() *mock.Transaction {
				mockTx := new(mock.Transaction)
				mockTx.On("Run", ctx, DeleteCommentCypherQuery, mmock.Anything).Return(nil, errors.New("run error"))
				return mockTx
			},
			expectedError: errors.New("run error"),
		},
		{
			name: "Peek unexpected return",
			mockTxSetup: func() *mock.Transaction {
				mockTx := new(mock.Transaction)
				mockResult := new(mock.Result)
				mockResult.On("Peek", ctx).Return(true)
				mockTx.On("Run", ctx, DeleteCommentCypherQuery, mmock.Anything).Return(mockResult, nil)
				return mockTx
			},
			expectedError: errors.New("there was a returned value, when deleting admin but there should be none"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			work := DeleteComment(ctx, 1, 2)
			result, err := work(tc.mockTxSetup())

			if tc.expectedError != nil {
				require.Error(t, err)
				require.Nil(t, result)
				require.EqualError(t, err, tc.expectedError.Error())
			} else {
				require.NoError(t, err)
				require.Nil(t, result)
			}
		})
	}
}

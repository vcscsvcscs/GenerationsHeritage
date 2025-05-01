package mock

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/stretchr/testify/mock"
)

type Transaction struct {
	neo4j.ManagedTransaction
	mock.Mock
}

func (m *Transaction) Run(ctx context.Context, cypher string, params map[string]any) (neo4j.ResultWithContext, error) {
	args := m.Called(ctx, cypher, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(neo4j.ResultWithContext), args.Error(1)
}

func (m *Transaction) legacy() neo4j.Transaction {
	return m.Called().Get(0).(neo4j.Transaction)
}

type ExplicitTransaction struct {
	neo4j.ExplicitTransaction
	mock.Mock
}

func (m *ExplicitTransaction) Run(ctx context.Context, cypher string, params map[string]any) (neo4j.ResultWithContext, error) {
	args := m.Called(ctx, cypher, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(neo4j.ResultWithContext), args.Error(1)
}

func (m *ExplicitTransaction) Commit(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *ExplicitTransaction) Rollback(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *ExplicitTransaction) Close(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *ExplicitTransaction) legacy() neo4j.Transaction {
	return m.Called().Get(0).(neo4j.Transaction)
}

package mock

import (
	"context"
	"sync"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/stretchr/testify/mock"
)

type SessionWithContext struct {
	neo4j.SessionWithContext
	ReturnOnce *sync.Once
	mock.Mock
}

var _ neo4j.SessionWithContext = &SessionWithContext{}

func (m *SessionWithContext) LastBookmarks() neo4j.Bookmarks {
	args := m.Called()
	if args.Get(0) != nil {
		return args.Get(0).(neo4j.Bookmarks)
	}

	return nil
}

func (m *SessionWithContext) lastBookmark() string {
	args := m.Called()
	if args.Get(0) != nil {
		return args.String(0)
	}

	return ""
}

func (m *SessionWithContext) BeginTransaction(ctx context.Context, configurers ...func(*neo4j.TransactionConfig)) (neo4j.ExplicitTransaction, error) {
	args := m.Called(ctx, configurers)
	if args.Get(0) != nil {
		return args.Get(0).(neo4j.ExplicitTransaction), args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *SessionWithContext) ExecuteRead(ctx context.Context, work neo4j.ManagedTransactionWork, configurers ...func(*neo4j.TransactionConfig)) (any, error) {
	args := m.Called(ctx, work, configurers)
	if len(args) > 2 && args.Get(2) != nil {
		returnValue1, returnValue2 := args.Get(2), args.Error(3)

		m.ReturnOnce.Do(func() {
			if args.Get(0) != nil {
				returnValue1, returnValue2 = args.Get(0), args.Error(1)
			}

			returnValue1, returnValue2 = nil, args.Error(1)
		})

		return returnValue1, returnValue2
	} else {
		if args.Get(0) != nil {
			return args.Get(0), args.Error(1)
		}

		return nil, args.Error(1)
	}
}

func (m *SessionWithContext) ExecuteWrite(ctx context.Context, work neo4j.ManagedTransactionWork, configurers ...func(*neo4j.TransactionConfig)) (any, error) {
	args := m.Called(ctx, work, configurers)
	if args.Get(0) != nil {
		return args.Get(0), args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *SessionWithContext) Run(ctx context.Context, cypher string, params map[string]any, configurers ...func(*neo4j.TransactionConfig)) (neo4j.ResultWithContext, error) {
	args := m.Called(ctx, cypher, params, configurers)
	if args.Get(0) != nil {
		return args.Get(0).(neo4j.ResultWithContext), args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *SessionWithContext) Close(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *SessionWithContext) executeQueryRead(ctx context.Context, work neo4j.ManagedTransactionWork, configurers ...func(*neo4j.TransactionConfig)) (any, error) {
	args := m.Called(ctx, work, configurers)
	if args.Get(0) != nil {
		return args.Get(0), args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *SessionWithContext) executeQueryWrite(ctx context.Context, work neo4j.ManagedTransactionWork, configurers ...func(*neo4j.TransactionConfig)) (any, error) {
	args := m.Called(ctx, work, configurers)
	if args.Get(0) != nil {
		return args.Get(0), args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *SessionWithContext) legacy() neo4j.Session {
	args := m.Called()
	if args.Get(0) != nil {
		return args.Get(0).(neo4j.Session)
	}

	return nil
}

func (m *SessionWithContext) getServerInfo(ctx context.Context) (neo4j.ServerInfo, error) {
	args := m.Called(ctx)
	if args.Get(0) != nil {
		return args.Get(0).(neo4j.ServerInfo), args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *SessionWithContext) verifyAuthentication(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

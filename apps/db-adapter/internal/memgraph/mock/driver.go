package mock

import (
	"context"
	"net/url"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/stretchr/testify/mock"
)

type DriverWithContext struct {
	mock.Mock
}

var _ neo4j.DriverWithContext = &DriverWithContext{}

func (m *DriverWithContext) ExecuteQueryBookmarkManager() neo4j.BookmarkManager {
	args := m.Called()
	return args.Get(0).(neo4j.BookmarkManager)
}

func (m *DriverWithContext) Target() url.URL {
	args := m.Called()
	return args.Get(0).(url.URL)
}

func (m *DriverWithContext) NewSession(ctx context.Context, config neo4j.SessionConfig) neo4j.SessionWithContext {
	args := m.Called(ctx, config)
	return args.Get(0).(neo4j.SessionWithContext)
}

func (m *DriverWithContext) VerifyConnectivity(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *DriverWithContext) VerifyAuthentication(ctx context.Context, auth *neo4j.AuthToken) error {
	args := m.Called(ctx, auth)
	return args.Error(0)
}

func (m *DriverWithContext) Close(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *DriverWithContext) IsEncrypted() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *DriverWithContext) GetServerInfo(ctx context.Context) (neo4j.ServerInfo, error) {
	args := m.Called(ctx)
	if args.Get(0) != nil {
		return args.Get(0).(neo4j.ServerInfo), args.Error(1)
	}

	return nil, args.Error(1)
}

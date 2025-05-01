package mock

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/stretchr/testify/mock"
)

type Result struct {
	neo4j.ResultWithContext
	mock.Mock
}

// Ensure Result implements neo4j.ResultWithContext interface
var _ neo4j.ResultWithContext = (*Result)(nil)

func (m *Result) Keys() ([]string, error) {
	args := m.Called()
	if keys, ok := args.Get(0).([]string); ok {
		return keys, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *Result) NextRecord(ctx context.Context, record **neo4j.Record) bool {
	args := m.Called(ctx, record)
	return args.Bool(0)
}

func (m *Result) Next(ctx context.Context) bool {
	args := m.Called(ctx)
	return args.Bool(0)
}

func (m *Result) PeekRecord(ctx context.Context, record **neo4j.Record) bool {
	args := m.Called(ctx, record)
	return args.Bool(0)
}

func (m *Result) Peek(ctx context.Context) bool {
	args := m.Called(ctx)
	return args.Bool(0)
}

func (m *Result) Err() error {
	args := m.Called()
	return args.Error(0)
}

func (m *Result) Record() *neo4j.Record {
	args := m.Called()
	if record, ok := args.Get(0).(*neo4j.Record); ok {
		return record
	}

	return nil
}

func (m *Result) Collect(ctx context.Context) ([]*neo4j.Record, error) {
	args := m.Called(ctx)
	if records, ok := args.Get(0).([]*neo4j.Record); ok {
		return records, args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *Result) Records(ctx context.Context) func(yield func(*neo4j.Record, error) bool) {
	args := m.Called(ctx)
	if recordsFunc, ok := args.Get(0).(func(yield func(*neo4j.Record, error) bool)); ok {
		return recordsFunc
	}

	return nil
}

func (m *Result) Single(ctx context.Context) (*neo4j.Record, error) {
	args := m.Called(ctx)
	if record, ok := args.Get(0).(*neo4j.Record); ok {
		return record, args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *Result) Consume(ctx context.Context) (neo4j.ResultSummary, error) {
	args := m.Called(ctx)
	if summary, ok := args.Get(0).(neo4j.ResultSummary); ok {
		return summary, args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *Result) IsOpen() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *Result) buffer(ctx context.Context) {
	m.Called(ctx)
}

func (m *Result) legacy() neo4j.Result {
	args := m.Called()
	if result, ok := args.Get(0).(neo4j.Result); ok {
		return result
	}

	return nil
}

func (m *Result) errorHandler(err error) {
	m.Called(err)
}

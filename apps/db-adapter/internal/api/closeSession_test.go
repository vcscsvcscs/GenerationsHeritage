package api

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	memgraphMock "github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/internal/memgraph/mock"
	"go.uber.org/zap/zaptest"
)

func TestCloseSession(t *testing.T) {
	logger := zaptest.NewLogger(t)
	mockSession := new(memgraphMock.SessionWithContext)
	ctx := context.Background()
	timeOut := 2 * time.Second

	t.Run("successful session close", func(t *testing.T) {
		mockSession.On("Close", mock.Anything).Return(nil).Once()

		closeSession(ctx, logger, mockSession, timeOut)

		mockSession.AssertCalled(t, "Close", mock.Anything)
		mockSession.AssertExpectations(t)
	})

	t.Run("error during session close", func(t *testing.T) {
		mockSession.On("Close", mock.Anything).Return(errors.New("close error")).Once()

		closeSession(ctx, logger, mockSession, timeOut)

		mockSession.AssertCalled(t, "Close", mock.Anything)
		mockSession.AssertExpectations(t)
	})
}

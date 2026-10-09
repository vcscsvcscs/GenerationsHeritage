package memgraph

import (
	"context"
	"errors"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

var (
	// ErrNotFound is returned when a query expected to match a record matched none.
	ErrNotFound = errors.New("resource not found")
	// ErrRecipeNotDeleted is returned when hard deleting a recipe that was not soft deleted first.
	ErrRecipeNotDeleted = errors.New("recipe must be soft deleted before it can be hard deleted")
)

// singleOrNotFound returns the only record of the result, or ErrNotFound if the result is empty.
func singleOrNotFound(ctx context.Context, result neo4j.ResultWithContext) (*neo4j.Record, error) {
	if !result.Peek(ctx) {
		if err := result.Err(); err != nil {
			return nil, err
		}

		return nil, ErrNotFound
	}

	return result.Single(ctx)
}

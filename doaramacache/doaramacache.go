// Package doaramacache provides caching of Doarama activities.
package doaramacache

import (
	"context"
	"io"

	"github.com/twpayne/go-doarama"
)

// An ActivityCreator can create activities.
type ActivityCreator interface {
	// Close releases any resources.
	Close() error
	// CreateActivityWithInfo creates a new doarama.Activity with the specified
	// doarama.ActivityInfo.
	CreateActivityWithInfo(context.Context, string, io.Reader, *doarama.ActivityInfo) (*doarama.Activity, error)
}

// Package doaramacache provides caching of Doarama activities.
package doaramacache

import (
	"io"

	"github.com/twpayne/go-doarama"
	"golang.org/x/net/context"
)

// An ActivityCreator can create activities.
type ActivityCreator interface {
	// CreateActivityWithInfo creates a new doarama.Activity with the specified
	// doarama.ActivityInfo.
	CreateActivityWithInfo(context.Context, string, io.Reader, *doarama.ActivityInfo) (*doarama.Activity, error)
}

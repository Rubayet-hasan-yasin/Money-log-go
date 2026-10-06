package domain

import (
	"fmt"

	pkgerrors "github.com/pkg/errors"
)

// WrapErrorf wraps an error with formatted context
func WrapErrorf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return pkgerrors.Wrap(err, fmt.Sprintf(format, args...))
}

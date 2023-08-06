package log

import "github.com/cockroachdb/errors"

func WithStackTracef(err error, format string, args ...any) error {
	return errors.WithHintf(err, format, args...)
}

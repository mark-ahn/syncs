package syncs

import (
	"errors"
	"fmt"
)

type ServeFailError struct{ error }

func ServeFailErrorf(fmtStr string, args ...interface{}) ServeFailError {
	return ServeFailError{error: fmt.Errorf("ServeFailError -> "+fmtStr, args...)}
}
func (__ ServeFailError) Unwrap() error {
	return errors.Unwrap(__.error)
}

type ServeBrokenError struct{ error }

func ServeBrokenErrorf(fmtStr string, args ...interface{}) ServeBrokenError {
	return ServeBrokenError{error: fmt.Errorf("ServeBrokenError -> "+fmtStr, args...)}
}
func (__ ServeBrokenError) Unwrap() error {
	return errors.Unwrap(__.error)
}

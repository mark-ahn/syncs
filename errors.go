package syncs

import (
	"errors"
	"fmt"
)

type ServeThreadError struct{ error }

func ServeThreadErrorf(fmtStr string, args ...interface{}) ServeThreadError {
	return ServeThreadError{error: fmt.Errorf("ServeThreadError -> "+fmtStr, args...)}
}
func (__ ServeThreadError) Unwrap() error {
	return errors.Unwrap(__.error)
}

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

type NoKeyInArgMapError struct{ error }

func NoKeyInArgMapErrorFrom(err error) NoKeyInArgMapError {
	return NoKeyInArgMapError{error: err}
}

func NoKeyInArgMapErrorOf(key string) NoKeyInArgMapError {
	return NoKeyInArgMapError{error: fmt.Errorf("NoKeyInArgMapError: %s", key)}
}
func (__ NoKeyInArgMapError) Unwrap() error {
	return errors.Unwrap(__.error)
}

type TrySpawnThreadOnContextDoneError struct{ error }

func TrySpawnThreadOnContextDoneErrorf(fmtStr string, args ...interface{}) TrySpawnThreadOnContextDoneError {
	return TrySpawnThreadOnContextDoneError{error: fmt.Errorf("TrySpawnThreadOnContextDoneError: "+fmtStr, args...)}
}
func (__ TrySpawnThreadOnContextDoneError) Unwrap() error {
	return errors.Unwrap(__.error)
}

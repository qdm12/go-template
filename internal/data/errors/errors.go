// Package errors contains database errors common to all implementations.
package errors

import "errors"

// Shared errors package for all implementation of the database

var (
	ErrReadFile  = errors.New("cannot read file")
	ErrWriteFile = errors.New("cannot write data to file")
	ErrEncoding  = errors.New("failed encoding data to write")
	ErrDecoding  = errors.New("failed decoding data read")

	ErrUserNotFound = errors.New("user not found")
)

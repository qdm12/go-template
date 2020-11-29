package errors

import "errors"

// Shared errors package for all implementation of the database

var (
	ErrCreation  = errors.New("cannot create database")
	ErrClose     = errors.New("cannot close database")
	ErrReadFile  = errors.New("cannot read file")
	ErrWriteFile = errors.New("cannot write data to file")

	ErrCreateUser   = errors.New("cannot create user")
	ErrGetUser      = errors.New("cannot get user")
	ErrUserNotFound = errors.New("user not found")
)

package util

import "github.com/pkg/errors"

type UtilError interface {
	error
	IsFault() bool
}

type NaturalError struct {
	err error
	UtilError
}

func (e NaturalError) IsFault() bool {
	return false
}

func (e NaturalError) Error() string {
	return e.err.Error()
}

func NewNaturalError(err error) NaturalError {
	return NaturalError{
		err: err,
	}
}

func NewNaturalErrorS(text string) NaturalError {
	return NaturalError{
		err: errors.New(text),
	}
}

type FaultError struct {
	err error
	UtilError
}

func (e FaultError) IsFault() bool {
	return true
}

func (e FaultError) Error() string {
	return e.err.Error()
}

func NewFaultError(err error) FaultError {
	return FaultError{
		err: err,
	}
}

func NewFaultErrorS(text string) FaultError {
	return FaultError{
		err: errors.New(text),
	}
}

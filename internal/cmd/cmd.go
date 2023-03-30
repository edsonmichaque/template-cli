package cmd

import (
	"errors"
	"io"
	"os"
)

type Error struct {
	Code int
	Err  error
}

func (e Error) Error() string {
	return e.Err.Error()
}

func newErrorFromString(code int, err string) Error {
	return newError(code, errors.New(err))
}

func newError(code int, err error) Error {
	return Error{
		Code: code,
		Err:  err,
	}
}

func NewOptions() (*Options, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	return &Options{
		Stdin:   os.Stdin,
		Stderr:  os.Stderr,
		Stdout:  os.Stdout,
		WorkDir: wd,
	}, nil
}

type Options struct {
	Stdout  io.Writer
	Stdin   io.Reader
	Stderr  io.Writer
	WorkDir string
}

func (c Options) Validate() error {
	return nil
}

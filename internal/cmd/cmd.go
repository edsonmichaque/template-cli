package cmd

import (
	"errors"
	"io"
	"os"

	"github.com/dnsimple/dnsimple-go/dnsimple"
)

type Error struct {
	Code int
	Err  error
}

func (e Error) Error() string {
	return e.Err.Error()
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
	Stdout        io.Writer
	Stdin         io.Reader
	Stderr        io.Writer
	WorkDir       string
	ClientBuilder func(string, string) *dnsimple.Client
}

func (c Options) Validate() error {
	if c.ClientBuilder == nil {
		return errors.New("invalid client builder")
	}

	return nil
}

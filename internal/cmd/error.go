package cmd

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

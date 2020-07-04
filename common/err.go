package common

type ApiErr struct {
	Code int
	Err  error
}

func (e ApiErr) Error() string {
	return e.Err.Error()
}

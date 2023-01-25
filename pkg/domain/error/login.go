package error

type ErrLoginUnauthorized struct{ Stack error }

func (e ErrLoginUnauthorized) Error() string {
	return e.Stack.Error()
}

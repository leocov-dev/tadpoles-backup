package async

type CanceledError struct{}

func NewCanceledError() *CanceledError {
	return &CanceledError{}
}

func (e *CanceledError) Error() string {
	return "Canceled..."
}

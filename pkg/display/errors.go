package display

type DisplayError struct {
	e string
}

func (err *DisplayError) Error() string {
	return err.e
}

func NewDisplayError(err error) DisplayError {
	return DisplayError{err.Error()}
}

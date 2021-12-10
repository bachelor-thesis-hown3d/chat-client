package errors

type TokenFileNotFound string

func (t TokenFileNotFound) Error() string {
	return string(t)
}

package bases

import "fmt"

type Entity interface {
	DataBase() string
}

type bError struct {
	arg string
}

func (e *bError) Error() string {
	return fmt.Sprintf(e.arg)
}

func Error(arg string) *bError {
	return &bError{arg: arg}
}

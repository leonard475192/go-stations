package model

import "fmt"

type ErrNotFound string

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("not Found: %v", e)
}

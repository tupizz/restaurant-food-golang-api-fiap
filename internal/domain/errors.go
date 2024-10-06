package domain

import (
	"fmt"
)

type NotFoundError struct {
	Entity string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("Entity %s Not Found", e.Entity)
}

func (e *NotFoundError) Is(target error) bool {
	_, ok := target.(*NotFoundError)
	return ok
}

func ErrNotFound(entity string) error {
	return &NotFoundError{Entity: entity}
}

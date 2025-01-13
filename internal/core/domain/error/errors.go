package domain

import (
	"fmt"
)

type NotFoundError struct {
	Entity string
}

type EntityNotProcessableError struct {
	Entity string
	Reason string
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

func (e *EntityNotProcessableError) Error() string {
	return fmt.Sprintf("Entity %s not processable: %s", e.Entity, e.Reason)
}

func (e *EntityNotProcessableError) Is(target error) bool {
	_, ok := target.(*EntityNotProcessableError)
	return ok
}

func NewEntityNotProcessableError(entity, reason string) error {
	return &EntityNotProcessableError{Entity: entity, Reason: reason}
}

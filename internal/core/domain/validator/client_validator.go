package validator

import (
	"errors"

	"github.com/paemuri/brdoc"
)

func IsValidCPF(cpf string) error {
	if brdoc.IsCPF(cpf) {
		return nil
	}

	return errors.New("CPF inválido")
}

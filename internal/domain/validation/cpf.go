package validation

import (
	"regexp"
	"strconv"
)

func IsValidCPF(cpf string) bool {
	// Remover caracteres não numéricos
	re := regexp.MustCompile("[^0-9]")
	cpf = re.ReplaceAllString(cpf, "")

	if len(cpf) != 11 {
		return false
	}

	// Verificar se todos os dígitos são iguais
	invalidCpfs := []string{
		"00000000000", "11111111111", "22222222222",
		"33333333333", "44444444444", "55555555555",
		"66666666666", "77777777777", "88888888888",
		"99999999999",
	}
	for _, invalid := range invalidCpfs {
		if cpf == invalid {
			return false
		}
	}

	// Cálculo dos dígitos verificadores
	sum := 0
	for i := 0; i < 9; i++ {
		digit, _ := strconv.Atoi(string(cpf[i]))
		sum += digit * (10 - i)
	}

	firstDigit := ((sum * 10) % 11) % 10
	sum = 0
	for i := 0; i < 10; i++ {
		digit, _ := strconv.Atoi(string(cpf[i]))
		sum += digit * (11 - i)
	}

	secondDigit := ((sum * 10) % 11) % 10

	return string(cpf[9]) == strconv.Itoa(firstDigit) && string(cpf[10]) == strconv.Itoa(secondDigit)
}

package gateways

import (
	"encoding/base64"

	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/domain/entities"
)

type mercadoPagoQRCodeMock struct{}

func NewQRCodePaymentGateway() PaymentGateway {
	return &mercadoPagoQRCodeMock{}
}

func (g *mercadoPagoQRCodeMock) Authorize(payment *entities.Payment) error {
	png, err := qrcode.Encode("https://https://www.fiap.com.br", qrcode.Medium, 256)
	if err != nil {
		return err
	}

	payment.QRData = base64.StdEncoding.EncodeToString(png)
	payment.ExternalReference = uuid.New().String()

	return nil
}

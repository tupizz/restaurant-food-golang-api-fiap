package entities

import (
	"encoding/base64"
	"time"

	"github.com/google/uuid"
	qrcode "github.com/skip2/go-qrcode"
)

type PaymentStatus string

const (
	PaymentStatusPending  PaymentStatus = "pending"
	PaymentStatusApproved PaymentStatus = "approved"
	PaymentStatusFailed   PaymentStatus = "failed"
)

type PaymentMethod string

const (
	PaymentMethodPix        PaymentMethod = "pix"
	PaymentMethodCreditCard PaymentMethod = "credit_card"
	PaymentMethodBillet     PaymentMethod = "billet"
	PaymentMethodQRCode     PaymentMethod = "qr_code"
)

type Payment struct {
	ID                int
	OrderID           int
	Status            PaymentStatus
	Method            PaymentMethod
	Amount            float64
	ExternalReference string
	QRData            string
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         *time.Time
}

func (p *Payment) Authorize() error {
	if p.Method == PaymentMethodQRCode {
		png, err := qrcode.Encode("https://https://www.fiap.com.br", qrcode.Medium, 256)
		if err != nil {
			return err
		}

		pngAsBase64 := base64.StdEncoding.EncodeToString(png)
		p.QRData = pngAsBase64
		p.ExternalReference = uuid.New().String()

		return nil
	}

	return nil
}

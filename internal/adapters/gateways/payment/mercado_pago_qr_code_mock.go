package gateways

import (
	"context"
	"encoding/base64"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/skip2/go-qrcode"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/domain/entities"
)

const QrCodeTTL = 10 * time.Minute

type mercadoPagoQRCodeMock struct {
	redisClient *redis.Client
}

func NewQRCodePaymentGateway(redisClient *redis.Client) PaymentGateway {
	return &mercadoPagoQRCodeMock{
		redisClient: redisClient,
	}
}

func (g *mercadoPagoQRCodeMock) Authorize(payment *entities.Payment) error {
	payment.ExternalReference = uuid.New().String()

	ctx := context.Background()
	redisKey := "qrcode:" + payment.ExternalReference

	cachedQR, err := g.redisClient.Get(ctx, redisKey).Result()
	if err == nil {
		payment.QRData = cachedQR
		return nil
	}

	png, err := qrcode.Encode("https://www.fiap.com.br", qrcode.Medium, 256)
	if err != nil {
		return err
	}

	qrData := base64.StdEncoding.EncodeToString(png)
	payment.QRData = qrData

	_ = g.redisClient.Set(ctx, redisKey, qrData, QrCodeTTL).Err()

	return nil
}

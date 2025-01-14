package handler

import (
	"net/http"

	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/usecase"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/usecase/dto"

	"github.com/gin-gonic/gin"
)

type HealthcheckHandler interface {
	Ping(c *gin.Context)
}

type healthcheckHandler struct {
	healthCheckPingUseCase usecase.HealthCheckPingUseCase
}

func NewHealthcheckHandler(healthCheckPingUseCase usecase.HealthCheckPingUseCase) HealthcheckHandler {
	return &healthcheckHandler{healthCheckPingUseCase: healthCheckPingUseCase}
}

// Ping godoc
// @Summary     Verifica se a aplicação está respondendo
// @Description Verifica se a aplicação está respondendo para o liveness e readiness probe do K8s
// @Tags        healthcheck
// @Accept      json
// @Produce     json
// @Success     200     {object}  []dto.HealthCheckOutput
// @Failure     500     {object}  handler.ErrorResponse
// @Router      /healthcheck [get]
func (h *healthcheckHandler) Ping(c *gin.Context) {
	err := h.healthCheckPingUseCase.Run(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.HealthCheckOutput{Status: "ok"})
}

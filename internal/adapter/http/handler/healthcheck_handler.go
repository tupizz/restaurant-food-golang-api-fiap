package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/application/dto"
)

type HealthcheckHandler interface {
	Ping(c *gin.Context)
}

type healthcheckHandler struct {
	db *pgxpool.Pool
}

func NewHealthcheckHandler(db *pgxpool.Pool) HealthcheckHandler {
	return &healthcheckHandler{db: db}
}

// Ping godoc
// @Summary     Verifica se a aplicação está respondendo
// @Description Verifica se a aplicação está respondendo para o liveness e readiness probe do K8s
// @Tags        healthcheck
// @Accept      json
// @Produce     json
// @Success     200     {object}  []dto.HealthCheckOutput
// @Failure     500     ErrorResponse
// @Router      /healthcheck [get]
func (h *healthcheckHandler) Ping(c *gin.Context) {
	err := h.db.Ping(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.HealthCheckOutput{Status: "ok"})
}

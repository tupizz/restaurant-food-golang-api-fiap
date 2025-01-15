package handler

import (
	"errors"
	"net/http"

	domainError "github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/domain/error"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/usecase"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/usecase/dto"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/usecase/mappers"

	"github.com/gin-gonic/gin"
)

type ClientHandler interface {
	Create(c *gin.Context)
	GetByCPF(c *gin.Context)
}

type clientHandler struct {
	createClientUseCase   usecase.CreateClientUseCase
	getClientByCPFUseCase usecase.GetClientByCPFUseCase
}

func NewClientHandler(createClientUseCase usecase.CreateClientUseCase, getClientByCPFUseCase usecase.GetClientByCPFUseCase) ClientHandler {
	return &clientHandler{createClientUseCase: createClientUseCase, getClientByCPFUseCase: getClientByCPFUseCase}
}

// Create CreateClient godoc
// @Summary      Cria um novo cliente
// @Description  Cria um novo cliente com os dados fornecidos
// @Tags         clients
// @Accept       json
// @Produce      json
// @Param        client  body      dto.ClientInput  true  "Dados do Cliente"
// @Success      201     {object}  dto.ClientOutput
// @Failure      400     {object}  handler.ErrorResponse
// @Failure      500     {object}  handler.ErrorResponse
// @Router       /clients [post]
func (h *clientHandler) Create(c *gin.Context) {
	var input dto.ClientInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client, err := h.createClientUseCase.Run(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, mappers.ToClientDTO(*client))
}

// GetByCPF SearchClient godoc
// @Summary      Obtém um cliente pelo CPF
// @Description  Retorna os dados de um cliente específico pelo CPF
// @Tags         clients
// @Accept       json
// @Produce      json
// @Param        cpf   path      string  true  "CPF do Cliente"
// @Success      200   {object}  dto.ClientOutput
// @Success      200   {object}  dto.ClientOutput
// @Failure      404   {object}  handler.ErrorResponse
// @Failure      500   {object}  handler.ErrorResponse
// @Router       /clients/{cpf} [get]
func (h *clientHandler) GetByCPF(c *gin.Context) {
	cpf := c.Param("cpf")
	client, err := h.getClientByCPFUseCase.Run(c.Request.Context(), cpf)
	if err != nil {
		if errors.Is(err, domainError.ErrNotFound("client")) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cliente não encontrado"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, mappers.ToClientDTO(*client))
}

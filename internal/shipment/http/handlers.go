package http

import (
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"

	"github.com/shopspring/decimal"
	"net/http"
	"test-task/internal/shipment/repo"
	"test-task/internal/shipment/service"
	"time"
)

type IShipmentHandler interface {
	CreateShipment(c *fiber.Ctx) error
	GetApplication(c *fiber.Ctx) error
}

type Handler struct {
	svc service.IShipmentService
}

func NewHandler(svc service.IShipmentService) IShipmentHandler {
	return &Handler{svc: svc}
}

type createRequest struct {
	Route    string          `json:"route"`
	Price    decimal.Decimal `json:"price"`
	Customer customerRequest `json:"customer"`
}

type customerRequest struct {
	Idn string `json:"idn"`
}

type createResponse struct {
	Id         string           `json:"id"`
	Status     repo.StatusTypes `json:"status"`
	CustomerId string           `json:"customerId"`
}

func (h Handler) CreateShipment(c *fiber.Ctx) error {
	req := new(createRequest)
	err := json.Unmarshal(c.Body(), req)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(err.Error())
	}

	model, err := mapAndValidateCreateRequest(*req)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(err.Error())
	}

	response, err := h.svc.Create(c.Context(), *model, req.Customer.Idn)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err.Error())
	}

	return c.Status(http.StatusCreated).JSON(mapCreateShipmentResponse(*response))
}

func mapAndValidateCreateRequest(request createRequest) (*repo.Shipment, error) {
	if len(request.Customer.Idn) != 12 {
		return nil, errors.New("invalid customer idn")
	}

	if len(request.Route) == 0 || request.Price == decimal.Zero {
		return nil, errors.New("invalid route or price")
	}

	return &repo.Shipment{
		Route: request.Route,
		Price: request.Price,
	}, nil
}

func mapCreateShipmentResponse(response service.CreateShipmentResponse) createResponse {
	return createResponse{
		Id:         response.ShipmentId,
		Status:     response.Status,
		CustomerId: response.CustId,
	}
}

type applicationResponse struct {
	Id         string           `json:"id"`
	Route      string           `json:"route"`
	Price      decimal.Decimal  `json:"price"`
	Status     repo.StatusTypes `json:"status"`
	ConsumerId string           `json:"consumerId"`
	CreatedAt  time.Time        `json:"createdAt"`
}

func (h Handler) GetApplication(c *fiber.Ctx) error {
	idn := c.Params("id")
	if len(idn) == 0 {
		return c.Status(http.StatusBadRequest).JSON("no id provided")
	}

	res, err := h.svc.GetShipment(c.Context(), idn)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err.Error())
	}

	return c.Status(http.StatusOK).JSON(mapApplicationResponse(*res))
}

func mapApplicationResponse(data repo.Shipment) applicationResponse {
	return applicationResponse{
		Id:         data.ID,
		Route:      data.Route,
		Price:      data.Price,
		Status:     data.Status,
		ConsumerId: data.CustomerID,
		CreatedAt:  data.CreatedAt,
	}
}

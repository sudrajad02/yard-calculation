package handlers

import (
	"fmt"
	"net/http"
	"yard-calculation/schemas"
	"yard-calculation/services"
	"yard-calculation/utils"

	"github.com/gofiber/fiber/v2"
)

type ContainerHandler struct {
	Service *services.ContainerService
}

func NewContainerHandler(service *services.ContainerService) *ContainerHandler {
	return &ContainerHandler{Service: service}
}

func (h *ContainerHandler) GetSuggestion(c *fiber.Ctx) error {
	req := new(schemas.SuggestContainerRequest)
	if err := c.BodyParser(req); err != nil {
		utils.ApiResponse(c, http.StatusBadRequest, "Cannot parse JSON", nil, "Cannot parse JSON")
		return nil
	}

	// Validasi input
	if req.Yard == "" || req.ContainerNumber == "" || (req.ContainerSize != 20 && req.ContainerSize != 40) {
		utils.ApiResponse(c, http.StatusBadRequest, "Invalid input", nil, "Invalid input: yard, container_number, and valid size (20 or 40) are required")
		return nil
	}

	suggestedContainer, err := h.Service.GetSuggestedPosition(req.Yard, req.ContainerNumber, req.ContainerSize, req.ContainerHeight, req.ContainerType)
	if err != nil {
		utils.ApiResponse(c, http.StatusInternalServerError, "Error Get Suggest", nil, err.Error())
		return nil
	}

	utils.ApiResponse(c, http.StatusOK, "Suggest Container", suggestedContainer, nil)
	return nil
}

func (h *ContainerHandler) PlaceContainer(c *fiber.Ctx) error {
	req := new(schemas.PlaceContainerRequest)
	if err := c.BodyParser(req); err != nil {
		utils.ApiResponse(c, http.StatusBadRequest, "Cannot parse JSON", nil, "Cannot parse JSON")
		return nil
	}

	// Validasi input
	if req.Yard == "" || req.ContainerNumber == "" || req.Block == "" || req.Slot <= 0 || req.Row <= 0 || req.Tier <= 0 || (req.ContainerSize != 20 && req.ContainerSize != 40) {
		utils.ApiResponse(c, http.StatusBadRequest, "Invalid input", nil, "Invalid input: yard, container_number, and valid size (20 or 40) are required")
		return nil
	}

	err := h.Service.PlaceContainerDetailed(req.Yard, req.ContainerNumber, req.Block, req.Slot, req.Row, req.Tier, req.ContainerSize, req.ContainerHeight, req.ContainerType)
	if err != nil {
		utils.ApiResponse(c, http.StatusInternalServerError, "Error Place Container", nil, err.Error())
		return nil
	}

	utils.ApiResponse(c, http.StatusOK, "Place Container Success", nil, nil)
	return nil
}

func (h *ContainerHandler) PickupContainer(c *fiber.Ctx) error {
	req := new(schemas.PickupContainerRequest)
	if err := c.BodyParser(req); err != nil {
		utils.ApiResponse(c, http.StatusBadRequest, "Cannot parse JSON", nil, "Cannot parse JSON")
		return nil
	}

	// Validasi input
	if req.Yard == "" || req.ContainerNumber == "" {
		utils.ApiResponse(c, http.StatusBadRequest, "Invalid input", nil, "Invalid input: yard and container_number are required")
		return nil
	}

	err := h.Service.PickupContainer(req.Yard, req.ContainerNumber)
	if err != nil {
		if err.Error() == fmt.Sprintf("container with number %s not found or not placed", req.ContainerNumber) {
			utils.ApiResponse(c, http.StatusNotFound, "Error Pickup Container", nil, err.Error())
			return nil
		}
		utils.ApiResponse(c, http.StatusInternalServerError, "Error Pickup Container", nil, err.Error())
		return nil
	}

	utils.ApiResponse(c, http.StatusOK, "Pickup Container Success", nil, nil)
	return nil
}

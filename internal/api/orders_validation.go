package api

import (
	"fmt"
	"strings"

	"github.com/Ashish-Barmaiya/torus-demo-backends/internal/data"
)

func validateCreateOrderRequest(req data.CreateOrderRequest) error {
	if strings.TrimSpace(req.CustomerID) == "" {
		return fmt.Errorf("customer_id must not be empty")
	}

	if strings.TrimSpace(req.Currency) == "" {
		return fmt.Errorf("currency must not be empty")
	}

	if len(req.Currency) != 3 {
		return fmt.Errorf("currency must be a 3-letter code")
	}

	if req.Total <= 0 {
		return fmt.Errorf("total must be greater than 0")
	}

	return nil
}

func validateUpdateOrderRequest(req data.UpdateOrderRequest) error {
	if req.Status == nil && req.Total == nil {
		return fmt.Errorf("at least one field must be provided")
	}

	if req.Status != nil {
		switch *req.Status {
		case "pending", "processing", "shipped", "completed", "cancelled":
		default:
			return fmt.Errorf("invalid status")
		}
	}

	if req.Total != nil && *req.Total <= 0 {
		return fmt.Errorf("total must be greater than 0")
	}

	return nil
}

package api

import (
	"fmt"
	"net/mail"
	"strings"

	"github.com/Ashish-Barmaiya/torus-demo-backends/internal/data"
)

func validateCreateUserRequest(req data.CreateUserRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return fmt.Errorf("name must not be empty")
	}

	if strings.TrimSpace(req.Email) == "" {
		return fmt.Errorf("email must not be empty")
	}

	if _, err := mail.ParseAddress(req.Email); err != nil {
		return fmt.Errorf("invalid email")
	}

	switch req.Plan {
	case "free", "pro", "enterprise":
		return nil
	default:
		return fmt.Errorf("invalid plan")
	}
}

func validateUpdateUserRequest(req data.UpdateUserRequest) error {
	if req.Name == nil &&
		req.Email == nil &&
		req.Plan == nil &&
		req.Status == nil {
		return fmt.Errorf("at least one field must be provided")
	}

	if req.Email != nil {
		if strings.TrimSpace(*req.Email) == "" {
			return fmt.Errorf("email must not be empty")
		}

		if _, err := mail.ParseAddress(*req.Email); err != nil {
			return fmt.Errorf("invalid email")
		}
	}

	if req.Plan != nil {
		switch *req.Plan {
		case "free", "pro", "enterprise":
		default:
			return fmt.Errorf("invalid plan")
		}
	}

	if req.Status != nil {
		switch *req.Status {
		case "active", "suspended":
		default:
			return fmt.Errorf("invalid status")
		}
	}

	return nil
}

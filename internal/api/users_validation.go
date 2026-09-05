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

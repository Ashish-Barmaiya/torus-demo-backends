package data

type User struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Plan   string `json:"plan"`
	Status string `json:"status"`
}

type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Plan  string `json:"plan"`
}

type Order struct {
	ID         string `json:"id"`
	CustomerID string `json:"customer_id"`
	Status     string `json:"status"`
	Currency   string `json:"currency"`
	Total      int64  `json:"total"`
}

type Meta struct {
	Service  string `json:"service"`
	Instance string `json:"instance"`
}

type Response[T any] struct {
	Data T    `json:"data"`
	Meta Meta `json:"meta"`
}

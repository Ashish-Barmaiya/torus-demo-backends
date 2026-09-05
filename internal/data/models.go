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

type UpdateUserRequest struct {
	Name   *string `json:"name,omitempty"`
	Email  *string `json:"email,omitempty"`
	Plan   *string `json:"plan,omitempty"`
	Status *string `json:"status,omitempty"`
}

type Order struct {
	ID         string `json:"id"`
	CustomerID string `json:"customer_id"`
	Status     string `json:"status"`
	Currency   string `json:"currency"`
	Total      int64  `json:"total"`
}

type CreateOrderRequest struct {
	CustomerID string `json:"customer_id"`
	Currency   string `json:"currency"`
	Total      int64  `json:"total"`
}

type UpdateOrderRequest struct {
	Status *string `json:"status,omitempty"`
	Total  *int64  `json:"total,omitempty"`
}

type Meta struct {
	Service  string `json:"service"`
	Instance string `json:"instance"`
}

type Response[T any] struct {
	Data T    `json:"data"`
	Meta Meta `json:"meta"`
}

func (r Response[T]) ResponseData() any {
	return r.Data
}

func (r Response[T]) ResponseMeta() any {
	return r.Meta
}

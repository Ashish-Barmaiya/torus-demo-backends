package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/Ashish-Barmaiya/torus-demo-backends/internal/config"
	"github.com/Ashish-Barmaiya/torus-demo-backends/internal/data"
	"github.com/Ashish-Barmaiya/torus-demo-backends/internal/payload"
)

type jsonResponse interface {
	ResponseData() any
	ResponseMeta() any
}

type Server struct {
	cfg config.Config
}

func NewServer(cfg config.Config) *Server {
	return &Server{cfg: cfg}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", s.handleHealth)

	switch s.cfg.Service {
	case "users":
		mux.HandleFunc("/api/v1/users", s.handleUsers)
		mux.HandleFunc("/api/v1/users/", s.handleUser)

	case "orders":
		mux.HandleFunc("/api/v1/orders", s.handleOrders)
		mux.HandleFunc("/api/v1/orders/", s.handleOrder)

	default:
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "unknown service", http.StatusInternalServerError)
		})
	}

	return loggingMiddleware(s.cfg, mux)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, r, http.StatusOK, map[string]any{
		"status":   "ok",
		"service":  s.cfg.Service,
		"instance": s.cfg.Instance,
	})
}

func (s *Server) handleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		writeJSON(
			w,
			r,
			http.StatusOK,
			data.Response[[]data.User]{
				Data: data.Users(),
				Meta: data.Meta{
					Service:  s.cfg.Service,
					Instance: s.cfg.Instance,
				},
			},
		)

	case http.MethodPost:
		s.handleCreateUser(w, r)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req data.CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if err := validateCreateUserRequest(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := data.User{
		ID:     "usr_demo_created",
		Name:   req.Name,
		Email:  req.Email,
		Plan:   req.Plan,
		Status: "active",
	}

	writeJSON(
		w,
		r,
		http.StatusCreated,
		data.Response[data.User]{
			Data: user,
			Meta: data.Meta{
				Service:  s.cfg.Service,
				Instance: s.cfg.Instance,
			},
		},
	)
}

func (s *Server) handleUser(w http.ResponseWriter, r *http.Request) {
	rawID := strings.TrimPrefix(r.URL.Path, "/api/v1/users/")

	id, err := strconv.Atoi(rawID)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	user, ok := data.UserByID(id)
	if !ok {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		writeJSON(
			w,
			r,
			http.StatusOK,
			data.Response[data.User]{
				Data: user,
				Meta: data.Meta{
					Service:  s.cfg.Service,
					Instance: s.cfg.Instance,
				},
			},
		)

	case http.MethodPatch:
		s.handleUpdateUser(w, r, user)

	case http.MethodDelete:
		s.handleDeleteUser(w, r, user)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleUpdateUser(
	w http.ResponseWriter,
	r *http.Request,
	user data.User,
) {
	defer r.Body.Close()

	var req data.UpdateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if err := validateUpdateUserRequest(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Name != nil {
		user.Name = *req.Name
	}

	if req.Email != nil {
		user.Email = *req.Email
	}

	if req.Plan != nil {
		user.Plan = *req.Plan
	}

	if req.Status != nil {
		user.Status = *req.Status
	}

	writeJSON(
		w,
		r,
		http.StatusOK,
		data.Response[data.User]{
			Data: user,
			Meta: data.Meta{
				Service:  s.cfg.Service,
				Instance: s.cfg.Instance,
			},
		},
	)
}

func (s *Server) handleDeleteUser(
	w http.ResponseWriter,
	r *http.Request,
	user data.User,
) {
	writeJSON(
		w,
		r,
		http.StatusOK,
		data.Response[map[string]any]{
			Data: map[string]any{
				"id":      user.ID,
				"deleted": true,
			},
			Meta: data.Meta{
				Service:  s.cfg.Service,
				Instance: s.cfg.Instance,
			},
		},
	)
}

func (s *Server) handleOrders(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		writeJSON(
			w,
			r,
			http.StatusOK,
			data.Response[[]data.Order]{
				Data: data.Orders(),
				Meta: data.Meta{
					Service:  s.cfg.Service,
					Instance: s.cfg.Instance,
				},
			},
		)

	case http.MethodPost:
		s.handleCreateOrder(w, r)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleOrder(w http.ResponseWriter, r *http.Request) {
	rawID := strings.TrimPrefix(r.URL.Path, "/api/v1/orders/")

	id, err := strconv.Atoi(rawID)
	if err != nil {
		http.Error(w, "invalid order id", http.StatusBadRequest)
		return
	}

	order, ok := data.OrderByID(id)
	if !ok {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		writeJSON(
			w,
			r,
			http.StatusOK,
			data.Response[data.Order]{
				Data: order,
				Meta: data.Meta{
					Service:  s.cfg.Service,
					Instance: s.cfg.Instance,
				},
			},
		)

	case http.MethodPatch:
		s.handleUpdateOrder(w, r, order)

	case http.MethodDelete:
		s.handleDeleteOrder(w, r, order)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleCreateOrder(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req data.CreateOrderRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if err := validateCreateOrderRequest(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	order := data.Order{
		ID:         "ord_demo_created",
		CustomerID: req.CustomerID,
		Status:     "pending",
		Currency:   strings.ToUpper(req.Currency),
		Total:      req.Total,
	}

	writeJSON(
		w,
		r,
		http.StatusCreated,
		data.Response[data.Order]{
			Data: order,
			Meta: data.Meta{
				Service:  s.cfg.Service,
				Instance: s.cfg.Instance,
			},
		},
	)
}

func (s *Server) handleUpdateOrder(
	w http.ResponseWriter,
	r *http.Request,
	order data.Order,
) {
	defer r.Body.Close()

	var req data.UpdateOrderRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if err := validateUpdateOrderRequest(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Status != nil {
		order.Status = *req.Status
	}

	if req.Total != nil {
		order.Total = *req.Total
	}

	writeJSON(
		w,
		r,
		http.StatusOK,
		data.Response[data.Order]{
			Data: order,
			Meta: data.Meta{
				Service:  s.cfg.Service,
				Instance: s.cfg.Instance,
			},
		},
	)
}

func (s *Server) handleDeleteOrder(
	w http.ResponseWriter,
	r *http.Request,
	order data.Order,
) {
	writeJSON(
		w,
		r,
		http.StatusOK,
		data.Response[map[string]any]{
			Data: map[string]any{
				"id":      order.ID,
				"deleted": true,
			},
			Meta: data.Meta{
				Service:  s.cfg.Service,
				Instance: s.cfg.Instance,
			},
		},
	)
}

func writeJSON(
	w http.ResponseWriter,
	r *http.Request,
	status int,
	value any,
) {
	w.Header().Set("Content-Type", "application/json")

	responseSizeHeader := r.Header.Values(payload.ResponseSizeHeader)

	if len(responseSizeHeader) == 0 {
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(value)
		return
	}

	size, err := payload.ParseSize(responseSizeHeader[0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if size == payload.SizeEmpty {
		w.WriteHeader(status)
		return
	}

	response, ok := value.(jsonResponse)
	if !ok {
		http.Error(
			w,
			"response does not support payload generation",
			http.StatusInternalServerError,
		)
		return
	}

	body, err := payload.GenerateJSONWithData(
		size,
		response.ResponseData(),
		response.ResponseMeta(),
	)
	if err != nil {
		http.Error(
			w,
			"failed to generate response payload",
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(status)
	_, _ = w.Write(body)
}

func loggingMiddleware(cfg config.Config, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		_ = cfg
	})
}

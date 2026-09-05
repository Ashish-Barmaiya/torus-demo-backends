package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/Ashish-Barmaiya/torus-demo-backends/internal/config"
	"github.com/Ashish-Barmaiya/torus-demo-backends/internal/data"
)

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
	writeJSON(w, http.StatusOK, map[string]any{
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
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	writeJSON(
		w,
		http.StatusOK,
		data.Response[[]data.Order]{
			Data: data.Orders(),
			Meta: data.Meta{
				Service:  s.cfg.Service,
				Instance: s.cfg.Instance,
			},
		},
	)
}

func (s *Server) handleOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

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

	writeJSON(
		w,
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

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(value)
}

func loggingMiddleware(cfg config.Config, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		_ = cfg
	})
}

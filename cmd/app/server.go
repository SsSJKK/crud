package app

import (
	"net/http"

	"github.com/SsSJKK/crud/cmd/app/middleware"

	"github.com/gorilla/mux"

	"github.com/SsSJKK/crud/pkg/customers"
	"github.com/SsSJKK/crud/pkg/managers"
	"github.com/SsSJKK/crud/pkg/security"
)

//Server ...
type Server struct {
	mux          *mux.Router
	customersSvc *customers.Service
	securitySvc  *security.Service
	managersSvc  *managers.Service
}

//NewServer ...
func NewServer(m *mux.Router, cSvc *customers.Service, sSvc *security.Service, mSvc *managers.Service) *Server {
	return &Server{mux: m, customersSvc: cSvc, securitySvc: sSvc, managersSvc: mSvc}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

const (
	//GET ...
	GET = "GET"
	//POST ...
	POST = "POST"
	//DELETE ...
	DELETE = "DELETE"
)

//Init ...
func (s *Server) Init() {
	//customersAythMd := middleware.Authenticate(s.customersSvc.IDByToken)

	customersSubRouter := s.mux.PathPrefix("/api/customers").Subrouter()
	//customersSubRouter.Use(customersAythMd)

	customersSubRouter.HandleFunc("/active", s.handleGetAllActiveCustomers).Methods(GET)
	customersSubRouter.HandleFunc("", s.handleGetAllCustomers).Methods(GET)
	customersSubRouter.HandleFunc("/{id}", s.handleGetCustomerByID).Methods(GET)
	//customersSubRouter.HandleFunc("", s.handleSave).Methods(POST)
	customersSubRouter.HandleFunc("/{id}", s.handleDelete).Methods(DELETE)
	customersSubRouter.HandleFunc("/{id}/block", s.handleBlockByID).Methods(POST)
	customersSubRouter.HandleFunc("/{id}/block", s.handleUnBlockByID).Methods(DELETE)
	customersSubRouter.HandleFunc("", s.apiSave).Methods(POST)
	customersSubRouter.HandleFunc("/token", s.apiToken).Methods(POST)
	customersSubRouter.HandleFunc("/token/validate", s.handleValidateToken).Methods(POST)
	customersSubRouter.HandleFunc("products", s.hCustGetProdeucts).Methods(GET)
	customersSubRouter.HandleFunc("/purchases", s.pass).Methods(GET)
	customersSubRouter.HandleFunc("/purchases", s.pass).Methods(POST)
	//s.mux.Use(middleware.Basic(s.securitySvc.Auth))

	managersAythMd := middleware.Authenticate(s.managersSvc.IDByToken)

	managersSubrouter := s.mux.PathPrefix("/api/managers").Subrouter()
	managersSubrouter.Use(managersAythMd)
	managersSubrouter.HandleFunc("", s.hManagerR).Methods(POST)
	managersSubrouter.HandleFunc("/token", s.apiTokenManager).Methods(POST)
	managersSubrouter.HandleFunc("/token/validate", s.pass).Methods(POST)
	managersSubrouter.HandleFunc("/sales", s.hGetSeles).Methods(GET)
	managersSubrouter.HandleFunc("/sales", s.hMakeSeles).Methods(POST)
	managersSubrouter.HandleFunc("/products", s.pass).Methods(GET)
	managersSubrouter.HandleFunc("/products", s.hChProduct).Methods(POST)
	managersSubrouter.HandleFunc("/customers", s.pass).Methods(GET)
	managersSubrouter.HandleFunc("/customers", s.pass).Methods(POST)
	managersSubrouter.HandleFunc("/customers/{id}", s.pass).Methods(DELETE)
}

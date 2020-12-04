package app

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/SsSJKK/gosql/pkg/customers"
)

//Server ...
type Server struct {
	mux          *http.ServeMux
	customersSvc *customers.Service
}

//NewServer ...
func NewServer(m *http.ServeMux, cSvc *customers.Service) *Server {
	return &Server{mux: m, customersSvc: cSvc}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

//Init ...
func (s *Server) Init() {
	s.mux.HandleFunc("/customers.getById", s.handleGetCustomerByID)
	s.mux.HandleFunc("/customers.getAll", s.handleGetAllCustomers)
	s.mux.HandleFunc("/customers.getAllActive", s.handleGetAllActiveCustomers)
	s.mux.HandleFunc("/customers.blockById", s.handleBlockByID)
	s.mux.HandleFunc("/customers.unblockById", s.handleUnBlockByID)
	s.mux.HandleFunc("/customers.removeById", s.handleDelete)
	s.mux.HandleFunc("/customers.save", s.handleSave)
}

func (s *Server) handleGetAllCustomers(w http.ResponseWriter, r *http.Request) {

	items, err := s.customersSvc.All(r.Context())
	if err != nil {
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}

	respondJSON(w, items)
}

func (s *Server) handleGetAllActiveCustomers(w http.ResponseWriter, r *http.Request) {

	items, err := s.customersSvc.AllActive(r.Context())
	if err != nil {
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}

	respondJSON(w, items)
}

func (s *Server) handleGetCustomerByID(w http.ResponseWriter, r *http.Request) {
	idP := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idP, 10, 64)
	if err != nil {
		errorWriter(w, http.StatusBadRequest, err)
		return
	}
	item, err := s.customersSvc.ByID(r.Context(), id)
	if errors.Is(err, customers.ErrNotFound) {
		errorWriter(w, http.StatusNotFound, err)
		return
	}
	if err != nil {
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}
	respondJSON(w, item)
}

func (s *Server) handleBlockByID(w http.ResponseWriter, r *http.Request) {
	idP := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idP, 10, 64)
	if err != nil {
		errorWriter(w, http.StatusBadRequest, err)
		return
	}
	item, err := s.customersSvc.ChangeActive(r.Context(), id, false)

	if errors.Is(err, customers.ErrNotFound) {
		errorWriter(w, http.StatusNotFound, err)
		return
	}
	if err != nil {
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}
	respondJSON(w, item)
}

func (s *Server) handleUnBlockByID(w http.ResponseWriter, r *http.Request) {
	idP := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idP, 10, 64)
	if err != nil {
		errorWriter(w, http.StatusBadRequest, err)
		return
	}
	item, err := s.customersSvc.ChangeActive(r.Context(), id, true)

	if errors.Is(err, customers.ErrNotFound) {
		errorWriter(w, http.StatusNotFound, err)
		return
	}
	if err != nil {
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}
	respondJSON(w, item)
}

func (s *Server) handleDelete(w http.ResponseWriter, r *http.Request) {
	idP := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idP, 10, 64)
	if err != nil {
		errorWriter(w, http.StatusBadRequest, err)
		return
	}
	item, err := s.customersSvc.Delete(r.Context(), id)

	if errors.Is(err, customers.ErrNotFound) {
		errorWriter(w, http.StatusNotFound, err)
		return
	}
	if err != nil {
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}
	respondJSON(w, item)
}

func (s *Server) handleSave(w http.ResponseWriter, r *http.Request) {
	idP := r.FormValue("id")
	name := r.FormValue("name")
	phone := r.FormValue("phone")

	id, err := strconv.ParseInt(idP, 10, 64)
	if err != nil {
		errorWriter(w, http.StatusBadRequest, err)
		return
	}
	if name == "" && phone == "" {
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	item := &customers.Customer{
		ID:    id,
		Name:  name,
		Phone: phone,
	}
	customer, err := s.customersSvc.Save(r.Context(), item)
	if err != nil {
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}
	respondJSON(w, customer)
}
func errorWriter(w http.ResponseWriter, httpSts int, err error) {
	log.Print(err)
	http.Error(w, http.StatusText(httpSts), httpSts)
}
func respondJSON(w http.ResponseWriter, iData interface{}) {
	data, err := json.Marshal(iData)
	if err != nil {
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		log.Print(err)
	}
}

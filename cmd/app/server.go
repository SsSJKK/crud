package app

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"

	"github.com/SsSJKK/crud/pkg/customers"
	"github.com/SsSJKK/crud/pkg/security"
)

//Server ...
type Server struct {
	mux          *mux.Router
	customersSvc *customers.Service
	securitySvc  *security.Service
}

//NewServer ...
func NewServer(m *mux.Router, cSvc *customers.Service, sSvc *security.Service) *Server {
	return &Server{mux: m, customersSvc: cSvc, securitySvc: sSvc}
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
	s.mux.HandleFunc("/customers/active", s.handleGetAllActiveCustomers).Methods(GET)
	s.mux.HandleFunc("/customers", s.handleGetAllCustomers).Methods(GET)
	s.mux.HandleFunc("/customers/{id}", s.handleGetCustomerByID).Methods(GET)
	s.mux.HandleFunc("/customers", s.handleSave).Methods(POST)
	s.mux.HandleFunc("/customers/{id}", s.handleDelete).Methods(DELETE)
	s.mux.HandleFunc("/customers/{id}/block", s.handleBlockByID).Methods(POST)
	s.mux.HandleFunc("/customers/{id}/block", s.handleUnBlockByID).Methods(DELETE)

	s.mux.HandleFunc("/api/customers", s.apiSave).Methods("POST")
	s.mux.HandleFunc("/api/customers/token", s.apiToken).Methods("POST")
	s.mux.HandleFunc("/api/customers/token/validate", s.handleValidateToken).Methods("POST")

	//s.mux.Use(middleware.Basic(s.securitySvc.Auth))

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
	idP, ok := mux.Vars(r)["id"]
	if !ok {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
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

func (s *Server) handleDelete(w http.ResponseWriter, r *http.Request) {
	idP, ok := mux.Vars(r)["id"]
	if !ok {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
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
	var item *customers.Customer
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		errorWriter(w, http.StatusBadRequest, err)
		return
	}
	customer, err := s.customersSvc.Save(r.Context(), item)
	respondJSON(w, customer)
}

func (s *Server) handleBlockByID(w http.ResponseWriter, r *http.Request) {
	idP, ok := mux.Vars(r)["id"]
	if !ok {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
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
	idP, ok := mux.Vars(r)["id"]
	if !ok {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
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

func (s *Server) apiSave(w http.ResponseWriter, r *http.Request) {
	var item *customers.Customer
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(item.Password), bcrypt.DefaultCost)
	if err != nil {
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	item.Password = string(hash)
	customer, err := s.customersSvc.APISave(r.Context(), item)
	if err != nil {
		errorWriter(w, http.StatusBadRequest, err)
		return
	}
	respondJSON(w, customer)
}

func (s *Server) apiToken(w http.ResponseWriter, r *http.Request) {
	var item *struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	token, err := s.securitySvc.TokenForCustomer(r.Context(), item.Login, item.Password)

	if err != nil {
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	respondJSON(w, map[string]interface{}{"status": "ok", "token": token})
}

func (s *Server) handleValidateToken(w http.ResponseWriter, r *http.Request) {
	var item *struct {
		Token string `json:"token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	id, err := s.securitySvc.AuthenticateCustomer(r.Context(), item.Token)

	if err != nil {
		status := http.StatusInternalServerError
		text := "internal error"
		if err == security.ErrNoSuchUser {
			status = http.StatusNotFound
			text = "not found"
		}
		if err == security.ErrExpireToken {
			status = http.StatusBadRequest
			text = "expired"
		}

		respondJSONWithCode(w, status, map[string]interface{}{"status": "fail", "reason": text})
		return
	}

	res := make(map[string]interface{})
	res["status"] = "ok"
	res["customerId"] = id

	respondJSONWithCode(w, http.StatusOK, res)
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


func respondJSONWithCode(w http.ResponseWriter, sts int, iData interface{}) {
	data, err := json.Marshal(iData)
	if err != nil {
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(sts)
	_, err = w.Write(data)
	if err != nil {
		log.Print(err)
	}
}
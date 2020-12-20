package app

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/SsSJKK/crud/cmd/app/middleware"

	"github.com/SsSJKK/crud/pkg/managers"
)

func (s *Server) hManagerR(w http.ResponseWriter, r *http.Request) {
	log.Println("hManagerR")
	var item *managers.Managers
	id, err := middleware.Authentication(r.Context())
	log.Println(id, err)
	if err != nil {
		errorWriter(w, http.StatusBadRequest, err)
		return
	}
	err = json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	log.Println(item)
	manager, err := s.managersSvc.Registration(r.Context(), item)
	if err != nil {
		errorWriter(w, http.StatusBadRequest, err)
		return
	}
	log.Println(item.ID)

	token, err := s.securitySvc.TokenWithOut(r.Context(), manager.ID)

	respondJSON(w, map[string]interface{}{"status": "ok", "token": token})

}

func (s *Server) apiTokenManager(w http.ResponseWriter, r *http.Request) {
	var item *struct {
		Login    string `json:"phone"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	token, err := s.securitySvc.TokenForManager(r.Context(), item.Login, item.Password)

	if err != nil {
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	respondJSON(w, map[string]interface{}{"status": "ok", "token": token})
}

func (s *Server) hChProduct(w http.ResponseWriter, r *http.Request) {
	var item *managers.Product
	id, err := middleware.Authentication(r.Context())
	log.Println(id, err)
	if err != nil {
		errorWriter(w, http.StatusBadRequest, err)
		return
	}
	err = json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	log.Println(item)
	product, err := s.managersSvc.ChangeProduct(r.Context(), item)
	if err != nil {
		errorWriter(w, http.StatusBadRequest, err)
		return
	}
	log.Println(item.ID)

	respondJSON(w, map[string]interface{}{"id": product.ID})

}

func (s *Server) hMakeSeles(w http.ResponseWriter, r *http.Request) {
	var SaleP *managers.SalePositions
	err := json.NewDecoder(r.Body).Decode(&SaleP)
	if err != nil {
		errorWriter(w, http.StatusBadRequest, err)
		return
	}
	id, err := middleware.Authentication(r.Context())
	err = s.managersSvc.MakeSele(r.Context(), SaleP, id)
	if err != nil {
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	respondJSON(w, map[string]interface{}{"id": id})
}

func (s *Server) hGetSeles(w http.ResponseWriter, r *http.Request) {
	id, err := middleware.Authentication(r.Context())
	if err != nil {
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	sum, err := s.managersSvc.GetSales(r.Context(), id)
	if err != nil {
		errorWriter(w, http.StatusBadRequest, err)
		return
	}
	respondJSON(w, map[string]interface{}{
		"manager_id":    id,
		"total": sum,
	})

}

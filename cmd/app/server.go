package app

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	//"time"

	"github.com/shodikhuja83/crud/pkg/customers"
	"github.com/shodikhuja83/crud/pkg/security"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type Server struct {
	mux          *mux.Router
	customersSvc *customers.Service
	security     *security.Service
}

func NewServer(mux *mux.Router, customersSvc *customers.Service, security *security.Service) *Server {
	return &Server{
		mux:          mux,
		customersSvc: customersSvc,
		security:     security}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

const (
	GET    = "GET"
	POST   = "POST"
	DELETE = "DELETE"
)

func (s *Server) Init() {
	// s.mux.HandleFunc("/customers.getById", s.handleGetCustomerByID)

	s.mux.HandleFunc("/customers", s.handleGetCustomerAll).Methods(GET)
	s.mux.HandleFunc("/customers", s.handleGetCustomerSave).Methods(POST)
	s.mux.HandleFunc("/customers/active", s.handleGetCustomerAllActive).Methods(GET)

	s.mux.HandleFunc("/customers/{id}", s.handleGetCustomerByID).Methods(GET)
	s.mux.HandleFunc("/customers/{id}/block", s.handleGetCustomerBlockByID).Methods(POST)
	s.mux.HandleFunc("/customers/{id}/block", s.handleGetCustomerUnBlockByID).Methods(DELETE)
	s.mux.HandleFunc("/customers/{id}", s.handleGetCustomerRemoveByID).Methods(DELETE)

	s.mux.HandleFunc("/api/customers", s.handleSave).Methods(POST)
	s.mux.HandleFunc("/api/customers/token", s.handleCreateToken).Methods("POST")
	s.mux.HandleFunc("/api/customers/token/validate", s.handleValidateToken).Methods("POST")

	//s.mux.Use(middleware.Basic(s.security.Auth))

}

func (s *Server) handleGetCustomerByID(w http.ResponseWriter, r *http.Request) {
	//idParam := r.URL.Query().Get("id")
	idParam, ok := mux.Vars(r)["id"]
	if !ok {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return
	}

	item, err := s.customersSvc.ByID(r.Context(), id)
	if errors.Is(err, customers.ErrNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(item)
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		return
	}

}
func (s *Server) handleGetCustomerAll(w http.ResponseWriter, r *http.Request) {

	item, err := s.customersSvc.All(r.Context())

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(item)
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		return
	}

}
func (s *Server) handleGetCustomerAllActive(w http.ResponseWriter, r *http.Request) {

	item, err := s.customersSvc.AllActive(r.Context())

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(item)
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		return
	}

}
func (s *Server) handleGetCustomerSave(w http.ResponseWriter, r *http.Request) {
	var item *customers.Customer
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	item, err = s.customersSvc.Save(r.Context(), item)

	data, err := json.Marshal(item)
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		return
	}

}
func (s *Server) handleGetCustomerRemoveByID(w http.ResponseWriter, r *http.Request) {
	idP := mux.Vars(r)["id"]

	id, err := strconv.ParseInt(idP, 10, 64)
	if err != nil {
		return
	}

	item, err := s.customersSvc.RemoveById(r.Context(), id)
	if errors.Is(err, customers.ErrNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(item)
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		return
	}

}
func (s *Server) handleGetCustomerBlockByID(w http.ResponseWriter, r *http.Request) {
	idP := mux.Vars(r)["id"]

	id, err := strconv.ParseInt(idP, 10, 64)
	if err != nil {
		return
	}

	item, err := s.customersSvc.BlockByID(r.Context(), id)
	if errors.Is(err, customers.ErrNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(item)
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		return
	}

}
func (s *Server) handleGetCustomerUnBlockByID(w http.ResponseWriter, r *http.Request) {
	idP := mux.Vars(r)["id"]

	id, err := strconv.ParseInt(idP, 10, 64)
	if err != nil {
		return
	}

	item, err := s.customersSvc.UnBlockByID(r.Context(), id)
	if errors.Is(err, customers.ErrNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(item)
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		return
	}

}

func (s *Server) handleSave(w http.ResponseWriter, r *http.Request) {

	var item *customers.Customer

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {

		errWriter(w, http.StatusBadRequest, err)
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(item.Password), bcrypt.DefaultCost)
	if err != nil {

		errWriter(w, http.StatusInternalServerError, err)
		return
	}

	item.Password = string(hashed)

	customer, err := s.customersSvc.Save(r.Context(), item)
	if err != nil {
		errWriter(w, http.StatusInternalServerError, err)
		return
	}
	resJson(w, customer)
}

func (s *Server) handleCreateToken(w http.ResponseWriter, r *http.Request) {
	var item *struct {
		Login    string `json:"login"`
		Password string `json:"password`
	}

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		errWriter(w, http.StatusBadRequest, err)
		return
	}

	token, err := s.security.TokenForCustomer(r.Context(), item.Login, item.Password)

	if err != nil {
		errWriter(w, http.StatusBadRequest, err)
		return
	}

	resJson(w, map[string]interface{}{"status": "ok", "token": token})
}

func (s *Server) handleValidateToken(w http.ResponseWriter, r *http.Request) {
	var item *struct {
		Token string `json:"token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		errWriter(w, http.StatusBadRequest, err)
		return
	}

	id, err := s.security.AuthenticateCustomer(r.Context(), item.Token)

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

		resJsonWithCode(w, status, map[string]interface{}{"status": "fail", "reason": text})
		return
	}

	res := make(map[string]interface{})
	res["status"] = "ok"
	res["customerId"] = id

	resJsonWithCode(w, http.StatusOK, res)
}

func resJson(w http.ResponseWriter, iData interface{}) {

	data, err := json.Marshal(iData)

	if err != nil {
		errWriter(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)

	if err != nil {

		log.Print(err)
	}
}

func resJsonWithCode(w http.ResponseWriter, sts int, iData interface{}) {

	data, err := json.Marshal(iData)

	if err != nil {

		errWriter(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(sts)

	_, err = w.Write(data)

	if err != nil {

		log.Print(err)
	}
}

// this is the function for writing the error to responseWriter
func errWriter(w http.ResponseWriter, httpSts int, err error) {
	log.Print(err)
	http.Error(w, http.StatusText(httpSts), httpSts)
}

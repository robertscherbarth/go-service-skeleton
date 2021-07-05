package ports

import (
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"net/http"
)

type Service interface {
	Add(user User) error
	Delete(id string) error
	FindAll() ([]User, error)
	FindByID(id string) (User, error)
}

type Http struct {
	logger  *zap.Logger
	service Service
}

func NewHttp(logger *zap.Logger, service Service) *Http {
	return &Http{logger: logger, service: service}
}

func (h *Http) FindUserByID(w http.ResponseWriter, r *http.Request, id string) {
	user, err := h.service.FindByID(id)
	if err != nil {
		h.logger.Error(err.Error(), zap.Error(err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		h.logger.Error(err.Error(), zap.Error(err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *Http) DeleteUser(w http.ResponseWriter, _ *http.Request, id string) {
	err := h.service.Delete(id)
	if err != nil {
		h.logger.Error(err.Error(), zap.Error(err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Http) FindAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.FindAll()
	if err != nil {
		h.logger.Error(err.Error(), zap.Error(err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		h.logger.Error(err.Error(), zap.Error(err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *Http) AddUser(w http.ResponseWriter, r *http.Request) {
	var user User
	body := r.Body
	defer body.Close()

	err := json.NewDecoder(body).Decode(&user)
	if err != nil {
		msg := h.enrichErrorMsg(err)
		h.logger.Error(msg, zap.Error(err))
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	err = h.service.Add(user)
	if err != nil {
		h.logger.Error(err.Error(), zap.Error(err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Http) enrichErrorMsg(err error) string {
	var syntaxError *json.SyntaxError
	var unmarshalTypeError *json.UnmarshalTypeError
	var msg string

	switch {
	case errors.As(err, &syntaxError):
		msg = fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
	case errors.As(err, &unmarshalTypeError):
		msg = fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
	default:
		msg = http.StatusText(http.StatusBadRequest)
	}

	return msg
}

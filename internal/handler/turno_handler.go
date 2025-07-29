package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/IvanMiranda1/gestion-de-turnos-peluqueria-unipersonal/internal/dto"
	"github.com/IvanMiranda1/gestion-de-turnos-peluqueria-unipersonal/internal/service/turno"
	"github.com/IvanMiranda1/gestion-de-turnos-peluqueria-unipersonal/pkg/web"
	"github.com/go-chi/chi/v5"
)

type TurnoHandler struct {
	s turno.TurnoService
}

func NewTurnoHandler(s turno.TurnoService) *TurnoHandler {
	return &TurnoHandler{s: s}
}

func (h *TurnoHandler) RegisterRoutes(r chi.Router) {
	r.Post("/", h.Create)
	r.Put("/{id}", h.Update)
	r.Get("/{fecha}", h.GetByFecha)
	r.Get("/", h.GetAll) //GET /turno
	r.Delete("/{id}", h.Delete)
}

func (h *TurnoHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.TurnoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		web.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	t, err := h.s.ToDomain(r.Context(), &req)
	if err != nil {
		web.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	res, err := h.s.Create(r.Context(), t)
	if err != nil {
		web.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	web.Success(w, http.StatusCreated, dto.TurnoFromDomain(res))
}

func (h *TurnoHandler) Update(w http.ResponseWriter, r *http.Request) {
	var req dto.TurnoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		web.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	t, err := h.s.ToDomain(r.Context(), &req)
	if err != nil {
		web.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	res, err := h.s.Update(r.Context(), t)
	if err != nil {
		web.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	web.Success(w, http.StatusOK, dto.TurnoFromDomain(res))
}

func (h *TurnoHandler) GetByFecha(w http.ResponseWriter, r *http.Request) {
	fecha := chi.URLParam(r, "fecha")
	if fecha == "" {
		web.Error(w, http.StatusBadRequest, "fecha is required")
		return
	}
	fechaparsed, err := time.Parse("2006/01/02", fecha)
	if err != nil {
		web.Error(w, http.StatusBadRequest, "formato de fecha invalido")
		return
	}
	res, err := h.s.GetByFecha(r.Context(), fechaparsed)
	if err != nil {
		web.Error(w, http.StatusNotFound, err.Error())
		return
	}
	turnoSlice := make([]any, 0, len(res))
	for _, t := range res {
		turnoSlice = append(turnoSlice, dto.TurnoFromDomain(t))
	}

	web.Success(w, http.StatusOK, turnoSlice)
}

func (h *TurnoHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	res, err := h.s.GetAll(r.Context())
	if err != nil {
		web.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	turnoSlice := make([]any, 0, len(res))
	for _, t := range res {
		turnoSlice = append(turnoSlice, dto.TurnoFromDomain(t))
	}
	web.Success(w, http.StatusOK, turnoSlice)
}

func (h *TurnoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		web.Error(w, http.StatusBadRequest, "id is required")
		return
	}
	if err := h.s.Delete(r.Context(), id); err != nil {
		web.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

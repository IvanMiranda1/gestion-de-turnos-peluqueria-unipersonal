package dto

import (
	"github.com/IvanMiranda1/gestion-de-turnos-peluqueria-unipersonal/internal/domain"
)

type ClienteRequest struct {
	ID                 string `json:"id" validate:"required"`
	Nombre             string `json:"nombre" validate:"required"`
	Telefono           string `json:"telefono" validate:"required"`
	PreferenciaHoraria string `json:"preferenciaHoraria" validate:"required"`
}

func (r *ClienteRequest) ToDomain() (*domain.Cliente, error) {
	preferenciaHoraria, err := domain.ParsePreferenciaHoraria(r.PreferenciaHoraria)
	if err != nil {
		return nil, err
	}
	return domain.NewCliente(
		r.ID,
		r.Nombre,
		r.Telefono,
		preferenciaHoraria,
	), nil
}

type ClienteResponse struct {
	ID                 string `json:"id"`
	Nombre             string `json:"nombre"`
	Telefono           string `json:"telefono"`
	PreferenciaHoraria string `json:"preferenciaHoraria"`
}

func ClienteFromDomain(c *domain.Cliente) *ClienteResponse {
	return &ClienteResponse{
		ID:                 c.ID,
		Nombre:             c.Nombre,
		Telefono:           c.Telefono,
		PreferenciaHoraria: c.PreferenciaHoraria.String(),
	}
}

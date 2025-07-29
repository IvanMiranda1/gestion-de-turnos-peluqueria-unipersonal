package dto

import (
	"github.com/IvanMiranda1/gestion-de-turnos-peluqueria-unipersonal/internal/domain"
)

type TurnoRequest struct {
	ID        string `json:"id" validate:"required"`
	Fecha     string `json:"fecha" validate:"required"`
	Hora      string `json:"hora" validate:"required"`
	ClienteID string `json:"clienteID" validate:"required"`
}

type TurnoResponse struct {
	ID        string `json:"id"`
	Fecha     string `json:"fecha"`
	Hora      string `json:"hora"`
	ClienteID string `json:"clienteID"`
}

func TurnoFromDomain(t *domain.Turno) *TurnoResponse {
	return &TurnoResponse{
		ID:        t.ID,
		Fecha:     t.Fecha.String(),
		Hora:      t.Hora.String(),
		ClienteID: t.Cliente.ID,
	}
}

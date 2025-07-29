package domain

import "errors"

type Cliente struct {
	ID                 string
	Nombre             string
	Telefono           string
	PreferenciaHoraria PreferenciaHoraria
}

func NewCliente(id, nombre, telefono string, preferenciahoraria PreferenciaHoraria) *Cliente {
	return &Cliente{
		ID:                 id,
		Nombre:             nombre,
		Telefono:           telefono,
		PreferenciaHoraria: preferenciahoraria,
	}
}

func (c *Cliente) Validate() error {
	if c.Nombre == "" || c.Telefono == "" || !IsValidPreferenciaHoraria(c.PreferenciaHoraria) {
		return errors.New("campos no válidos")
	}
	return nil
}

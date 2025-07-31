package domain

import (
	"errors"
	"fmt"
	"time"
)

type Turno struct {
	ID      string
	Fecha   time.Time
	Hora    TimeOfDay
	Cliente Cliente
}

func NewTurno(id string, fecha time.Time, hora TimeOfDay, cliente Cliente) *Turno {
	return &Turno{
		ID:      id,
		Fecha:   fecha,
		Hora:    hora,
		Cliente: cliente,
	}
}

func (t *Turno) Validate() error {
	if t.Fecha.IsZero() {
		return errors.New("fecha no puede ser cero")
	}
	if !t.Hora.IsValid() {
		return errors.New("hora inválida")
	}
	if err := t.Cliente.Validate(); err != nil {
		return fmt.Errorf("cliente inválido: %w", err)
	}
	return nil
}

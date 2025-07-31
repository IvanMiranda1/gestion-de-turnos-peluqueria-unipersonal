package repository

import (
	"context"
	"time"

	"github.com/IvanMiranda1/gestion-de-turnos-peluqueria-unipersonal/internal/domain"
)

type ClienteRepository interface {
	CreateOrUpdate(ctx context.Context, c *domain.Cliente) (*domain.Cliente, error)
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*domain.Cliente, error)
	GetAll(ctx context.Context) ([]*domain.Cliente, error)
}

type TurnoRepository interface {
	CreateOrUpdate(ctx context.Context, t *domain.Turno) (*domain.Turno, error)
	Delete(ctx context.Context, id string) error
	GetByFecha(ctx context.Context, fecha time.Time) ([]*domain.Turno, error)
	GetAll(ctx context.Context) ([]*domain.Turno, error)
}

/*
ctx context.Context es un objeto que transporta información de control a través de llamadas. Para cliente:

— El contexto permite cancelar operaciones si el usuario cierra la conexión o se agota el tiempo.
— También permite pasar valores auxiliares, como ID de usuario autenticado o logs.
— Es una forma de controlar el ciclo de vida de una operación completa.

Ejemplo: si una petición HTTP al crear un cliente tarda mucho, ctx puede avisar al repositorio que interrumpa la consulta a la base.

Sin contexto: la operación sigue aunque el cliente ya se fue.
Con contexto: se puede cortar a mitad, ahorrando recursos.

Siempre se propaga desde el handler hacia abajo.
*/

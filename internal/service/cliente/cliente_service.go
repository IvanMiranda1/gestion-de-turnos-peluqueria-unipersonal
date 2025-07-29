package cliente

import (
	"context"
	"errors"

	"github.com/IvanMiranda1/gestion-de-turnos-peluqueria-unipersonal/internal/domain"
	"github.com/IvanMiranda1/gestion-de-turnos-peluqueria-unipersonal/repository"
	"github.com/google/uuid"
)

type ClienteService interface {
	Create(ctx context.Context, c *domain.Cliente) (*domain.Cliente, error)
	Update(ctx context.Context, c *domain.Cliente) (*domain.Cliente, error)
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*domain.Cliente, error)
	GetAll(ctx context.Context) ([]*domain.Cliente, error)
}

type clienteService struct {
	repo repository.ClienteRepository
}

func NewClienteService(repo repository.ClienteRepository) *clienteService {
	return &clienteService{repo: repo}
}

func (s clienteService) Create(ctx context.Context, c *domain.Cliente) (*domain.Cliente, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}
	if c.ID == "" {
		c.ID = uuid.New().String()
	}

	return s.repo.CreateOrUpdate(ctx, c)

}

func (s clienteService) Update(ctx context.Context, c *domain.Cliente) (*domain.Cliente, error) {
	if c.ID == "" {
		return nil, errors.New("ID requerido para actualizar")
	}
	if err := c.Validate(); err != nil {
		return nil, err
	}
	return s.repo.CreateOrUpdate(ctx, c)
}

func (s clienteService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s clienteService) GetByID(ctx context.Context, id string) (*domain.Cliente, error) {
	return s.repo.GetByID(ctx, id)
}

func (s clienteService) GetAll(ctx context.Context) ([]*domain.Cliente, error) {
	return s.repo.GetAll(ctx)
}

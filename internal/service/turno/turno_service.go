package turno

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/IvanMiranda1/gestion-de-turnos-peluqueria-unipersonal/internal/domain"
	"github.com/IvanMiranda1/gestion-de-turnos-peluqueria-unipersonal/internal/dto"
	"github.com/IvanMiranda1/gestion-de-turnos-peluqueria-unipersonal/internal/repository"
	service "github.com/IvanMiranda1/gestion-de-turnos-peluqueria-unipersonal/internal/service/cliente"
	"github.com/google/uuid"
)

type TurnoService interface {
	Create(ctx context.Context, t *domain.Turno) (*domain.Turno, error)
	Update(ctx context.Context, t *domain.Turno) (*domain.Turno, error)
	Delete(ctx context.Context, id string) error
	GetByFecha(ctx context.Context, fecha time.Time) ([]*domain.Turno, error)
	GetAll(ctx context.Context) ([]*domain.Turno, error)
	ToDomain(ctx context.Context, t *dto.TurnoRequest) (*domain.Turno, error)
}

type turnoService struct {
	repo           repository.TurnoRepository
	clienteService service.ClienteService
}

func NewTurnoService(repo repository.TurnoRepository, cs service.ClienteService) *turnoService {
	return &turnoService{
		repo:           repo,
		clienteService: cs,
	}
}

func (s turnoService) Create(ctx context.Context, t *domain.Turno) (*domain.Turno, error) {
	if err := t.Validate(); err != nil {
		return nil, err
	}
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	return s.repo.CreateOrUpdate(ctx, t)
}

func (s turnoService) Update(ctx context.Context, t *domain.Turno) (*domain.Turno, error) {
	if err := t.Validate(); err != nil {
		return nil, err
	}
	if t.ID == "" {
		return nil, errors.New("ID requerido para actualizar")
	}
	return s.repo.CreateOrUpdate(ctx, t)
}

func (s turnoService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s turnoService) GetByFecha(ctx context.Context, fecha time.Time) ([]*domain.Turno, error) {
	return s.repo.GetByFecha(ctx, fecha)
}

func (s turnoService) GetAll(ctx context.Context) ([]*domain.Turno, error) {
	return s.repo.GetAll(ctx)
}

// este no recibe test
func (s turnoService) ToDomain(ctx context.Context, t *dto.TurnoRequest) (*domain.Turno, error) {
	fecha, err := time.Parse("2006/01/02", t.Fecha)
	if err != nil {
		return nil, fmt.Errorf("error de parse de fecha time.Time: %w", err)
	}
	hora, err := domain.ParseTimeOfDay(t.Hora)
	if err != nil {
		return nil, fmt.Errorf("error de parse timeofday: %w", err)
	}
	cliente, _ := s.clienteService.GetByID(ctx, t.ClienteID)
	if cliente == nil {
		return nil, fmt.Errorf("cliente vacio")
	}
	return domain.NewTurno(
		t.ID,
		fecha,
		hora,
		*cliente,
	), nil

}

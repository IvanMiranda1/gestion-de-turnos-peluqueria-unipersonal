package turno_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/IvanMiranda1/gestion-de-turnos-peluqueria-unipersonal/internal/domain"
	"github.com/IvanMiranda1/gestion-de-turnos-peluqueria-unipersonal/internal/service/turno"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTurnoRepository struct {
	mock.Mock
}

func (m *MockTurnoRepository) CreateOrUpdate(ctx context.Context, t *domain.Turno) (*domain.Turno, error) {
	args := m.Called(ctx, t)
	if args.Get(0) != nil {
		return args.Get(0).(*domain.Turno), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTurnoRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockTurnoRepository) GetByFecha(ctx context.Context, fecha time.Time) ([]*domain.Turno, error) {
	args := m.Called(ctx, fecha)
	if args.Get(0) != nil {
		return args.Get(0).([]*domain.Turno), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTurnoRepository) GetAll(ctx context.Context) ([]*domain.Turno, error) {
	args := m.Called(ctx)
	if args.Get(0) != nil {
		return args.Get(0).([]*domain.Turno), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestTurnoService_Create(t *testing.T) {
	t.Run("Create Return Error Validate()", func(t *testing.T) {
		s := turno.NewTurnoService(nil, nil)
		res, err := s.Create(context.Background(), &domain.Turno{
			ID:    "",
			Fecha: time.Time{},
			Hora:  domain.TimeOfDay{},
		})
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.EqualError(t, err, "fecha no puede ser cero")
	})
	t.Run("Create asigna UUID si ID esta vacio", func(t *testing.T) {
		mockRepo := new(MockTurnoRepository)
		s := turno.NewTurnoService(mockRepo, nil)
		fechaprueba, _ := time.Parse("2006/01/02", "2025/08/15")
		horaprueba, _ := domain.ParseTimeOfDay("10:30")
		turnoNuevo := &domain.Turno{
			ID:    "",
			Fecha: fechaprueba,
			Hora:  horaprueba,
			Cliente: domain.Cliente{
				ID:                 "",
				Nombre:             "Cliente Test",
				Telefono:           "123456789",
				PreferenciaHoraria: domain.PreferenciaHoraria(1),
			},
		}
		//mock anything espera cualquier valor de tipo *domain.Turno
		mockRepo.On("CreateOrUpdate", mock.Anything, mock.AnythingOfType("*domain.Turno")).Return(turnoNuevo, nil).Run(func(args mock.Arguments) {
			turnoNuevo := args.Get(1).(*domain.Turno)
			assert.NotEmpty(t, turnoNuevo.ID) // verifica que el ID no esté vacío, si no da error
		})
		s.Create(context.Background(), turnoNuevo)
	})
	//test driven table
	tests := []struct {
		name     string
		mockData *domain.Turno
		mockErr  error
		WantErr  bool
	}{
		{"Success", makeTurno("01"), nil, false},
		{"RepoError", makeTurno("01"), assert.AnError, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, mockRepo := setupTurnoServiceWithMock(t)
			mockRepo.On("CreateOrUpdate", mock.Anything, tt.mockData).Return(tt.mockData, tt.mockErr)
			got, err := s.Create(context.Background(), tt.mockData)

			if tt.WantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockData, got)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestTurnoService_Update(t *testing.T) {
	t.Run("Update Return Error Validate()", func(t *testing.T) {
		s := turno.NewTurnoService(nil, nil)
		res, err := s.Update(context.Background(), &domain.Turno{
			ID:    "",
			Fecha: time.Time{},
			Hora:  domain.TimeOfDay{},
			Cliente: domain.Cliente{
				ID:                 "123",
				Nombre:             "Cliente Test",
				Telefono:           "123456789",
				PreferenciaHoraria: domain.PreferenciaHoraria(1),
			},
		})
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.EqualError(t, err, "fecha no puede ser cero")
	})
	t.Run("Validar error si ID esta vacio", func(t *testing.T) {
		s := turno.NewTurnoService(nil, nil)
		res, err := s.Update(context.Background(), &domain.Turno{
			ID:    "",
			Fecha: time.Now(),
			Hora:  domain.TimeOfDay{},
			Cliente: domain.Cliente{
				ID:                 "123",
				Nombre:             "Cliente Test",
				Telefono:           "123456789",
				PreferenciaHoraria: domain.PreferenciaHoraria(1),
			},
		})
		assert.NoError(t, err)
		assert.NotEmpty(t, res.ID)
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.EqualError(t, err, "ID requerido para actualizar")
	})

	//Table Driven Tests
	tests := []struct {
		name     string
		mockData *domain.Turno
		mockErr  error
		WantErr  bool
	}{
		{"Success", makeTurno("01"), nil, false},
		{"RepoError", makeTurno("01"), assert.AnError, true},
		{"EmptyResult", makeTurno("01"), nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, mockRepo := setupTurnoServiceWithMock(t)
			mockRepo.On("CreateOrUpdate", mock.Anything, tt.mockData).Return(tt.mockData, tt.mockErr)
			got, err := s.Update(context.Background(), tt.mockData)

			if tt.WantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockData, got)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestTurnoService_Delete(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"Delete Return Exitoso", false},
		{"Delete Return Error", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTurnoRepository)
			s := turno.NewTurnoService(mockRepo, nil)
			if tt.wantErr {
				mockRepo.On("Delete", mock.Anything, "123").Return(assert.AnError)
			} else {
				mockRepo.On("Delete", mock.Anything, "123").Return(nil)
			}
			err := s.Delete(context.Background(), "123")
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, assert.AnError, err)
			} else {
				assert.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestTurnoService_GetByFecha(t *testing.T) {
	tests := []struct {
		name     string
		mockData []*domain.Turno
		mockErr  error
		wantErr  bool
	}{
		{"Success", []*domain.Turno{makeTurno("15"), makeTurno("16")}, nil, false},
		{"RepoError", nil, errors.New("repo failure"), true},
		{"EmptyResult", []*domain.Turno{}, nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, mockRepo := setupTurnoServiceWithMock(t)
			fechaprueba, _ := time.Parse("2006/01/02", "2025/08/15")
			mockRepo.On("GetByFecha", mock.Anything, fechaprueba).Return(tt.mockData, tt.mockErr)

			got, err := s.GetByFecha(context.Background(), fechaprueba)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockData, got)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestTurnoService_GetAll(t *testing.T) {
	/* mockData: lo que devuelve el mock si no hay error. Simula la respuesta esperada del repositorio ([]*domain.Turno en este caso).

	mockErr: el error que el mock debe simular. Si es distinto de nil, simula un fallo del repositorio.

	wantErr: expectativa del test. true si esperás que falle el método GetAll, false si esperás éxito. */
	tests := []struct {
		name     string
		mockData []*domain.Turno
		mockErr  error
		wantErr  bool
	}{
		/* se definen los tests en  orden
		name = success
		mockdata = []*domain.Turno{makeTurno("01"), makeTurno("02")}
		mockErr = nil
		wantErr = false
		*/
		{"Success", []*domain.Turno{makeTurno("01"), makeTurno("02")}, nil, false},
		{"RepoError", nil, errors.New("repo failure"), true},
		{"EmptyResult", []*domain.Turno{}, nil, false},
	} //los test definidos se ejecutan en el orden en que fueron creados

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, mockRepo := setupTurnoServiceWithMock(t)
			mockRepo.On("GetAll", mock.Anything).Return(tt.mockData, tt.mockErr)

			got, err := s.GetAll(context.Background())

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockData, got)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// funciones auxiliares
func makeTurno(dia string) *domain.Turno {
	fecha, _ := time.Parse("2006/01/02", fmt.Sprintf("2025/06/%s", dia))
	return &domain.Turno{
		ID:    uuid.NewString(),
		Fecha: fecha,
		Hora:  domain.TimeOfDay{Hour: 10, Minute: 30},
		Cliente: domain.Cliente{
			ID:                 "123",
			Nombre:             "Cliente Test",
			Telefono:           "123456789",
			PreferenciaHoraria: domain.PreferenciaHoraria(1),
		},
	}
}

func setupTurnoServiceWithMock(t *testing.T) (turno.TurnoService, *MockTurnoRepository) {
	mockRepo := new(MockTurnoRepository)
	s := turno.NewTurnoService(mockRepo, nil)
	return s, mockRepo
}

/* explicacion de mocktest

Ubicación del foco: pantalla dividida entre implementación (turno/service.go) y prueba (turno/service_test.go).

Identificación del objetivo: observar el método del servicio (Create, Update, etc.) y recorrer su lógica línea por línea.

Detección de validaciones: por cada if err := validar(...), marcar que se requiere un t.Run("error si X es inválido").

Estructura de prueba espejo: cada bloque if o decisión lógica relevante en el método del servicio debe tener su correspondiente t.Run() en el test.

Nombre de subtests: usar descripciones claras que indiquen qué parte de la lógica se está verificando.

Mocks obligatorios: si el servicio llama al repositorio, se mockea. Solo es necesario definir el On(...).Return(...) si se espera que esa rama se ejecute en ese test.

Cobertura lógica:

Test 1: validar error si ID está vacío.

Test 2: validar error si Fecha es cero.

Test 3: validar éxito si todo es correcto.

Test 4: validar transformación si hay lógica que modifique el input antes de guardar.

Test 5: validar errores devueltos por el repositorio si corresponde.

Assert y control: usar assert.Error o assert.NoError, assert.Equal, assert.Nil, assert.NotEmpty según lo que corresponda a cada test.

Minimización visual: si cada t.Run cubre una lógica aislada, se puede colapsar para ver solo los títulos, verificando que haya un espejo entre la lógica del servicio y los tests.

Desacople de dependencias: nunca testear con una base de datos real, siempre con mocks que permitan inyectar respuestas controladas.

Resultado final: la estructura del test se convierte en una representación declarativa de las decisiones y rutas lógicas del servicio. Cada ruta esperada o error anticipado debe estar cubierto.
*/

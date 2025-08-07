package cliente_test

import (
	"context"
	"errors"
	"testing"

	"github.com/IvanMiranda1/gestion-de-turnos-peluqueria-unipersonal/internal/domain"
	cliente "github.com/IvanMiranda1/gestion-de-turnos-peluqueria-unipersonal/internal/service/cliente"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockClienteRepository struct {
	mock.Mock
}

func (m *MockClienteRepository) CreateOrUpdate(ctx context.Context, c *domain.Cliente) (*domain.Cliente, error) {
	// m.Called guarda el llamado con sus argumentos y devuelve lo configurado en el test
	args := m.Called(ctx, c)
	if args.Get(0) != nil {
		// args.Get(0) es el primer valor que se configuró como retorno (el cliente creado)
		// args.Error(1) es el segundo valor, que sería un error
		return args.Get(0).(*domain.Cliente), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockClienteRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockClienteRepository) GetByID(ctx context.Context, id string) (*domain.Cliente, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*domain.Cliente), args.Error(1)
	}
	return nil, args.Error(1)
	// args.Get(0).(*domain.Cliente) type assertion, interpretando el primer valor como *domain.Cliente
	// el segundo valor es el error, args.Error(1), es como hacer el returno cliente, nil
}

func (m *MockClienteRepository) GetAll(ctx context.Context) ([]*domain.Cliente, error) {
	args := m.Called(ctx)
	if args.Get(0) != nil {
		return args.Get(0).([]*domain.Cliente), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestClienteService_Create(t *testing.T) {
	t.Run("Error validate()", func(t *testing.T) {
		s, _ := setupClienteServiceWithMock(t)
		res, err := s.Create(context.Background(), makeCliente("123", "")) //nombre vacío
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.EqualError(t, err, "campos no válidos")
	})
	t.Run("Asigna UUID  si ID esta vacío", func(t *testing.T) {
		s, _, := setupClienteServiceWithMock(t)
		res, err := s.Create(context.Background(), makeCliente("", "Ivan"))
		assert.NoError(t, err)
		assert.NotEmpty(t, res.ID)
	})

	tests := []struct {
		name string
		mockData *domain.Cliente
		mockErr  error
		WantErr  bool
	}{
		{"Success", makeCliente("01", "Pepe"), nil, false},
		{"RepoError", makeCliente("01", "Pepe"), assert.AnError, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, mockRepo := setupClienteServiceWithMock(t)
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

func TestClienteService_Update(t *testing.T) {
	t.Run("Return error si ID está vacío", func(t *testing.T) {
		s, _ := setupClienteServiceWithMock(t)
		res, err := s.Update(context.Background(), makeCliente("", "Pepe"))
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.EqualError(t, err, "ID requerido para actualizar")
	})

	t.Run("Error validate()", func(t *testing.T) {
		s, _ := setupClienteServiceWithMock(t)
		res, err := s.Update(context.Background(), makeCliente("123", "")) //nombre vacío
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.EqualError(t, err, "campos no válidos")
	})
	
	//Table Driven Tests
	tests := []struct {
		name string
		mockData *domain.Cliente
		mockErr  error
		WantErr  bool
	}{
		{"Success", makeCliente("01", "Pepe"), nil, false},
		{"RepoError", makeCliente("01", "Pepe"), assert.AnError, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, mockRepo := setupClienteServiceWithMock(t)
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

func TestClienteService_Delete(t *testing.T) {
	tests := []struct {
		name    string
		mockID  string
		mockErr error
		WantErr bool
	}{
		{"Success", "123", nil, false},
		{"RepoError", "123", assert.AnError, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, mockRepo := setupClienteServiceWithMock(t)
			mockRepo.On("Delete", mock.Anything, tt.mockID).Return(tt.mockErr)
			err := s.Delete(context.Background(), tt.mockID)

			if tt.WantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestClienteService_GetByID(t *testing.T) {
	tests := []struct {
		name     string
		mockID   string
		mockData *domain.Cliente
		mockErr  error
		WantErr  bool
	}{
		{"Success", "123", makeCliente("123", "Pepe"), nil, false},
		{"RepoError", "123", nil, assert.AnError, true},
		{"NotFound", "123", nil, nil, false}, //simula no encontrado
		{"EmptyID", "", makeCliente("", "Pepe"), assert.AnError, true},   //simula id vacío
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, mockRepo := setupClienteServiceWithMock(t)
			mockRepo.On("GetByID", mock.Anything, tt.mockID).Return(tt.mockData, tt.mockErr)
			got, err := s.GetByID(context.Background(), tt.mockID)

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

// funciones auxiliares
func makeCliente(id string, name string) *domain.Cliente {
	return &domain.Cliente{
		ID:                 id,
		Nombre:             name,
		Telefono:           "123456789",
		PreferenciaHoraria: domain.PreferenciaHoraria(1),
	}
}

func setupClienteServiceWithMock(t *testing.T) (cliente.ClienteService, *MockClienteRepository) {
	mockRepo := new(MockClienteRepository)
	s := cliente.NewClienteService(mockRepo)
	return s, mockRepo
}

/* func TestClienteService_Create(t *testing.T) {
	//instancia del mock
	//en este caso el mock, remplaza a la implementación real de ClienteRepository
	mockRepo := new(MockClienteRepository)

	/*func NewClienteService(repo repository.ClienteRepository) *clienteService {
	return &clienteService{repo: repo}
	}
	s := cliente.NewClienteService(mockRepo) //se instancia el service, es como usar

	input := &domain.Cliente{
		ID:                 "123",
		Nombre:             "Ivan",
		Telefono:           "123456789",
		PreferenciaHoraria: domain.PreferenciaHoraria(1),
	}

	// se configura el mock: si se llama Create con cualquier contexto y el input,
	// entonces debe devolver el mismo input y error nil
	mockRepo.On("CreateOrUpdate", mock.Anything, input).Return(input, nil)

	// Llamar al método Create del servicio
	result, err := s.Create(context.Background(), input)
	if err != nil {
		t.Fatalf("Error al crear cliente: %v", err)
	}
	//verifica que no hay error
	require.NoError(t, err)

	//require
	require.Equal(t, input, result) // Verifica que el resultado sea igual al input

	// Verifica que se llamaron todos los métodos esperados en el mock
	mockRepo.AssertExpectations(t)
}
*/

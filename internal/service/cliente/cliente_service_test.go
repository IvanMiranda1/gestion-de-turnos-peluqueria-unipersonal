package cliente_test

import (
	"context"
	"testing"

	"github.com/IvanMiranda1/gestion-de-turnos-peluqueria-unipersonal/internal/domain"
	cliente "github.com/IvanMiranda1/gestion-de-turnos-peluqueria-unipersonal/internal/service/cliente"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
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
	return args.Get(0).(*domain.Cliente), args.Error(1)
}

func (m *MockClienteRepository) GetAll(ctx context.Context) ([]*domain.Cliente, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.Cliente), args.Error(1)
}

func TestClienteService_Create(t *testing.T) {
	//instancia del mock
	//en este caso el mock, remplaza a la implementación real de ClienteRepository
	mockRepo := new(MockClienteRepository)

	/*func NewClienteService(repo repository.ClienteRepository) *clienteService {
	return &clienteService{repo: repo}
	}	*/
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

func TestClienteService_Update(t *testing.T) {
	mockRepo := new(MockClienteRepository)
	s := cliente.NewClienteService(mockRepo)
	input := &domain.Cliente{
		ID:                 "123",
		Nombre:             "Ivan",
		Telefono:           "123456789",
		PreferenciaHoraria: domain.PreferenciaHoraria(1),
	}

	clienteActualizado := &domain.Cliente{
		ID:                 "123",
		Nombre:             "Ivan Actualizado",
		Telefono:           "987654321",
		PreferenciaHoraria: domain.PreferenciaHoraria(2),
	}

	mockRepo.On("CreateOrUpdate", mock.Anything, input).Return(clienteActualizado, nil)
	res, err := s.Update(context.Background(), input)
	assert.NoError(t, err)
	assert.Equal(t, clienteActualizado, res)

	assert.Error(t, err)
	assert.EqualError(t, err, "error simulado")

	mockRepo.AssertExpectations(t)
}

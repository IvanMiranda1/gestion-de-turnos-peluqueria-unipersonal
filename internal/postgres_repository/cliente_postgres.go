package postgresrepository

import (
	"context"
	"database/sql"

	"github.com/IvanMiranda1/gestion-de-turnos-peluqueria-unipersonal/internal/domain"
)

type ClientePostgresRepository struct {
	db *sql.DB
}

func NewClientePostgresRepository(db *sql.DB) *ClientePostgresRepository {
	return &ClientePostgresRepository{db: db}
}

func (r *ClientePostgresRepository) CreateOrUpdate(ctx context.Context, c *domain.Cliente) (*domain.Cliente, error) {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO cliente(id, nombre, telefono, preferenciahoraria)
	 VALUES ($1, $2, $3, $4)
	 ON CONFLICT (id)
	 DO UPDATE SET nombre = EXCLUDED.nombre,
	               telefono = EXCLUDED.telefono,
	               preferenciahoraria = EXCLUDED.preferenciahoraria`,
		c.ID, c.Nombre, c.Telefono, c.PreferenciaHoraria)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *ClientePostgresRepository) GetByID(ctx context.Context, id string) (*domain.Cliente, error) {
	var c domain.Cliente
	err := r.db.QueryRowContext(ctx,
		`SELECT id, nombre, telefono, preferenciahoraria FROM cliente WHERE id = $1 `, id).
		Scan(&c.ID, &c.Nombre, &c.Telefono, &c.PreferenciaHoraria)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *ClientePostgresRepository) GetAll(ctx context.Context) ([]*domain.Cliente, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, nombre, telefono, preferenciahoraria FROM cliente`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clientes []*domain.Cliente
	for rows.Next() {
		var c domain.Cliente
		if err := rows.Scan(&c.ID, &c.Nombre, &c.Telefono, &c.PreferenciaHoraria); err != nil {
			return nil, err
		}
		clientes = append(clientes, &c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return clientes, nil
}

func (r *ClientePostgresRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM cliente WHERE id = $1`, id)
	return err
}

package postgresrepository

import (
	"context"
	"database/sql"
	"time"

	"github.com/IvanMiranda1/gestion-de-turnos-peluqueria-unipersonal/internal/domain"
)

type TurnoPostgresRepository struct {
	db *sql.DB
}

func NewTurnoPostgresRepository(db *sql.DB) *TurnoPostgresRepository {
	return &TurnoPostgresRepository{db: db}
}

func (r *TurnoPostgresRepository) CreateOrUpdate(ctx context.Context, t *domain.Turno) (*domain.Turno, error) {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO turno(id, fecha, hora, cliente_id)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT(id)
	DO UPDATE SET fecha = EXCLUDED.fecha,
	hora = EXCLUDED.hora,
	cliente_id = EXCLUDED.cliente_id`,
		t.ID, t.Fecha, t.Hora.String(), t.Cliente)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (r *TurnoPostgresRepository) GetAll(ctx context.Context) ([]*domain.Turno, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, fecha, hora, cliente_id FROM turno`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var turnos []*domain.Turno
	for rows.Next() {
		var t domain.Turno
		if err := rows.Scan(&t.ID, &t.Fecha, &t.Hora, &t.Cliente); err != nil {
			return nil, err
		}
		turnos = append(turnos, &t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return turnos, nil
}

func (r *TurnoPostgresRepository) GetByFecha(ctx context.Context, fecha time.Time) ([]*domain.Turno, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT 
		t.id, t.fecha, t.hora, t.cliente_id, 
		c.id, c.nombre, c.telefono, c.preferenciahoraria 
		FROM turno t 
		INNER JOIN cliente c ON t.cliente_id = c.id 
		WHERE t.fecha = $1`, fecha)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var turnos []*domain.Turno
	for rows.Next() {
		var t domain.Turno
		if err := rows.Scan(&t.ID, &t.Fecha, &t.Hora, &t.Cliente.ID, &t.Cliente.Nombre, &t.Cliente.Telefono, &t.Cliente.PreferenciaHoraria); err != nil {
			return nil, err
		}
		turnos = append(turnos, &t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return turnos, nil
}

func (r *TurnoPostgresRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM turno WHERE id = $1`, id)
	return err
}

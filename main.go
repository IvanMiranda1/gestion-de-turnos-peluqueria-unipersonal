package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/IvanMiranda1/gestion-de-turnos-peluqueria-unipersonal/internal/handler"
	postgresrepository "github.com/IvanMiranda1/gestion-de-turnos-peluqueria-unipersonal/internal/postgres_repository"
	"github.com/IvanMiranda1/gestion-de-turnos-peluqueria-unipersonal/internal/service/cliente"
	"github.com/IvanMiranda1/gestion-de-turnos-peluqueria-unipersonal/internal/service/turno"
	"github.com/go-chi/chi/v5"
)

func main() {
	db, err := sql.Open("postgres", "postgres://admin:admin123@localhost:5432/turnos?sslmode=disable")
	if err != nil {
		panic(err)
	}
	clienteRepo := postgresrepository.NewClientePostgresRepository(db)
	turnoRepo := postgresrepository.NewTurnoPostgresRepository(db)

	clienteService := cliente.NewClienteService(clienteRepo)
	turnoService := turno.NewTurnoService(turnoRepo, clienteService)

	clienteHandler := handler.NewClienteHandler(clienteService)
	turnoHandler := handler.NewTurnoHandler(turnoService)

	router := chi.NewRouter()
	router.Route("/cliente", clienteHandler.RegisterRoutes)
	router.Route("/turno", turnoHandler.RegisterRoutes)

	log.Printf("Server is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

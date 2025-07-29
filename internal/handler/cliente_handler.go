package handler

import (
	"encoding/json"
	"net/http"

	"github.com/IvanMiranda1/gestion-de-turnos-peluqueria-unipersonal/internal/dto"
	"github.com/IvanMiranda1/gestion-de-turnos-peluqueria-unipersonal/internal/service/cliente"
	"github.com/IvanMiranda1/gestion-de-turnos-peluqueria-unipersonal/pkg/web"
	"github.com/go-chi/chi/v5"
)

type ClienteHandler struct {
	s cliente.ClienteService
}

func NewClienteHandler(s cliente.ClienteService) *ClienteHandler {
	return &ClienteHandler{s: s}
}

func (h *ClienteHandler) RegisterRoutes(r chi.Router) {
	r.Post("/", h.Create)
	r.Put("/{id}", h.Update)
	r.Get("/{id}", h.GetByID)
	r.Get("/", h.GetAll) //GET /cliente
	r.Delete("/{id}", h.Delete)
}

func (h *ClienteHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.ClienteRequest // lo que se espera recibir
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		web.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	c, err := req.ToDomain()
	if err != nil {
		web.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.s.Create(r.Context(), c)
	if err != nil {
		web.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	web.Success(w, http.StatusCreated, dto.ClienteFromDomain(res))
}

/*
 	r.Body contiene el cuerpo crudo de la solicitud HTTP.

    json.NewDecoder(r.Body).Decode(&req) intenta deserializar ese cuerpo (en formato JSON) y cargarlo en la estructura req.

    Si el JSON está mal formado o no coincide con la estructura esperada, devuelve error.

    Si hay error, se responde con HTTP 400 (Bad Request) y se corta la ejecución del handler.

	En resumen: parsea el JSON recibido del cliente y lo convierte en una estructura de Go (ClienteRequest).
*/

func (h *ClienteHandler) Update(w http.ResponseWriter, r *http.Request) {
	var req dto.ClienteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		web.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	c, err := req.ToDomain()
	if err != nil {
		web.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	res, err := h.s.Update(r.Context(), c)
	if err != nil {
		web.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	web.Success(w, http.StatusOK, dto.ClienteFromDomain(res))
}

func (h *ClienteHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		web.Error(w, http.StatusBadRequest, "id is required")
		return
	}

	res, err := h.s.GetByID(r.Context(), id)
	if err != nil {
		web.Error(w, http.StatusNotFound, err.Error())
		return
	}
	web.Success(w, http.StatusOK, dto.ClienteFromDomain(res))
}

func (h *ClienteHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	res, err := h.s.GetAll(r.Context())
	if err != nil {
		web.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	clienteSlice := make([]any, 0, len(res))
	for _, c := range res {
		clienteSlice = append(clienteSlice, dto.ClienteFromDomain(c))
	}
	web.Success(w, http.StatusOK, clienteSlice)
}

func (h *ClienteHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		web.Error(w, http.StatusBadRequest, "id is required")
		return
	}

	if err := h.s.Delete(r.Context(), id); err != nil {
		web.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

/*
`http.ResponseWriter` y `*http.Request` son los componentes centrales en un handler HTTP en Go:

* `*http.Request` representa la solicitud entrante del cliente. Contiene:

  * URL, método (GET, POST, etc.)
  * Headers
  * Cuerpo (`Body`) con datos JSON, formulario, etc.
  * Parámetros (query, path)
  * Contexto (`Context()`) para cancelación, deadlines, valores compartidos.

* `http.ResponseWriter` es el mecanismo para construir la respuesta HTTP:

  * Escribir código de estado (`WriteHeader`).
  * Escribir headers (Content-Type, Authorization, etc.).
  * Escribir el cuerpo (por ejemplo, JSON serializado).
  * Se usa para enviar datos al cliente.

Un handler es una función con firma:

```go
func(w http.ResponseWriter, r *http.Request)
```

Cada vez que llega una petición HTTP, el servidor llama a esta función con esos parámetros.

En resumen:

* Usás `r *http.Request` para leer datos del cliente.
* Usás `w http.ResponseWriter` para enviar la respuesta.

Todo el ciclo de comunicación HTTP se basa en esos dos objetos.

*/

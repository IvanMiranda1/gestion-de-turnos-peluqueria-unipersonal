# ✂️ Gestión de Turnos — Peluquería Unipersonal

API REST para la gestión de turnos y clientes de una peluquería unipersonal. Desarrollada en Go con arquitectura en capas (handler → service → repository) y base de datos PostgreSQL.

> ⚠️ Proyecto en desarrollo. Algunas funcionalidades pueden estar incompletas.

---

## 🛠️ Tecnologías

- **Go** — lenguaje principal
- **PostgreSQL** — base de datos
- **Chi** — router HTTP
- **Docker / Docker Compose** — para levantar la base de datos fácilmente

---

## 📁 Estructura del proyecto

```
.
├── main.go                  # Punto de entrada, configuración de rutas
├── internal/
│   ├── handler/             # Handlers HTTP (cliente, turno)
│   ├── service/
│   │   ├── cliente/         # Lógica de negocio de clientes
│   │   └── turno/           # Lógica de negocio de turnos
│   └── postgres_repository/ # Acceso a la base de datos
├── database/                # Scripts SQL / migraciones
├── pkg/web/                 # Utilidades web compartidas
├── docker-compose.yml
├── go.mod
└── go.sum
```

---

## 🚀 Cómo correrlo

### 1. Levantar la base de datos

```bash
docker-compose up -d
```

Esto levanta una instancia de PostgreSQL con las credenciales configuradas en `main.go`.

### 2. Correr la aplicación

```bash
go run main.go
```

El servidor queda escuchando en `http://localhost:8080`.

---

## 🔌 Endpoints

### Clientes

| Método | Ruta | Descripción |
|--------|------|-------------|
| `GET` | `/cliente` | Listar clientes |
| `POST` | `/cliente` | Crear un cliente |
| `GET` | `/cliente/{id}` | Obtener un cliente |
| `PUT` | `/cliente/{id}` | Actualizar un cliente |
| `DELETE` | `/cliente/{id}` | Eliminar un cliente |

### Turnos

| Método | Ruta | Descripción |
|--------|------|-------------|
| `GET` | `/turno` | Listar turnos |
| `POST` | `/turno` | Crear un turno |
| `GET` | `/turno/{id}` | Obtener un turno |
| `PUT` | `/turno/{id}` | Actualizar un turno |
| `DELETE` | `/turno/{id}` | Eliminar un turno |

> Los endpoints exactos pueden variar según el estado actual del desarrollo.

---

## ⚙️ Configuración de la base de datos

Por defecto la app se conecta con:

```
postgres://admin:admin123@localhost:5432/turnos?sslmode=disable
```

Para cambiar la cadena de conexión, modificar el `main.go` o configurar una variable de entorno (pendiente de implementar).

---

## 📝 Notas

- El proyecto no está terminado. Fue desarrollado como práctica de Go con arquitectura en capas.
- Las migraciones de base de datos están en la carpeta `database/`.

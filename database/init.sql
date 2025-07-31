CREATE TABLE cliente (
    id TEXT PRIMARY KEY,
    nombre TEXT NOT NULL,
    telefono TEXT NOT NULL,
    preferenciahoraria TEXT NOT NULL
);

CREATE TABLE turno (
    id TEXT PRIMARY KEY,
    fecha DATE NOT NULL,
    hora TEXT NOT NULL,
    cliente_id TEXT NOT NULL REFERENCES cliente(id)
);

# ImageFlow - Aplicación Go

Prototipo cliente-servidor para cargar, consultar, buscar, descargar y eliminar imágenes. El cliente y el servidor están implementados en Go y utilizan la biblioteca estándar; no hay dependencias externas.

## Requisitos

- Go 1.22 o superior.
- macOS, Linux o Windows con Go instalado.

Comprobar la instalación:

```bash
go version
```

## Iniciar el servidor

Desde esta carpeta (`app-go`):

```bash
go run ./cmd/server
```

El servidor escucha por defecto en `0.0.0.0:8080` y prepara los directorios de datos automáticamente.

## Iniciar el cliente

Abre una segunda terminal en esta misma carpeta:

```bash
go run ./cmd/client
```

El cliente se conecta por defecto a `http://localhost:8080` y presenta un menú de operaciones.

## Configuración

| Variable | Valor predeterminado | Descripción |
|---|---|---|
| `IMAGEFLOW_HOST` | `0.0.0.0` | Interfaz en la que escucha el servidor. |
| `IMAGEFLOW_PORT` | `8080` | Puerto HTTP del servidor. |
| `IMAGEFLOW_DATA_DIR` | `./data` | Carpeta de almacenamiento del servidor. |
| `IMAGEFLOW_SERVER_URL` | `http://localhost:8080` | Dirección del servidor para el cliente. |

Ejemplo en macOS/Linux para cambiar el puerto:

```bash
IMAGEFLOW_PORT=8081 go run ./cmd/server
```

En la terminal del cliente:

```bash
IMAGEFLOW_SERVER_URL=http://localhost:8081 go run ./cmd/client
```

Para conectar desde otro equipo de la red, reemplaza `localhost` por la IP del equipo servidor, por ejemplo `http://192.168.1.20:8080`. Permite el puerto en el firewall y comprueba que ambos equipos tengan conectividad.

## Operaciones del menú

1. Cargar una imagen. El cliente solicita la ruta local y una categoría.
2. Consultar todas las imágenes.
3. Buscar por nombre y/o categoría.
4. Descargar una imagen por su ID y guardarla en una ruta indicada.
5. Eliminar una imagen por su ID, solicitando confirmación.
0. Salir.

## Demostración reproducible

Con el servidor encendido, en el cliente:

1. Selecciona `1`.
2. Ingresa `testdata/imagen_demo.png` (si la terminal está situada en `app-go`).
3. Ingresa la categoría `demo`.
4. Selecciona `2` para comprobar que quedó registrada.
5. Selecciona `3`, escribe `imagen_demo` como nombre y `demo` como categoría.
6. Selecciona `4`, pega el ID devuelto y escribe `descargas/imagen_demo.png` como destino.
7. Para comprobar una validación, selecciona `1` e ingresa `testdata/no_es_imagen.txt`. El servidor debe rechazar el archivo por formato.
8. Para comprobar la eliminación, selecciona `5`, pega el ID y confirma con `s`.

## Formatos y límites

Formatos aceptados: JPG/JPEG, PNG, GIF y WEBP. Tamaño máximo: 10 MiB. El servidor detecta el tipo a partir del contenido recibido y no confía únicamente en la extensión del archivo.

## Almacenamiento

- Archivos binarios: `data/images/`.
- Metadatos: `data/metadata.json`.
- Los datos permanecen después de reiniciar el servidor.

## Pruebas

Desde esta carpeta:

```bash
go test ./...
go vet ./...
```

Las pruebas cubren el endpoint de salud, comunicación cliente-servidor, carga, consulta, descarga, eliminación, rechazo de formatos, límite de tamaño, validación del nombre, persistencia y recurso inexistente.

## Estructura

```text
app-go/
├── cmd/
│   ├── client/main.go
│   └── server/main.go
├── internal/
│   ├── client/api.go
│   ├── handlers/image_handler.go
│   ├── models/image.go
│   ├── services/image_service.go
│   ├── storage/local.go
│   └── transport/server.go
├── data/
├── testdata/
└── go.mod
```

## Limitaciones de esta fase

- La persistencia se realiza con archivos locales y JSON; no se utiliza una base de datos relacional.
- No hay autenticación ni gestión de roles en este prototipo. La guía ubica la integración con el servicio de autenticación en una fase posterior.
- La conectividad entre sedes se desarrollará en las siguientes fases.
- No se generan miniaturas ni se realiza análisis antivirus.

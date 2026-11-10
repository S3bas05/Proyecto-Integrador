# ImageFlow

**Prototipo empresarial para la gestión de imágenes**  
Segundo entregable del Proyecto Integrador - Semana 10  
Asignatura: Programación Orientada a Objetos

ImageFlow es una aplicación cliente-servidor para cargar, registrar, consultar, buscar, descargar y eliminar imágenes desde un cliente de consola desarrollado en Go. El servidor también está desarrollado en Go.

## Alcance de esta entrega

Esta versión corresponde al prototipo funcional inicial de la semana 10. El código ejecutable de la aplicación está escrito en **Go** y utiliza únicamente la biblioteca estándar del lenguaje.

La integración con el servicio de autenticación FastAPI y la conexión entre sedes no forman parte de esta semana; la guía académica las ubica en fases posteriores. En esta entrega no se duplica el componente de autenticación.

## Funcionalidades implementadas

- Cargar imágenes desde el cliente al servidor.
- Validar formato de imagen y tamaño máximo de 10 MiB.
- Registrar ID, nombre, tipo MIME, tamaño, categoría y fecha.
- Guardar los archivos binarios por separado de los metadatos JSON.
- Consultar todas las imágenes.
- Buscar por nombre y/o categoría.
- Descargar una imagen por su identificador.
- Eliminar una imagen y su registro.
- Informar errores de conexión, datos inválidos, formatos no permitidos, archivos demasiado grandes y recursos inexistentes.

## Requisitos

- Go 1.22 o superior.
- Git, si se desea clonar el repositorio.
- No se requieren dependencias externas de Go.

## Estructura

```text
imageflow/
├── app-go/
│   ├── cmd/
│   │   ├── client/main.go
│   │   └── server/main.go
│   ├── internal/
│   │   ├── client/api.go
│   │   ├── handlers/image_handler.go
│   │   ├── models/image.go
│   │   ├── services/image_service.go
│   │   ├── storage/local.go
│   │   └── transport/server.go
│   ├── data/
│   ├── testdata/
│   ├── go.mod
│   └── README.md
├── docs/
│   ├── phase2/
│   ├── requirements/
│   ├── traceability/
│   └── use-cases/
├── LICENSE
└── README.md
```

## Configuración

| Variable | Valor predeterminado | Uso |
|---|---|---|
| `IMAGEFLOW_HOST` | `0.0.0.0` | Interfaz donde escucha el servidor. |
| `IMAGEFLOW_PORT` | `8080` | Puerto HTTP. |
| `IMAGEFLOW_DATA_DIR` | `./data` | Carpeta donde se guardan imágenes y metadatos. |
| `IMAGEFLOW_SERVER_URL` | `http://localhost:8080` | URL a la que se conecta el cliente. |

## Endpoints HTTP

| Método | Ruta | Función |
|---|---|---|
| GET | `/health` | Comprueba la disponibilidad del servidor. |
| POST | `/api/v1/images` | Carga una imagen mediante `multipart/form-data`; campos `file` y `category`. |
| GET | `/api/v1/images` | Consulta imágenes; admite filtros `name` y `category`. |
| GET | `/api/v1/images/{id}/download` | Descarga una imagen por ID. |
| DELETE | `/api/v1/images/{id}` | Elimina una imagen por ID. |

## Almacenamiento y decisiones técnicas

- Archivos binarios: `app-go/data/images/`.
- Metadatos: `app-go/data/metadata.json`.
- Formatos aceptados: JPEG, PNG, GIF y WEBP.
- Tamaño máximo: 10 MiB (10 × 1024 × 1024 bytes).
- Los metadatos se serializan como JSON y los archivos se guardan en nombres internos generados por el servidor.
- `sync.RWMutex` protege el acceso concurrente al almacenamiento local.

## Requisitos pendientes y limitaciones

| Elemento | Estado |
|---|---|
| Cliente y servidor independientes en Go | Implementado |
| Comunicación HTTP real | Implementado |
| Carga, consulta, búsqueda, descarga y eliminación | Implementado |
| Persistencia local de archivos y metadatos | Implementado |
| Validaciones y errores | Implementado |
| Integración con autenticación FastAPI | Fase posterior según la guía académica |
| Comunicación entre sedes y configuración de red | Fases posteriores |
| Base de datos relacional, miniaturas y análisis antivirus | Fuera del alcance de este prototipo |

## Autoría, licencia 

**Autor:** Sebastián Alexis Cumbal Molineros.  
La presente entrega se desarrolla de forma individual. El autor asume el análisis, programación del cliente y servidor, modelado, validaciones, pruebas y documentación. La modalidad individual debe quedar registrada y autorizada por el docente.

Licencia MIT. Consultar `LICENSE`.

## Enlaces de entrega

- Repositorio Git: (https://github.com/S3bas05/Proyecto-Integrador.git)
- Video demostrativo: **PENDIENTE: pegar URL del video.**

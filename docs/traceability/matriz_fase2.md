# Matriz de trazabilidad - ImageFlow Fase 2

| Requisito | Componentes principales | Evidencia de verificación |
|---|---|---|
| P2-01 Cargar imagen | `internal/client`, `internal/handlers`, `internal/services` | Prueba `TestClientServerImageOperations` y menú opción 1. |
| P2-02 Guardar metadatos | `internal/models`, `internal/storage` | `data/metadata.json` y respuesta JSON de carga. |
| P2-03 Consultar imágenes | `internal/handlers`, `internal/client` | Menú opción 2 y prueba de integración. |
| P2-04 Buscar por nombre/categoría | `internal/services`, `internal/client` | Menú opción 3 y parámetros `name` y `category`. |
| P2-05 Descargar imagen | `internal/handlers`, `internal/client` | Menú opción 4 y comparación de bytes en prueba. |
| P2-06 Eliminar imagen | `internal/storage`, `internal/services`, `internal/client` | Menú opción 5 y verificación de lista vacía. |
| P2-07 Validar formato/tamaño/nombre | `internal/services`, `internal/handlers` | Pruebas de formato no permitido, tamaño excesivo y nombre con ruta. |
| P2-08 Comunicación independiente | `cmd/client`, `cmd/server`, `internal/transport` | Cliente y servidor iniciados en procesos separados; pruebas HTTP. |
| P2-09 Configurar IP y puerto | `cmd/client`, `cmd/server` | Variables `IMAGEFLOW_SERVER_URL`, `IMAGEFLOW_HOST` e `IMAGEFLOW_PORT`. |
| P2-10 Persistencia | `internal/storage` | Prueba `TestLocalStoragePersistsImageAndMetadata`. |
| P2-11 Manejo de errores | `internal/handlers`, `internal/client` | Formato incorrecto, ID inexistente y error de conexión. |
| P2-12 Autenticación | No implementado en semana 10 | Fase posterior según la guía académica. |
| P2-13 Comunicación entre sedes | No implementado en semana 10 | Fases posteriores de infraestructura de red. |

## Nota de alcance

La entrega de semana 10 implementa el prototipo de gestión de imágenes con cliente y servidor Go. La integración del servicio de autenticación y la infraestructura entre sedes se mantiene fuera del alcance de esta fase, conforme a las instrucciones de la actividad.

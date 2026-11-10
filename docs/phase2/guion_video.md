# Guion breve para el video demostrativo - ImageFlow Fase 2

Duración sugerida: 3 a 5 minutos. Grabar la pantalla mostrando el funcionamiento real, no solo diapositivas.

## Secuencia de demostración

1. **Presentación (15-20 s):** indicar que ImageFlow es un prototipo de gestión de imágenes con cliente y servidor independientes, ambos escritos en Go.
2. **Estructura del proyecto (30 s):** mostrar `cmd/client`, `cmd/server` e `internal` con los paquetes `client`, `handlers`, `models`, `services`, `storage` y `transport`.
3. **Servidor (20 s):** abrir una terminal en `app-go` y ejecutar `go run ./cmd/server`. Señalar la IP/interfaz, el puerto y la ruta de salud.
4. **Cliente (20 s):** abrir otra terminal y ejecutar `go run ./cmd/client`. Mostrar que se establece la conexión HTTP.
5. **Carga y consulta (45 s):** elegir la opción 1, cargar `testdata/imagen_demo.png` con categoría `demo`, copiar el ID y elegir la opción 2.
6. **Búsqueda y descarga (40 s):** elegir la opción 3 y filtrar por nombre/categoría. Luego elegir opción 4, usar el ID y descargar a `descargas/imagen_demo.png`.
7. **Validación (25 s):** intentar cargar `testdata/no_es_imagen.txt` y mostrar el rechazo por formato no permitido.
8. **Eliminación (20 s):** eliminar la imagen usando la opción 5 y confirmar con `s`.
9. **Cierre técnico (20 s):** mostrar `go test ./...`, el almacenamiento diferenciado y explicar qué queda para fases posteriores.

## Comandos para preparar la grabación

Terminal 1:

```bash
cd app-go
go run ./cmd/server
```

Terminal 2:

```bash
cd app-go
go run ./cmd/client
```

## Guion de explicación técnica

- El cliente no accede directamente a los archivos del servidor; envía solicitudes HTTP.
- El servidor recibe las solicitudes, el handler delega al servicio y el servicio utiliza el almacenamiento.
- `Image` utiliza campos privados y un constructor para evitar objetos inválidos.
- Los archivos binarios se guardan en `data/images/` y los metadatos en `data/metadata.json`.
- El proyecto informa errores sin utilizar `panic` para validaciones normales.

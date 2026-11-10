# Requisitos del prototipo ImageFlow - Fase 2

Los siguientes requisitos corresponden a la opción C: sistema de gestión de imágenes. El alcance de esta matriz es el prototipo funcional de la semana 10.

| ID | Requisito | Criterio de aceptación | Estado |
|---|---|---|---|
| P2-01 | Cargar una imagen | El cliente envía un archivo y el servidor confirma su registro. | Implementado |
| P2-02 | Guardar metadatos | Se registra ID, nombre, tipo MIME, tamaño, categoría y fecha. | Implementado |
| P2-03 | Consultar imágenes | El cliente recibe el listado de imágenes disponibles. | Implementado |
| P2-04 | Buscar o clasificar | Se filtra el listado por nombre y/o categoría. | Implementado |
| P2-05 | Descargar una imagen | El cliente recupera el archivo original mediante su ID. | Implementado |
| P2-06 | Eliminar una imagen | Se elimina el archivo y su registro de metadatos. | Implementado |
| P2-07 | Validar archivos | Se rechazan formatos no permitidos y archivos mayores a 10 MiB. | Implementado |
| P2-08 | Comunicación cliente-servidor | Cliente y servidor son procesos independientes que intercambian solicitudes HTTP. | Implementado |
| P2-09 | Configurar IP y puerto | El servidor y el cliente admiten configuración por variables de entorno. | Implementado |
| P2-10 | Persistir la información | Los binarios y metadatos se conservan tras reiniciar el servidor. | Implementado |
| P2-11 | Manejar errores | Se informan datos inválidos, formato no permitido, archivo grande, recurso inexistente y errores de comunicación. | Implementado |
| P2-12 | Integrar autenticación | Consumir el servicio de autenticación sin duplicar su lógica en Go. | Fase posterior según la guía |
| P2-13 | Comunicación entre sedes | Acceder al servidor mediante la infraestructura de red del proyecto. | Fases posteriores |

## Decisiones técnicas

- Lenguaje de programación de cliente y servidor: Go.
- Dependencias externas: ninguna; se utiliza la biblioteca estándar.
- Protocolo: HTTP.
- Formato de metadatos transmitidos: JSON.
- Carga de archivos: `multipart/form-data`.
- Almacenamiento binario: directorio `data/images/`.
- Almacenamiento de metadatos: `data/metadata.json`.
- Formatos aceptados: JPEG, PNG, GIF y WEBP.
- Tamaño máximo: 10 MiB.
- Búsqueda: coincidencia parcial, sin distinguir mayúsculas/minúsculas, por nombre y categoría.

## Requisitos pendientes y límites

La autenticación y la autorización se integrarán en la fase indicada por la guía académica. La conectividad entre sedes se incorporará cuando se implemente la infraestructura de red. En esta versión no se utiliza una base de datos relacional ni se generan miniaturas.
